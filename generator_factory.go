package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type generatorFactory struct {
	dir        string
	pkg        string
	generators map[string]*Generator

	curFset  *token.FileSet
	lastType string
}

func NewGenerators(symbols []string, dir string) ([]*Generator, error) {
	factory := &generatorFactory{dir: dir}

	if err := factory.walkDir(factory.inspectPkg); err != nil {
		return nil, fmt.Errorf("factory.walkDir: %w", err)
	}

	if err := factory.initGenerators(symbols); err != nil {
		return nil, fmt.Errorf("factory.initGenerators: %w", err)
	}

	if err := factory.walkDir(factory.inspectDeclaration); err != nil {
		return nil, fmt.Errorf("factory.walkDir: %w", err)
	}

	var result []*Generator
	for _, generator := range factory.generators {
		result = append(result, generator)
	}
	return result, nil
}

func (f *generatorFactory) walkDir(fn func(file *ast.File) error) error {
	if f.dir == "" {
		return fmt.Errorf("no dir specified")
	}
	return filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}
		f.curFset = token.NewFileSet()
		file, err := parser.ParseFile(f.curFset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parser.ParseFile %s: %w", path, err)
		}
		return fn(file)
	})
}

func (f *generatorFactory) inspectPkg(file *ast.File) error {
	switch pkg := file.Name.Name; {
	case f.pkg == "":
		f.pkg = pkg
	case f.pkg != pkg:
		return fmt.Errorf("package name mismatch: %s!= %s", f.pkg, pkg)
	}
	return nil
}

func (f *generatorFactory) initGenerators(symbols []string) error {
	if len(symbols) == 0 || f.dir == "" || f.pkg == "" {
		return fmt.Errorf("these fields must be non-empty, symbols %s, f.dir %s, f.pkg %s", symbols, f.dir, f.pkg)
	}

	generators := make(map[string]*Generator, len(symbols))
	for _, symbol := range symbols {
		generators[symbol] = &Generator{
			Name: symbol,
			Dir:  f.dir,
			Pkg:  f.pkg,
		}
	}

	f.generators = generators
	return nil
}

func (f *generatorFactory) inspectDeclaration(file *ast.File) error {
	var err error
	ast.Inspect(file, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.GenDecl:
			err = f.inspectGenericDeclaration(decl)
			if err != nil {
				err = fmt.Errorf("f.inspectGenericDeclaration: %w", err)
				return false
			}
		case *ast.FuncDecl:
			err = f.inspectFunctionDeclaration(decl)
			if err != nil {
				err = fmt.Errorf("f.inspectFunctionDeclaration: %w", err)
				return false
			}
		}
		return true
	})
	return err
}

func (f *generatorFactory) inspectGenericDeclaration(decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.ValueSpec:
			if err := f.inspectValueSpec(spec); err != nil {
				return fmt.Errorf("f.inspectValueSpec: %w", err)
			}
		case *ast.TypeSpec:
			if err := f.inspectTypeSpec(spec); err != nil {
				return fmt.Errorf("f.inspectTypeSpec: %w", err)
			}
		}
	}
	return nil
}

func (f *generatorFactory) inspectValueSpec(spec *ast.ValueSpec) error {
	var specType string
	if t := spec.Type; t != nil {
		s, err := parseNode(f.curFset, t)
		if err != nil {
			return fmt.Errorf("parseNode: %w", err)
		}
		specType = s
	}

	lit, _ := parseNode(f.curFset, spec)
	debug.Printf("Inspecting %s", lit)

	for i, name := range spec.Names {
		debug.Printf("Inspecting %dth name %s ......\n", i, name.Name)
		generator, ok := f.generators[name.Name]
		if !ok {
			continue
		}

		generator.Type = specType
		if generator.Type == "" {
			if i >= len(spec.Values) {
				// for const declaration list, the expression may be omitted,
				// which leads to a shorter Values slice.
				generator.Type = f.lastType
				debug.Printf("Assign lastType %s to %s", f.lastType, name.Name)
				continue
			}
			debug.Printf("spec.Values[%d] is %T\n", i, spec.Values[i])
			kind := ""
			switch v := spec.Values[i].(type) {
			case *ast.BasicLit:
				kind = v.Kind.String()
			case *ast.UnaryExpr:
				if bl, ok := v.X.(*ast.BasicLit); ok {
					kind = bl.Kind.String()
				}
			case *ast.Ident: // case for iota
				if v.Name == "iota" {
					kind = "INT"
				}
			}

			switch kind {
			case "INT":
				generator.Type = "int"
			case "FLOAT":
				generator.Type = "float64"
			case "IMAG":
				generator.Type = "complex128"
			case "CHAR":
				generator.Type = "rune"
			case "STRING":
				generator.Type = "string"
			}
		}
		generator.Type = strings.TrimSpace(generator.Type)
		debug.Printf("generator.Type = %s\n", generator.Type)
		if generator.Type == "" {
			return fmt.Errorf("can't infer type for '%s'", name.Name)
		}
		debug.Printf("Type of '%s' is %s\n", name.Name, generator.Type)
		f.lastType = generator.Type
	}
	return nil
}

func (f *generatorFactory) inspectTypeSpec(spec *ast.TypeSpec) error {
	// TODO implement this later
	return nil
}

func (f *generatorFactory) inspectFunctionDeclaration(decl *ast.FuncDecl) error {
	// TODO implement this later
	return nil
}

func parseNode(fset *token.FileSet, node any) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return "", fmt.Errorf("printer.Fprint: %w", err)
	}
	return buf.String(), nil
}

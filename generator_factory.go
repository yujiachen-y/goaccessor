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

	curFileName string
	curFset     *token.FileSet
	lastType    string
}

func NewGenerators(targets []string, dir string, field bool) ([]*Generator, error) {
	factory := &generatorFactory{dir: dir}

	if err := factory.walkDir(factory.inspectPkg); err != nil {
		return nil, fmt.Errorf("factory.walkDir: %w", err)
	}

	if err := factory.initGenerators(targets); err != nil {
		return nil, fmt.Errorf("factory.initGenerators: %w", err)
	}

	if err := factory.walkDir(factory.inspectDeclaration); err != nil {
		return nil, fmt.Errorf("factory.walkDir: %w", err)
	}

	if field {
		variables, err := factory.replaceVariableGenerators()
		if err != nil {
			return nil, fmt.Errorf("factory.replaceVariableGenerators: %w", err)
		}

		if err := factory.walkDir(factory.inspectDeclaration); err != nil {
			return nil, fmt.Errorf("factory.walkDir: %w", err)
		}

		if err := factory.insertVariablesGenerators(variables); err != nil {
			return nil, fmt.Errorf("insertVariablesGenerators: %w", err)
		}
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
		if info.IsDir() && path != f.dir {
			return filepath.SkipDir
		}
		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_goaccessor.go") {
			return nil
		}
		f.curFileName = strings.TrimSuffix(info.Name(), ".go")
		f.curFset = token.NewFileSet()
		file, err := parser.ParseFile(f.curFset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parser.ParseFile %s: %w", path, err)
		}
		debug.Printf("begin to parse file %s\n", path)
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

func (f *generatorFactory) initGenerators(targets []string) error {
	if len(targets) == 0 || f.dir == "" || f.pkg == "" {
		return fmt.Errorf("these fields must be non-empty, targets %s, f.dir %s, f.pkg %s", targets, f.dir, f.pkg)
	}

	generators := make(map[string]*Generator, len(targets))
	for _, target := range targets {
		generators[target] = &Generator{
			Name:    target,
			Dir:     f.dir,
			Pkg:     f.pkg,
			Fields:  make([]Field, 0),
			Methods: make(map[string]struct{}),
		}
	}

	f.generators = generators
	return nil
}

func (f *generatorFactory) replaceVariableGenerators() (variables map[string]*Generator, err error) {
	variables = f.generators
	types := make(map[string]*Generator, len(variables))
	for _, v := range variables {
		if v.GeneratorType != GeneratorTypeVariable {
			return nil, fmt.Errorf("unexpected variable generator type %s", v.GeneratorType)
		}

		referType := v.Type
		// simple solution for type parameter
		if idx := strings.Index(referType, "["); idx != -1 {
			referType = referType[:idx]
		}
		// simple solution for pointer type
		referType = strings.TrimPrefix(referType, "*")
		if _, ok := types[referType]; ok {
			continue
		}

		types[referType] = &Generator{
			Name:    referType,
			Dir:     f.dir,
			Pkg:     f.pkg,
			Fields:  make([]Field, 0),
			Methods: make(map[string]struct{}),
		}
		debug.Printf("Replace variable %+v with %+v\n", v, types[referType])
	}
	f.generators = types
	return
}

func (f *generatorFactory) insertVariablesGenerators(variables map[string]*Generator) error {
	newVariables := make(map[string]*Generator, len(variables))
	for k, v := range variables {
		if v.GeneratorType != GeneratorTypeVariable {
			return fmt.Errorf("unexpected variable generator type %s", v.GeneratorType)
		}

		// skip anonymous types
		if len(v.Fields) > 0 {
			v.GeneratorType = GeneratorTypeField
			newVariables[k] = v
			continue
		}

		referType := v.Type
		// simple solution for type parameter
		if idx := strings.Index(referType, "["); idx != -1 {
			referType = referType[:idx]
		}
		// simple solution for pointer type
		referType = strings.TrimPrefix(referType, "*")
		t, ok := f.generators[referType]
		if !ok {
			continue
		}

		if t.GeneratorType != GeneratorTypeStructure {
			return fmt.Errorf("unexpected structure generator type %s", t.GeneratorType)
		}

		newVariables[k] = &Generator{
			Name:          k,
			Dir:           v.Dir,
			Pkg:           v.Pkg,
			Type:          referType,
			TypeParams:    t.TypeParams,
			TypeArguments: v.TypeArguments,
			ReceiverName:  t.ReceiverName,
			Fields:        t.Fields,
			Methods:       t.Methods,
			GeneratorType: GeneratorTypeField,
			FileName:      v.FileName,
		}
		debug.Printf("Insert new variable %+v\n", newVariables[k])
	}
	f.generators = newVariables
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

// TODO refactor this
func (f *generatorFactory) inspectValueSpec(spec *ast.ValueSpec) error {
	var specType string
	var typeArguments []string
	var fields []Field
	if t := spec.Type; t != nil {
		s, err := parseNode(f.curFset, t)
		if err != nil {
			return fmt.Errorf("parseNode: %w", err)
		}
		specType = s

		if expr, ok := t.(*ast.StarExpr); ok {
			t = expr.X
		}
		// handler type arguments
		typeArguments = parseTypeArguments(t)

		// handler anonymous structs
		if structType, ok := t.(*ast.StructType); ok {
			fields, err = f.parseFields(structType.Fields.List)
			if err != nil {
				return fmt.Errorf("f.parseFields: %w", err)
			}
		}
	}

	lit, _ := parseNode(f.curFset, spec)
	debug.Printf("Inspecting %s", lit)

	for i, name := range spec.Names {
		debug.Printf("Inspecting %dth name %s ......\n", i, name.Name)
		generator, ok := f.generators[name.Name]
		if !ok {
			continue
		}

		generator.FileName = f.curFileName
		generator.GeneratorType = GeneratorTypeVariable
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
				s, err := f.parseUnaryExpression(v)
				if err != nil {
					return fmt.Errorf("f.parseUnaryExpression: %w", err)
				}
				kind = s

				if lit, ok := v.X.(*ast.CompositeLit); ok {
					// handler type arguments
					debug.Printf("type of lit.Type is %T", lit.Type)
					generator.TypeArguments = parseTypeArguments(lit.Type)
					// handler anonymous struct fields
					if structType, ok := lit.Type.(*ast.StructType); ok {
						fields, err := f.parseFields(structType.Fields.List)
						if err != nil {
							return fmt.Errorf("f.parseFields: %w", err)
						}
						generator.Fields = fields
					}
				}
			case *ast.Ident: // case for iota
				if v.Name == "iota" {
					kind = "INT"
				}
			case *ast.CompositeLit:
				s, err := parseNode(f.curFset, v.Type)
				if err != nil {
					return fmt.Errorf("parseNode: %w", err)
				}
				kind = s

				// handler type arguments
				generator.TypeArguments = parseTypeArguments(v.Type)
				// handler anonymous struct fields
				if structType, ok := v.Type.(*ast.StructType); ok {
					fields, err := f.parseFields(structType.Fields.List)
					if err != nil {
						return fmt.Errorf("f.parseFields: %w", err)
					}
					generator.Fields = fields
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
			default:
				generator.Type = kind
			}
		}
		generator.Type = strings.TrimSpace(generator.Type)
		debug.Printf("generator.Type = %s\n", generator.Type)
		if generator.Type == "" {
			return fmt.Errorf("can't infer type for '%s'", name.Name)
		}
		debug.Printf("Type of '%s' is %s\n", name.Name, generator.Type)
		f.lastType = generator.Type

		if len(typeArguments) > 0 && len(generator.TypeArguments) == 0 {
			generator.TypeArguments = typeArguments
		}
		if len(fields) > 0 && len(generator.Fields) == 0 {
			generator.Fields = fields
		}
	}
	return nil
}

func (f *generatorFactory) parseUnaryExpression(expr *ast.UnaryExpr) (string, error) {
	var kind string
	switch v := expr.X.(type) {
	case *ast.BasicLit:
		kind = v.Kind.String()
	case *ast.CompositeLit:
		s, err := parseNode(f.curFset, v.Type)
		if err != nil {
			return "", fmt.Errorf("parseNode: %w", err)
		}
		kind = s
	default:
		return "", fmt.Errorf("unsupported expression type: %T", v)
	}
	if kind == "" {
		return "", nil
	}
	if expr.Op == token.AND {
		kind = "*" + kind
	}
	return kind, nil
}

func (f *generatorFactory) inspectTypeSpec(spec *ast.TypeSpec) error {
	generator, ok := f.generators[spec.Name.Name]
	if !ok {
		return nil
	}

	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	generator.Type = spec.Name.Name
	if spec.TypeParams != nil {
		for _, typeParam := range spec.TypeParams.List {
			for _, name := range typeParam.Names {
				generator.TypeParams = append(generator.TypeParams, name.Name)
			}
		}
	}

	generator.FileName = f.curFileName
	generator.GeneratorType = GeneratorTypeStructure
	fields, err := f.parseFields(structType.Fields.List)
	if err != nil {
		return fmt.Errorf("f.parseFields: %w", err)
	}

	generator.Fields = fields
	return nil
}

func (f *generatorFactory) parseFields(fieldList []*ast.Field) ([]Field, error) {
	var fields []Field
	for _, field := range fieldList {
		if len(field.Names) == 0 {
			continue
		}

		typeStr, err := parseNode(f.curFset, field.Type)
		if err != nil {
			return nil, fmt.Errorf("parseNode: %w", err)
		}

		for _, name := range field.Names {
			debug.Printf("parse field name: %s, type: %s", name, typeStr)
			fields = append(fields, Field{name.Name, typeStr})
		}
	}
	return fields, nil
}

func (f *generatorFactory) inspectFunctionDeclaration(decl *ast.FuncDecl) error {
	if decl.Recv == nil {
		return nil
	}

	if len(decl.Recv.List) != 1 {
		return fmt.Errorf("expected one receiver, got %d", len(decl.Recv.List))
	}
	receiver := decl.Recv.List[0]

	t := receiver.Type
	// handler pointer
	if starExpr, ok := t.(*ast.StarExpr); ok {
		debug.Printf("inspect t as *ast.StarExpr\n")
		t = starExpr.X
	}
	// handler type parameters
	if indexExpr, ok := t.(*ast.IndexExpr); ok {
		debug.Printf("inspect t as *ast.IndexExpr\n")
		t = indexExpr.X
	}
	if indexListExpr, ok := t.(*ast.IndexListExpr); ok {
		debug.Printf("inspect t as *ast.IndexListExpr\n")
		t = indexListExpr.X
	}
	ident, ok := t.(*ast.Ident)
	if !ok {
		return fmt.Errorf("unexpected receiver type: %T", t)
	}
	recvTypeName := ident.Name

	generator, ok := f.generators[recvTypeName]
	if !ok {
		return nil
	}

	if names := receiver.Names; len(names) > 0 {
		generator.ReceiverName = names[0].Name
	}
	generator.Methods[decl.Name.Name] = struct{}{}

	return nil
}

// help functions

func parseTypeArguments(expr ast.Expr) []string {
	args := make([]string, 0)
	if expr, ok := expr.(*ast.IndexExpr); ok {
		args = append(args, fmt.Sprint(expr.Index))
	}
	if expr, ok := expr.(*ast.IndexListExpr); ok {
		for _, index := range expr.Indices {
			args = append(args, fmt.Sprint(index))
		}
	}
	return args
}

func parseNode(fset *token.FileSet, node any) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return "", fmt.Errorf("printer.Fprint: %w", err)
	}
	return buf.String(), nil
}

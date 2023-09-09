package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type options struct {
	getter     bool
	setter     bool
	pureGetter bool
	prefix     string
	includes   map[string]struct{}
	excludes   map[string]struct{}
}

type optionsFn func(*options)

func WithGetter(v bool) optionsFn {
	return func(o *options) {
		o.getter = v
	}
}

func WithSetter(v bool) optionsFn {
	return func(o *options) {
		o.setter = v
	}
}

func WithPureGetter(v bool) optionsFn {
	return func(o *options) {
		o.pureGetter = v
	}
}

func WithPrefix(p string) optionsFn {
	return func(o *options) {
		o.prefix = p
	}
}

func WithIncludes(includes []string) optionsFn {
	return func(o *options) {
		if len(includes) == 0 {
			o.includes = nil
			return
		}

		o.includes = make(map[string]struct{}, len(includes))
		for _, include := range includes {
			o.includes[include] = struct{}{}
		}
	}
}

func WithExcludes(excludes []string) optionsFn {
	return func(o *options) {
		if len(excludes) == 0 {
			o.excludes = nil
			return
		}

		o.excludes = make(map[string]struct{}, len(excludes))
		for _, exclude := range excludes {
			o.excludes[exclude] = struct{}{}
		}
	}
}

type Generator struct {
	Name          string
	Dir           string
	Pkg           string
	Type          string
	TypeParams    []string
	TypeArguments []string
	ReceiverName  string
	Fields        []Field
	Methods       map[string]struct{}
	FileName      string
	GeneratorType GeneratorType
	Imports       []string

	opts *options
}

type GeneratorType int

//go:generate go run golang.org/x/tools/cmd/stringer -type GeneratorType
const (
	GeneratorTypeUnknown GeneratorType = iota
	GeneratorTypeVariable
	GeneratorTypeStructure
	GeneratorTypeField
)

type Field struct {
	Name, Type string
}

func (g *Generator) InspectImports(unnamedImports []string, namedImports map[string]string) error {
	seenPkgs := make(map[string]struct{})
	if pkgName := getPackageNameFromType(g.Type); pkgName != "" {
		if err := g.inspectImport(pkgName, unnamedImports, namedImports, seenPkgs); err != nil {
			return fmt.Errorf("g.inspectImport %s %v %v: %w", pkgName, unnamedImports, namedImports, err)
		}
	}

	for _, typeArg := range g.TypeArguments {
		pkgName := getPackageNameFromType(typeArg)
		if pkgName == "" {
			continue
		}
		if err := g.inspectImport(pkgName, unnamedImports, namedImports, seenPkgs); err != nil {
			return fmt.Errorf("g.inspectImport %s %v %v: %w", pkgName, unnamedImports, namedImports, err)
		}
	}

	for _, field := range g.Fields {
		pkgName := getPackageNameFromType(field.Type)
		if pkgName == "" {
			continue
		}
		if err := g.inspectImport(pkgName, unnamedImports, namedImports, seenPkgs); err != nil {
			return fmt.Errorf("g.inspectImport %s %v %v: %w", pkgName, unnamedImports, namedImports, err)
		}
	}
	return nil
}

func (g *Generator) inspectImport(pkgName string, unnamedImports []string, namedImports map[string]string, seenPkgs map[string]struct{}) error {
	if _, ok := seenPkgs[pkgName]; ok {
		return nil
	}
	seenPkgs[pkgName] = struct{}{}

	if ipt, ok := namedImports[pkgName]; ok {
		g.Imports = append(g.Imports, ipt)
		return nil
	}

	for i, ctn := 0, true; ctn; i++ {
		ctn = false
		for _, ipt := range unnamedImports {
			trimmedIPT := strings.Trim(ipt, "\"")
			names := strings.Split(trimmedIPT, "/")
			p := len(names) - 1 - i
			if p < 0 {
				continue
			}

			ctn = true
			if pkgName == names[p] {
				g.Imports = append(g.Imports, ipt)
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find package %s", pkgName)
}

func (g *Generator) Generate(optsFn ...optionsFn) error {
	g.opts = &options{}
	for _, o := range optsFn {
		o(g.opts)
	}

	debug.Printf("Generator.Name %s", g.Name)
	debug.Printf("Generator.Dir %s", g.Dir)
	debug.Printf("Generator.Pkg %s", g.Pkg)
	debug.Printf("Generator.Type %s", g.Type)
	debug.Printf("Generator.TypeParams %s", g.TypeParams)
	debug.Printf("Generator.TypeArguments %s", g.TypeArguments)
	debug.Printf("Generator.ReceiverName %s", g.ReceiverName)
	debug.Printf("Generator.Fields %s", g.Fields)
	debug.Printf("Generator.Methods %s", g.Methods)
	debug.Printf("Generator.FileName %s", g.FileName)
	debug.Printf("Generator.GeneratorType %s", g.GeneratorType)
	debug.Printf("Generator.Imports %s", g.Imports)

	// TODO replace the method route by interface.
	var err error
	switch g.GeneratorType {
	case GeneratorTypeVariable:
		err = g.WriteVarAccessor()
	case GeneratorTypeStructure:
		err = g.WriteStructAccessor()
	case GeneratorTypeField:
		err = g.WriteFieldAccessor()
	default:
		err = fmt.Errorf("unknown generator type %s", g.GeneratorType)
	}
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) WriteVarAccessor() error {
	cl := append(g.getPackageCodeLines(), g.getVarCodeLines()...)
	return g.writeFile(cl)
}

func (g *Generator) WriteStructAccessor() error {
	cl := append(g.getPackageCodeLines(), g.getStructCodeLines()...)
	return g.writeFile(cl)
}

func (g *Generator) WriteFieldAccessor() error {
	cl := append(g.getPackageCodeLines(), g.getFieldCodeLines()...)
	return g.writeFile(cl)
}

type codeLines []struct {
	format string
	a      []interface{}
}

func (c codeLines) Append(format string, a ...interface{}) codeLines {
	return append(c, struct {
		format string
		a      []interface{}
	}{format, a})
}

func (g *Generator) getPackageCodeLines() (cl codeLines) {
	cl = cl.Append("// Code generated by \"goaccessor %s\". DO NOT EDIT.", strings.Join(os.Args[1:], " "))
	cl = cl.Append("")
	cl = cl.Append("package %s", g.Pkg)

	if len(g.Imports) == 1 {
		cl = cl.Append("")
		cl = cl.Append("import %s", g.Imports[0])
		cl = cl.Append("")
	}
	if len(g.Imports) > 1 {
		cl = cl.Append("")
		cl = cl.Append("import (")
		for _, ipt := range g.Imports {
			cl = cl.Append("        %s", ipt)
		}
		cl = cl.Append(")")
		cl = cl.Append("")
	}

	return
}

func (g *Generator) getVarCodeLines() (cl codeLines) {
	getMethodName := concat(g.GetPrefix(), g.opts.prefix, g.Name)
	if g.opts.getter && getMethodName != g.Name {
		cl = cl.Append("")
		cl = cl.Append("func %s() %s {", getMethodName, g.Type)
		cl = cl.Append("        return %s", g.Name)
		cl = cl.Append("}")
	} else if g.opts.getter {
		cl = cl.Append("")
		cl = cl.Append("// %s already exists", getMethodName)
	}

	if g.opts.setter {
		cl = cl.Append("")
		cl = cl.Append("func %s(%s %s) {", concat("set", g.opts.prefix, g.Name), g.newValueName(), g.Type)
		cl = cl.Append("        %s = %s", g.Name, g.newValueName())
		cl = cl.Append("}")
	}
	return
}

func (g *Generator) getStructCodeLines() (cl codeLines) {
	for _, field := range g.Fields {
		fieldName, fieldType := field.Name, field.Type
		if includes := g.opts.includes; len(includes) > 0 {
			if _, ok := includes[fieldName]; !ok {
				continue
			}
		}

		if excludes := g.opts.excludes; len(excludes) > 0 {
			if _, ok := excludes[fieldName]; ok {
				continue
			}
		}

		getMethodName := concat(g.GetPrefix(), g.opts.prefix, fieldName)
		if g.opts.getter && !g.IsNameExist(getMethodName) {
			cl = cl.Append("")
			cl = cl.Append("func (%s *%s) %s() %s {", g.getReceiverName(), g.getReceiverType(), getMethodName, fieldType)
			cl = cl.Append("        return %s.%s", g.getReceiverName(), fieldName)
			cl = cl.Append("}")
		} else if g.opts.getter {
			cl = cl.Append("")
			cl = cl.Append("// %s already exists", getMethodName)
		}

		setMethodName := concat("set", g.opts.prefix, fieldName)
		if g.opts.setter && !g.IsNameExist(setMethodName) {
			cl = cl.Append("")
			cl = cl.Append("func (%s *%s) %s(%s %s) {", g.getReceiverName(), g.getReceiverType(), setMethodName, g.newValueName(), fieldType)
			cl = cl.Append("        %s.%s = %s", g.getReceiverName(), fieldName, g.newValueName())
			cl = cl.Append("}")
		} else if g.opts.setter {
			cl = cl.Append("")
			cl = cl.Append("// %s already exists", setMethodName)
		}
	}
	return
}

func (g *Generator) getFieldCodeLines() (cl codeLines) {
	for _, field := range g.Fields {
		fieldName, fieldType := field.Name, field.Type
		if includes := g.opts.includes; len(includes) > 0 {
			if _, ok := includes[fieldName]; !ok {
				continue
			}
		}

		if excludes := g.opts.excludes; len(excludes) > 0 {
			if _, ok := excludes[fieldName]; ok {
				continue
			}
		}

		// check type arguments
		fieldType = g.fillTypeArguments(fieldType)

		getMethodName := concat(g.GetPrefix(), g.opts.prefix, fieldName)
		if g.opts.getter && getMethodName != g.Name && getMethodName != g.Type {
			cl = cl.Append("")
			cl = cl.Append("func %s() %s {", getMethodName, fieldType)
			cl = cl.Append("        return %s.%s", g.Name, fieldName)
			cl = cl.Append("}")
		} else if g.opts.getter {
			cl = cl.Append("")
			cl = cl.Append("// %s already exists", getMethodName)
		}

		setMethodName := concat("set", g.opts.prefix, fieldName)
		if g.opts.setter {
			cl = cl.Append("")
			cl = cl.Append("func %s(%s %s) {", setMethodName, g.newValueName(), fieldType)
			cl = cl.Append("        %s.%s = %s", g.Name, fieldName, g.newValueName())
			cl = cl.Append("}")
		}
	}
	return
}

func (g *Generator) getReceiverName() string {
	if g.ReceiverName != "" {
		return g.ReceiverName
	}
	g.ReceiverName = strings.ToLower(g.Type[:1])
	return g.ReceiverName
}

func (g *Generator) newValueName() string {
	if g.Name != "v" && g.getReceiverName() != "v" {
		return "v"
	}
	return "val"
}

func (g *Generator) getReceiverType() string {
	if len(g.TypeParams) == 0 {
		return g.Name
	}
	return g.Name + "[" + strings.Join(g.TypeParams, ", ") + "]"
}

func (g *Generator) fillTypeArguments(t string) string {
	for i, param := range g.TypeParams {
		t = fillTypeArguments(t, param, g.TypeArguments[i])
	}
	return t
}

func (g *Generator) writeFile(cl codeLines) error {
	f, err := os.Create(g.FilePath())
	if err != nil {
		return fmt.Errorf("os.Create(%s): %w", g.FilePath(), err)
	}
	defer f.Close()

	for _, line := range cl {
		_, err := fmt.Fprintf(f, line.format+"\n", line.a...)
		if err != nil {
			return fmt.Errorf("fmt.Fprintf(%s, %v): %w", line.format, line.a, err)
		}
	}
	return nil
}

func (g *Generator) FilePath() string {
	return filepath.Join(g.Dir, g.FileName+strings.ToLower(g.Name)+"_goaccessor.go")
}

func (g *Generator) GetPrefix() string {
	if g.opts.pureGetter {
		return ""
	}
	return "get"
}

func (g *Generator) IsNameExist(name string) bool {
	if _, ok := g.Methods[name]; ok {
		return true
	}
	for _, field := range g.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

// helper functions

func concat(strs ...string) (s string) {
	for _, str := range strs {
		if str == "" {
			continue
		}
		s += upper(str)
	}
	return
}

func upper(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}

func fillTypeArguments(t, param, arg string) string {
	re := regexp.MustCompile(`(^|[^_0-9\p{L}])` + regexp.QuoteMeta(param) + `($|[^_0-9\p{L}])`)
	return re.ReplaceAllStringFunc(t, func(s string) string {
		if strings.HasPrefix(s, param) {
			return arg + s[len(param):]
		} else if strings.HasSuffix(s, param) {
			return s[:len(s)-len(param)] + arg
		} else {
			return strings.ReplaceAll(s, param, arg)
		}
	})
}

func getPackageNameFromType(typeName string) string {
	if !strings.Contains(typeName, ".") {
		return ""
	}

	// hacky way to handle anonymous types and generic types
	if strings.Contains(typeName, "{") || strings.Contains(typeName, "[") {
		return ""
	}

	pkgName := strings.Split(typeName, ".")[0]

	// handle pointer type
	pkgName = strings.TrimLeft(pkgName, "*")
	return pkgName
}

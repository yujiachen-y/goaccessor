// Goaccessor provides a Go tool designed to automate the generation of getter and setter boilerplate
// code for your data types and variables.
//
// Usage:
//
//	go install github.com/yujiachen-y/goaccessor@latest
//
// After installing goaccessor, you can use it in the command line interface (CLI) or with //go:generate directives.
//
// Examples:
//
// Consider the file book.go:
//
//	package main
//
//	//go:generate goaccessor --target Book --getter --setter
//	type Book struct {
//	    Title  string
//	    Author string
//	}
//
// When we run go generate, it creates a new file book_goaccessor.go:
//
//	package main
//
//	func (b *Book) GetTitle() string {
//	    return b.Title
//	}
//
//	func (b *Book) SetTitle(title string) {
//	    b.Title = title
//	}
//
//	func (b *Book) GetAuthor() string {
//	    return b.Author
//	}
//
//	func (b *Book) SetAuthor(author string) {
//	    b.Author = author
//	}
//
// goaccessor isn't just for struct types; it can also handle top-level constants and variables. For instance:
//
//	//go:generate goaccessor --target books --getter --setter
//	var books = map[string]*Book {
//	    ...
//	}
//
// After executing go generate, we get:
//
//	func GetBooks() map[string]*Book {
//	    return books
//	}
//
//	func SetBooks(newBooks map[string]*Book) {
//	    books = newBooks
//	}
//
// In certain cases, you might want to export specific fields of a top-level variable with a prefix:
//
//	//go:generate goaccessor --target bestSellingBook --field --getter --include Author --prefix BestSelling
//	var bestSellingBook = &Book{ ... }
//
// This directive will generate:
//
//	func GetBestSellingAuthor() string {
//	    return bestSellingBook.Author
//	}
//
// Options:
//
// Here are the available options for goaccessor:
//
//	--target | -t: Specify the target to be handled.
//	--getter | -g: Generate getter for the target.
//	--setter | -s: Generate setter for the target.
//	--accessor | -a: Generate both getter and setter for the target.
//	--prefix | -p: Add a prefix to the generated methods/functions.
//	--field | -f: Apply the flag (getter, setter, accessor) to each field of the target (only applicable for struct type variables).
//	--include | -i: Generate methods only for the specified fields (fields should be comma-separated).
//	--exclude | -e: Exclude specified fields from method generation (fields should be comma-separated).
//
// Dependency Management:
//
// If you do not want to install goaccessor and want to use it as a dependency for your project, follow these steps:
//
// 1. Go to your project directory and add the goaccessor dependency via go mod:
//
//	go get github.com/yujiachen-y/goaccessor@latest
//
// 2. Create a new file named tools.go (or any other name you like) with the following code:
//
//	//go:build tools
//
//	package main
//
//	import (
//	    _ "github.com/yujiachen-y/goaccessor"
//	)
//
// 3. If you want to use the goaccessor CLI command, use go run github.com/yujiachen-y/goaccessor@latest instead.
//
// 4. If you use go:generate directives to generate code, change the directives like this:
//
//	//go:generate go run github.com/yujiachen-y/goaccessor@latest -target book
//	var book Book
//
// 5. Note: even though goaccessor is a dependency of your project, it will not—and should not—be a part of your project's build result.
// We use a build flag in tools.go to ensure goaccessor is ignored during build.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var debug *log.Logger

func setupLogger() {
	if os.Getenv("DEBUG") != "" {
		debug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(ioutil.Discard, "", 0)
	}
	log.SetFlags(0)
	log.SetPrefix("goaccessor: ")
}

var (
	flagTargets    []string
	flagGetter     bool
	flagSetter     bool
	flagPureGetter bool
	flagField      bool
	flagPrefix     string
	flagIncludes   []string
	flagExcludes   []string
	argDir         string
)

func parseFlags() {
	t := flag.String("t", "", "")
	target := flag.String("target", "", "")
	g := flag.Bool("g", false, "")
	getter := flag.Bool("getter", false, "")
	s := flag.Bool("s", false, "")
	setter := flag.Bool("setter", false, "")
	a := flag.Bool("a", false, "")
	accessor := flag.Bool("accessor", false, "")
	pg := flag.Bool("pg", false, "")
	pureGetter := flag.Bool("pure-getter", false, "")
	f := flag.Bool("f", false, "")
	field := flag.Bool("field", false, "")
	p := flag.String("p", "", "")
	prefix := flag.String("prefix", "", "")
	i := flag.String("i", "", "")
	include := flag.String("include", "", "")
	e := flag.String("e", "", "")
	exclude := flag.String("exclude", "", "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of goaccessor:\n")
		fmt.Fprintf(os.Stderr, "\t--target -t string\n")
		fmt.Fprintf(os.Stderr, "\t\tSpecify the target to be handled.\n")
		fmt.Fprintf(os.Stderr, "\t--getter -g getter\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate `getter` for the target.\n")
		fmt.Fprintf(os.Stderr, "\t--setter -s getter\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate `setter` for the target.\n")
		fmt.Fprintf(os.Stderr, "\t--accessor -a getter\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate `accessor` for the target.\n")
		fmt.Fprintf(os.Stderr, "\t--pure-getter -pg getter\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate `getter` without 'Get' prefix for the target.\n")
		fmt.Fprintf(os.Stderr, "\t--field -f getter\n")
		fmt.Fprintf(os.Stderr, "\t\tApply the command (`getter`, `setter`, `accessor`) to each field of the target (only works for struct type variables).\n")
		fmt.Fprintf(os.Stderr, "\t--prefix -p string\n")
		fmt.Fprintf(os.Stderr, "\t\tAdd a prefix to the generated methods/functions.\n")
		fmt.Fprintf(os.Stderr, "\t--include -i string\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate methods only for the specified fields (fields should be comma-separated).\n")
		fmt.Fprintf(os.Stderr, "\t--exclude -e string\n")
		fmt.Fprintf(os.Stderr, "\t\tExclude specified fields from method generation (fields should be comma-separated).\n")
		fmt.Fprintf(os.Stderr, "For more information, see:\n")
		fmt.Fprintf(os.Stderr, "\thttps://www.github.com/yujiachen-y/goaccessor\n")
	}

	flag.Parse()

	if len(*t) != 0 {
		flagTargets = strings.Split(*t, ",")
	} else if len(*target) != 0 {
		flagTargets = strings.Split(*target, ",")
	}
	if len(flagTargets) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	flagGetter = *g || *getter
	flagSetter = *s || *setter
	if *a || *accessor {
		flagGetter = true
		flagSetter = true
	}
	if *pg || *pureGetter {
		flagGetter = true
		flagPureGetter = true
	}
	if !flagGetter && !flagSetter {
		flag.Usage()
		os.Exit(2)
	}

	flagField = *f || *field

	if *p != "" {
		flagPrefix = *p
	} else if *prefix != "" {
		flagPrefix = *prefix
	}

	if len(*i) != 0 {
		flagIncludes = strings.Split(*i, ",")
	} else if len(*include) != 0 {
		flagIncludes = strings.Split(*include, ",")
	}

	if len(*e) != 0 {
		flagExcludes = strings.Split(*e, ",")
	} else if len(*exclude) != 0 {
		flagExcludes = strings.Split(*exclude, ",")
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	path := args[0]
	pathInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if pathInfo.IsDir() {
		argDir = path
	} else {
		argDir = filepath.Dir(path)
	}
}

func main() {
	setupLogger()
	parseFlags()

	debug.Printf("Received arguments:\n")
	debug.Printf("\t\tflagTargets %s\n", flagTargets)
	debug.Printf("\t\tflagGetter %t\n", flagGetter)
	debug.Printf("\t\tflagSetter %t\n", flagSetter)
	debug.Printf("\t\tflagPureGetter %t\n", flagPureGetter)
	debug.Printf("\t\tflagField %t\n", flagField)
	debug.Printf("\t\tflagPrefix %s\n", flagPrefix)
	debug.Printf("\t\tflagIncludes %s\n", flagIncludes)
	debug.Printf("\t\tflagExcludes %s\n", flagExcludes)
	debug.Printf("\t\targDir %s\n", argDir)

	generators, err := NewGenerators(flagTargets, argDir, flagField)
	if err != nil {
		log.Fatalf("Failed to create generators, error: %s", err.Error())
	}

	for _, generator := range generators {
		log.Printf("generate %s ...\n", generator.Name)
		err := generator.Generate(
			WithGetter(flagGetter),
			WithSetter(flagSetter),
			WithPureGetter(flagPureGetter),
			WithPrefix(flagPrefix),
			WithIncludes(flagIncludes),
			WithExcludes(flagExcludes),
		)
		if err != nil {
			log.Fatalf("Failed to generate, error: %s", err.Error())
		}
	}
}

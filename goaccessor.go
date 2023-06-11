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

func init() {
	if os.Getenv("DEBUG") != "" {
		debug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(ioutil.Discard, "", 0)
	}
	log.SetFlags(0)
	log.SetPrefix("goaccessor: ")

	s := flag.String("s", "", "")
	symbol := flag.String("symbol", "", "")
	g := flag.Bool("g", false, "")
	getter := flag.Bool("getter", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of goaccessor:\n")
		fmt.Fprintf(os.Stderr, "\t--symbol -s string\n")
		fmt.Fprintf(os.Stderr, "\t\tSpecify the symbol to be handled.\n")
		fmt.Fprintf(os.Stderr, "\t--getter -g getter\n")
		fmt.Fprintf(os.Stderr, "\t\tGenerate `getter` for the symbol.\n")
		fmt.Fprintf(os.Stderr, "For more information, see:\n")
		fmt.Fprintf(os.Stderr, "\thttps://www.github.com/yjc567/goaccessor\n")
	}

	flag.Parse()

	if len(*s) != 0 {
		flagSymbols = strings.Split(*s, ",")
	} else if len(*symbol) != 0 {
		flagSymbols = strings.Split(*symbol, ",")
	}
	if len(flagSymbols) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	flagGetter = *g || *getter
	if !flagGetter {
		flag.Usage()
		os.Exit(2)
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

var (
	flagSymbols []string
	flagGetter  bool
	argDir      string
)

func main() {
	log.Printf("Received arguments:\n")
	log.Printf("\t\tflagSymbols %s\n", flagSymbols)
	log.Printf("\t\tflagGetter %t\n", flagGetter)
	log.Printf("\t\targDir %s\n", argDir)

	generators, err := NewGenerators(flagSymbols, argDir)
	if err != nil {
		log.Fatalf("Failed to create generators, error: %s", err.Error())
	}

	log.Println()
	for _, generator := range generators {
		log.Printf("generate %s ...\n", generator.Name)
		err := generator.Generate(WithGetter(flagGetter))
		if err != nil {
			log.Fatalf("Failed to generate, error: %s", err.Error())
		}
	}
}

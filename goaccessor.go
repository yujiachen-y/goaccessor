package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func init() {
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
		flagSymbol = *s
	} else if len(*symbol) != 0 {
		flagSymbol = *symbol
	}
	if len(flagSymbol) == 0 {
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
	flagSymbol string
	flagGetter bool
	argDir     string
)

func main() {
	log.Printf("flagSymbol %s\nflagGetter %t\nargDir %s\n", flagSymbol, flagGetter, argDir)
}

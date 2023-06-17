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

	t := flag.String("t", "", "")
	target := flag.String("target", "", "")
	g := flag.Bool("g", false, "")
	getter := flag.Bool("getter", false, "")
	s := flag.Bool("s", false, "")
	setter := flag.Bool("setter", false, "")
	a := flag.Bool("a", false, "")
	accessor := flag.Bool("accessor", false, "")

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
		fmt.Fprintf(os.Stderr, "For more information, see:\n")
		fmt.Fprintf(os.Stderr, "\thttps://www.github.com/yjc567/goaccessor\n")
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
	if !flagGetter && !flagSetter {
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
	flagTargets []string
	flagGetter  bool
	flagSetter  bool
	argDir      string
)

func main() {
	log.Printf("Received arguments:\n")
	log.Printf("\t\tflagTargets %s\n", flagTargets)
	log.Printf("\t\tflagGetter %t\n", flagGetter)
	log.Printf("\t\targDir %s\n", argDir)

	generators, err := NewGenerators(flagTargets, argDir)
	if err != nil {
		log.Fatalf("Failed to create generators, error: %s", err.Error())
	}

	log.Println()
	for _, generator := range generators {
		log.Printf("generate %s ...\n", generator.Name)
		err := generator.Generate(WithGetter(flagGetter), WithSetter(flagSetter))
		if err != nil {
			log.Fatalf("Failed to generate, error: %s", err.Error())
		}
	}
}

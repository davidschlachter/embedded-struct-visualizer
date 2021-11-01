// The embedded-struct-visualizer command builds a Graphviz DOT file
// representing the tree of embedded structs in a Go project
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Struct struct {
	Name     string
	Package  string
	FilePath string
	Embeds   map[string]bool
}

var structsList []Struct
var verbose *bool

func main() {
	var (
		searchPath = "./"
		outputFile = os.Stdout
		err        error
		flags      flag.FlagSet
	)

	outputPath := flags.String("out", "", "write to file instead of stdout")
	verbose = flags.Bool("v", false, "verbose logging")
	flags.Usage = help
	flags.Parse(os.Args[1:])

	if len(flags.Args()) == 1 {
		searchPath = flags.Arg(0)
	}
	if *outputPath != "" {
		outputFile, err = os.OpenFile(*outputPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			fmt.Printf("Error writing output file: %v", err)
			return
		}
	}

	_ = filepath.WalkDir(searchPath, findGoFiles)

	graph := buildDOTFile()

	writer := bufio.NewWriter(outputFile)
	_, err = writer.WriteString(graph)
	if err != nil {
		fmt.Printf("Error writing output file: %v", err)
		return
	}
	writer.Flush()
}

func help() {
	fmt.Printf("Usage: %s [OPTIONS] DirToScan\n", os.Args[0])
	fmt.Printf("If the directory to scan is not provided, it defaults to './'\n")
	fmt.Printf("OPTIONS:\n")
	fmt.Printf("  -out <file>  path to output file (default: write to stdout)\n")
	fmt.Printf("  -v           verbose logging\n")
	os.Exit(1)
}

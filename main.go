// The embedded-struct-visualizer command builds a Graphviz DOT file
// representing the tree of embedded structs in a Go project
package main

import (
	"fmt"
	"path/filepath"
)

const searchPath = "./"

type Struct struct {
	Name     string
	Package  string
	FilePath string
	Embeds   []string
}

var structsList []Struct

func main() {
	_ = filepath.WalkDir(searchPath, walk)
	graph := buildDOTFile()
	fmt.Println(graph)
}

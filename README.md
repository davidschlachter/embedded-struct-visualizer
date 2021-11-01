# Visualize the hierarchy of embedded structs

This tool scans a directory of Go source code files to create a visualization of struct embedding in the project. This can be useful for navigating the data structures of complex projects, or identifying dependencies on particular structures.

# Example

Given the following input file:

```go
package main

import (
	"time"

	"github.com/davidschlachter/embedded-struct-visualizer/user"
)

type A struct {
	B
	C map[string]D
}

type B struct {
	E, F  string
	G     user.Status
	Timer H
}

type D struct {
	I uint64
}

type H struct {
	Timer time.Ticker
	J     chan D
}

```

The following figure would be generated:

[IMAGE PLACEHOLDER]

# Usage

Install:

`$ go install github.com/davidschlachter/embedded-struct-visualizer@latest`

Or,

```shell
$ git clone https://github.com/davidschlachter/embedded-struct-visualizer
$ cd embedded-struct-visualizer
$ go install github.com/davidschlachter/embedded-struct-visualizer
```

Options:

```
$ embedded-struct-visualizer -h
Usage: embedded-struct-visualizer [OPTIONS] DirToScan                                                        
If the directory to scan is not provided, it defaults to './'
OPTIONS:
  -out <file>  path to output file (default: write to stdout)
  -v           verbose logging
```

To open a generated DOT file, you could use [Graphviz](https://graphviz.org/download/) or [xdot](https://github.com/jrfonseca/xdot.py).
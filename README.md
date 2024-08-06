# Visualize a hierarchy of embedded Go structs

This tool scans a directory of Go source code files to create a visualization of struct embedding in the project. This can be useful for navigating the data structures of complex projects, or identifying dependencies on particular structures.

# Example

This figure would be generated from the input file below:

![visualizer screenshot](https://user-images.githubusercontent.com/8271268/139614622-01693a08-4aff-4c47-81e2-764a7e4ca0ef.png)

```go
package main

import (
	"time"
	"domain.tld/user"
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
  -rankdir <direction> graphs direction (default: LR, enum: TB,LR,BT,RL)
  -exclude-pkg path to exclude pkg config file, format(default: empty):
     eg: exclude gopkg.in/guregu/null.v3 and models/MyStuct
       prefix:null.
       models.MyStuct
  -v           verbose logging
```

To open a generated DOT file, you could use [Graphviz](https://graphviz.org/download/) or [xdot](https://github.com/jrfonseca/xdot.py).

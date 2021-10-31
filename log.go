package main

import "fmt"

func log(s string) {
	if verbose != nil && *verbose {
		fmt.Printf("%s\n", s)
	}
}

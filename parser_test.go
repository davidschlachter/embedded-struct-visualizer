package main

import (
	"strings"
	"testing"
)

func TestParseStructs(t *testing.T) {
	readFile(strings.NewReader(testString), "./testString")
	if len(structsList) != len(expectedStructs) {
		t.Fatalf("len(structsList) != len(expectedStructs)")
	}
	// n^2 complexity, improve this
	for _, s := range structsList {
		found := false
		for _, e := range expectedStructs {
			if s.Name == e.Name && s.Package == e.Package && s.FilePath == e.FilePath && len(s.Embeds) == len(e.Embeds) {
				found = true
			}
		}
		if !found {
			t.Fatalf("Could not find expected Struct %s", s.Name)
		}
	}
}

const testString = `
package main

type A struct {
	Int        int
	Bool       bool
	Map        map[string]bool
	ChanString chan string
	ChanStruct chan M
}

type M struct {
	Q int
}

type Z struct {
	A
	SliceOfStructs []M
}
`

var expectedStructs = []Struct{
	{Name: "A", Package: "main", FilePath: "./testString", Embeds: []string{"M"}},
	{Name: "M", Package: "main", FilePath: "./testString", Embeds: []string{}},
	{Name: "Z", Package: "main", FilePath: "./testString", Embeds: []string{"A", "M"}},
}

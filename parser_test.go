package main

import (
	"strings"
	"testing"
)

const testString1 = `
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

const testString2 = `
package cmd

type A struct {
	M map[string]map[string]*N
	I map[string][]string
}

type N struct {
	Name                                     string
	node                                     q.Node
	Tags                                     map[string][]string
	W, X, Y, Z string
}

type B struct {
	Item1, Item2    N
	A map[string]*A
	F
}

type F struct {
	G   map[string]*time.Time
	H map[string]string
	I  string
}
`

func TestParseStructs(t *testing.T) {
	readFile(strings.NewReader(testString1), "./testString1")
	readFile(strings.NewReader(testString2), "./testString2")
	if len(structsList) != len(expectedStructs) {
		t.Fatalf("len(structsList) != len(expectedStructs)")
	}
	// n^2 complexity, improve this
	// TODO: actually check values in Embeds
	for _, s := range structsList {
		found := false
		for _, e := range expectedStructs {
			if s.Name == e.Name && s.Package == e.Package && s.FilePath == e.FilePath && len(s.Embeds) == len(e.Embeds) {
				found = true
			}
		}
		if !found {
			t.Fatalf("Could not find expected Struct: %+v", s)
		}
	}
}

var expectedStructs = []Struct{
	{Name: "A", Package: "main", FilePath: "./testString1", Embeds: []string{"M"}},
	{Name: "M", Package: "main", FilePath: "./testString1", Embeds: []string{}},
	{Name: "Z", Package: "main", FilePath: "./testString1", Embeds: []string{"A", "M"}},
	{Name: "A", Package: "cmd", FilePath: "./testString2", Embeds: []string{"N"}},
	{Name: "N", Package: "cmd", FilePath: "./testString2", Embeds: []string{"q.Node"}},
	{Name: "B", Package: "cmd", FilePath: "./testString2", Embeds: []string{"N", "A", "F"}},
	{Name: "F", Package: "cmd", FilePath: "./testString2", Embeds: []string{"time.Time"}},
}

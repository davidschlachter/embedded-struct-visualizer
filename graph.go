package main

import (
	"fmt"
	"strings"
)

func buildDOTFile(rankdir string) string {
	g := []byte{}

	g = append(g, fmt.Sprintf(`
digraph {
rankdir="%s"
`, rankdir)...)
	for _, s := range structsList {
		if len(s.Embeds) == 0 {
			continue
		}
		g = append(g, getFullStructName(s.Name, s.Package)+" -> { "...)
		for e := range s.Embeds {
			if IsExcludePkg(e, s.Package) {
				continue
			}
			g = append(g, getFullStructName(e, s.Package)+" "...)
		}
		g = append(g, "};\n"...)
	}
	g = append(g, "}"...)
	return string(g)
}

func getFullStructName(s string, pkg string) string {
	if strings.Contains(s, ".") { // struct name includes package name
		return "\"" + s + "\""
	}

	return "\"" + pkg + "." + s + "\""
}

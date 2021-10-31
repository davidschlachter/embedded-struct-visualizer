package main

import "strings"

func buildDOTFile() string {
	g := []byte{}
	g = append(g, "digraph {\n"...)
	for _, s := range structsList {
		if len(s.Embeds) == 0 {
			continue
		}
		g = append(g, getFullStructName(s.Name, s.Package)+" -> { "...)
		for _, e := range s.Embeds {
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

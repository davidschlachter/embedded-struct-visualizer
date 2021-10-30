package main

import "strings"

func buildDOTFile() string {
	g := []byte{}
	for _, s := range structsList {
		if len(s.Embeds) == 0 {
			continue
		}
		g = append(g, s.Package+"."+s.Name+" -> { "...)
		for _, e := range s.Embeds {
			g = append(g, getFullStructName(e, s.Package)+" "...)
		}
		g = append(g, "}\n"...)
	}
	return string(g)
}

func getFullStructName(s string, pkg string) string {
	if strings.Contains(s, ".") { // struct name includes package name
		return s
	}

	return pkg + "." + s
}

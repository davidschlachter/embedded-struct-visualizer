package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// recurse through a file tree, finding Go source files
func findGoFiles(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Skip this file if not Go source code
	if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open file %s: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	readFile(file, path)

	return nil
}

func readFile(file io.Reader, path string) {
	var (
		pkg  string
		line string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = removeComments(scanner, scanner.Text())

		// Identify the package name
		if pkg == "" && strings.HasPrefix(line, "package ") {
			pkg = strings.Replace(line, "package ", "", 1)
			continue
		}

		// Find all structs defined in this file
		if strings.HasPrefix(line, "type ") && strings.HasSuffix(line, " struct {") {
			log(line)
			parseStruct(scanner, line, pkg, path)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		os.Exit(1)
	}
}

func parseStruct(scanner *bufio.Scanner, line string, pkg string, path string) {
	name := getName(line)
	s := Struct{
		Name:     name,
		Package:  pkg,
		FilePath: path,
		Embeds:   make(map[string]bool),
	}
	// Read the fields of this struct
	for scanner.Scan() {
		structLine := removeComments(scanner, scanner.Text())
		log(structLine)
		if structLine == "}" {
			break
		}
		if structLine == "" {
			continue
		}
		if strings.HasSuffix(structLine, " struct {") {
			for scanner.Scan() {
				if strings.TrimSpace(scanner.Text()) == "}," {
					break
				}
			}
			continue
		}
		fields, hasStructs := getStructs(structLine)
		if hasStructs {
			for _, f := range fields {
				s.Embeds[f] = true
			}
		}
	}

	// Add to the list
	log(fmt.Sprintf("%+v\n", s))
	structsList = append(structsList, s)
}

func getStructs(s string) ([]string, bool) {
	tokens := strings.Fields(cleanTags(s))
	t := tokens[len(tokens)-1]
	if types[t] || strings.Contains(s, "func") {
		return []string{""}, false
	}
	if strings.HasPrefix(t, "map[") || strings.HasPrefix(t, "*map[") {
		return parseMap(t)
	}

	return []string{cleanPointers(t)}, true
}

func parseMap(s string) ([]string, bool) {
	var structs []string
	var hasStructs bool
	mapFields := strings.FieldsFunc(s, func(r rune) bool {
		return r == '[' || r == ']'
	})
	for _, f := range mapFields {
		if !types[f] {
			structs = append(structs, cleanPointers(f))
		}
	}
	if len(structs) > 0 {
		hasStructs = true
	}
	return structs, hasStructs
}

func getName(s string) string {
	s = strings.Replace(s, "type ", "", 1)
	s = strings.Replace(s, " struct {", "", 1)
	return s
}

func removeComments(scanner *bufio.Scanner, line string) string {
	s := []byte(line)
	j := 0
	for i, b := range s {
		if b == '/' && i+1 < len(s) && s[i+1] == '/' {
			break
		}
		if b == '/' && i+1 < len(s) && s[i+1] == '*' {
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "*/") {
					break
				}
			}
			break
		}
		j++
	}
	return strings.TrimSpace(string(s[:j]))
}

func cleanPointers(token string) string {
	s := []byte(token)
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '.' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func cleanTags(token string) string {
	s := []byte(token)
	j := 0
	for _, b := range s {
		if b == '`' {
			break
		}
		j++
	}
	return strings.TrimSpace(string(s[:j]))
}

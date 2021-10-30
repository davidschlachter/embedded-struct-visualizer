package main

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"strings"
)

func walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Skip this file if not Go source code
	if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
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
		line = scanner.Text()

		// Identify the package name
		if pkg == "" && strings.HasPrefix(line, "package ") {
			pkg = strings.Replace(line, "package ", "", 1)
			pkg = strings.TrimSpace(pkg)
		}

		// Find all structs defined in this file
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type ") && strings.HasSuffix(line, " struct {") {
			parseStruct(scanner, line, pkg, path)
		}

	}

	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
}

func parseStruct(scanner *bufio.Scanner, line string, pkg string, path string) {
	name := getName(line)
	s := Struct{
		Name:     name,
		Package:  pkg,
		FilePath: path,
	}
	// Read the fields of this struct
	for scanner.Scan() {
		structLine := strings.TrimSpace(scanner.Text())
		if structLine == "}" {
			break
		}
		if strings.HasPrefix(structLine, "//") || structLine == "" {
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
		if !isStruct(structLine) {
			continue
		}
		s.Embeds = append(s.Embeds, getStruct(structLine))
	}

	// Add to the list
	structsList = append(structsList, s)
}

func getName(s string) string {
	s = strings.Replace(s, "type ", "", 1)
	s = strings.Replace(s, " struct {", "", 1)
	return s
}

func isStruct(s string) bool {
	tokens := strings.Fields(cleanTags(s))
	if len(tokens) == 1 {
		return true
	}
	if len(tokens) >= 2 {
		t := cleanExtras(tokens[1])
		if types[t] || strings.HasPrefix(tokens[1], "map[") || strings.HasPrefix(tokens[1], "*map[") {
			if t == "chan" {
				return isStruct(tokens[1] + " " + tokens[2])
			}
			return false
		}
	}
	return true
}

func getStruct(s string) string {
	tokens := strings.Fields(cleanTags(s))
	if len(tokens) == 1 {
		return cleanExtras(tokens[0])
	}
	if len(tokens) == 3 {
		if tokens[1] == "chan" {
			return cleanExtras(tokens[2])
		}
	}
	if len(tokens) >= 2 {
		return cleanExtras(tokens[1])
	}
	return "INVALID: " + s
}

func cleanExtras(q string) string {
	s := []byte(q)
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '.' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func cleanTags(q string) string {
	s := []byte(q)
	slashCount := 0
	j := 0
	for _, b := range s {
		if b == '`' {
			break
		}
		if b == '/' {
			slashCount++
			if slashCount > 1 {
				break
			}
		}
		s[j] = b
		j++
	}
	return string(s[:j])
}

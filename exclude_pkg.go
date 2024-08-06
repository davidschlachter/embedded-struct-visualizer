package main

import (
	"os"
	"strings"
)

func parseExcludeConfig(filename string) error {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			excludePkg = append(excludePkg, line)
		}
	}
	return nil
}

func IsExcludePkg(structName, pack string) bool {
	var full = structName
	if !strings.Contains(structName, ".") {
		full = pack + "." + structName
	}
	for _, expkg := range excludePkg {
		if strings.HasPrefix(expkg, "prefix:") {
			ex := strings.TrimPrefix(expkg, "prefix:")
			if strings.HasPrefix(full, ex) {
				return true
			}
		} else {
			if expkg == full {
				return true
			}
		}
	}
	return false
}

//go:build ignore
// +build ignore

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	parts := []string{}

	err := filepath.Walk("../../internal/actions/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() ||
			filepath.Ext(path) != ".go" ||
			filepath.Base(path) == "actions.go" {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}

		out, err := format(string(raw))
		if err != nil {
			return fmt.Errorf("could not get headers from %s: %w", path, err)
		}

		parts = append(parts, out)

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("../../internal/cmd/actions-help.txt", []byte(strings.Join(parts, "\n\n")), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func format(content string) (string, error) {
	parts := strings.Split(content, "/*")
	if len(parts) < 2 {
		return "", errors.New("no header found")
	}

	parts = strings.Split(parts[1], "*/")
	if len(parts) < 2 {
		return "", errors.New("no header end found")
	}

	header := strings.Trim(parts[0], "\n")
	parts = strings.SplitN(header, "\n", 2)

	re := regexp.MustCompile(`Package (\w+) implements an \[actions.Runner\] that (.*)\.`)
	matches := re.FindStringSubmatch(parts[0])

	if len(matches) < 3 {
		return "", fmt.Errorf("header does not match the expected format")
	}

	outParts := []string{
		fmt.Sprintf("%s: %s", matches[1], matches[2]),
	}

	if len(parts) == 2 {
		tail := strings.Trim(parts[1], "\n")
		tail = indent(tail)
		outParts = append(outParts, tail)
	}

	return strings.Join(outParts, "\n"), nil
}

func indent(s string) string {
	return "  " + strings.ReplaceAll(s, "\n", "\n  ")
}

package main

import (
	"fmt"
	"os"

	"github.com/nobe4/gh-not/internal/cmd"
)

//go:generate go run ../gen-docs/main.go

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

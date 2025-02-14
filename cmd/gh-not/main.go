package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/nobe4/gh-not/internal/cmd"
)

//go:generate go run ../gen-docs/main.go

func main() {
	info, _ := debug.ReadBuildInfo()
	fmt.Println("Go version:", info.GoVersion)
	fmt.Println("App version:", info.Main.Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

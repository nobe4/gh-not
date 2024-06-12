package main

import (
	"fmt"
	"os"

	"github.com/nobe4/gh-not/internal/cmd"
)

func main() {
	fmt.Print("A")
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

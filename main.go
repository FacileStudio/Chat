// Package main is the entry point for the chat CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/FacileStudio/Chat/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

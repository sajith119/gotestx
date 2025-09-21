package main

import (
	"os"

	"github.com/entiqon/gotestx/internal"
)

func main() {
	code := gotestx.Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

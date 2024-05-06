package main

import (
	"fmt"
	"os"

	"github.com/dangermike/roku_toy/cmd"
)

func main() {
	if err := cmd.Cmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

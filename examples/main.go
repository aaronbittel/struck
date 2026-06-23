package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aaronbittel/struck"
)

type Options struct {
	Name    string  `long:"name" short:"n" help:"User name"`
	Age     uint64  `long:"age" help:"User age"`
	Verbose bool    `long:"verbose" short:"v" help:"Enable verbose output"`
	Input   float32 `arg:"input" help:"Input value"`
}

func main() {
	var opts Options

	parser := struck.NewParser(&opts)

	args := []string{"--name", "Bob", "--age", "42", "-v", "141.531"}

	if err := parser.Parse(args...); err != nil {
		switch {
		case errors.Is(err, struck.HelpRequested):
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if opts.Verbose {
		fmt.Printf("%+v\n", opts)
	}
}

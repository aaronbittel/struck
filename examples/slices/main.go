package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aaronbittel/struck"
)

type SlicesOptions struct {
	Names     []string `long:"name" short:"n" help:"people that should be greeted"`
	Verbosity []bool   `long:"verbose" short:"v" help:"set verbosity level (can be added multiple times)"`
}

func main() {
	opts := SlicesOptions{
		Names: []string{"Bob"},
	}

	parser := struck.NewParser(&opts)

	if err := parser.Parse(); err != nil {
		switch {
		case errors.Is(err, struck.HelpRequested):
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// TODO: currently each slice needs a value, so greet Bob, Alice and Charlie with
	// verbosity 3, you would need to call it like this:
	// go run ./examples/slices.go -v true --verbose true --verbose true --name alice -n charlie
	// here it would not matter if you put "true" or "false" because we just take the
	// length of the verbosity slice

	switch len(opts.Verbosity) {
	case 0:
		fmt.Println("hello")
	case 1:
		fmt.Println("hello and welcome:", strings.Join(opts.Names, ", "))
	case 2:
		fmt.Printf(
			"A warm and enthusiastic welcome to %s! We're delighted to have you here.\n",
			strings.Join(opts.Names, ", "))
	default:
		fmt.Printf(
			"🎉 A most magnificent welcome to %s! Thank you for gracing this humble program with your presence. 🎉\n",
			strings.Join(opts.Names, ", "),
		)
	}
}

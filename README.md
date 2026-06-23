# struck

A command-line argument parser for Go using struct tags. Inspired by
[go-flags](https://github.com/jessevdk/go-flags).

## Features

- Struct-tag based CLI definition
- Long (`--flag`) and short (`-f`) flags
- Positional arguments
- Automatic type parsing
- Built-in help generation
- `--help`, `-h` support

## Installation
```console
    go get github.com/aaronbittel/struck
```

## Quick Start

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

type Options struct {
	Name    string  `long:"name" short:"n" help:"User name"`
	Age     uint64  `long:"age" help:"User age"`
	Verbose bool    `long:"verbose" short:"v" help:"Enable verbose output"`
	Input   float32 `arg:"input" help:"Input value"`
}

func main() {
	var opts Options

	args := []string{"--name", "Bob", "--age", "42", "-v", "141.531"}

	if err := Parse(&opts, args...); err != nil {
		switch {
		case errors.Is(err, HelpRequested):
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if opts.Verbose {
		// Prints: {Name:Bob Age:42 Verbose:true Input:141.531}
		fmt.Printf("%+v\n", opts)
	}
}
```

## Behavior

### Flags

Flags are defined using:
- `long:"name"` -> `--name`
- `short:"n"` -> `-n`

- Flags (except `bool`) consume the next argument as their value
- Boolean flags are `true` if present, `false` if absent

### Positional Arguments

- Fields without long or short tags are treated as positional arguments.
- They are assigned in declaration order.
- Optional override:
    - `arg:"name"`
    - This changes the display name in help output.

### Help Flag

If any of the following is provided: `--help`, `-help`, `-h`
Then:
- Help output is printed
- `Parse()` returns `HelpRequested` error

## Struct Tags Reference

| Tag     | Meaning                               |
|---------|---------------------------------------|
| `long`  | Long flag name (`--name`)             |
| `short` | Short flag name (`-n`)                |
| `help`  | Help text                             |
| `arg`   | Positional argument name override     |

## Example

```go
type Options struct {
	Name    string  `long:"name" short:"n" help:"User name"`
	Age     uint64  `long:"age" help:"User age"`
	Verbose bool    `long:"verbose" short:"v" help:"Verbose output"`

	Path string `arg:"path" help:"Input file path"`
}

`./app --name bob --age 23 --verbose ./file.txt`

Result:

Options{
	Name: "bob",
	Age: 23,
	Verbose: true,
	Path: "./file.txt",
}
```

## Design Notes

- Non-exported fields are ignored
- Flags and positional arguments are mutually exclusive. Conflicting definitions (`long`
  and `arg` tag) are invalid and cause a panic

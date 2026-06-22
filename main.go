package main

import (
	"fmt"
	"os"
)

type Options struct {
	Name string  `long:"name" short:"n"`
	Age  uint64  `long:"age"`
	Pos1 float32 `arg:"pos"`
}

func main() {
	var opts Options

	args := []string{"-n", "bob", "--age", "23", "123.51"}

	if err := Parse(&opts, args...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(opts)
}

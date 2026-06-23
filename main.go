package main

import (
	"fmt"
	"os"
	"reflect"
)

type Options struct {
	Name               string  `long:"name" short:"n" help:"name input"`
	Age                uint64  `long:"age" help:"specificy the age" age:"no conflict"`
	Verbose            bool    `long:"verbose" short:"v" help:"if set output verbosely"`
	Dummy              string  `short:"asdf" long:""`
	Pos1               float32 `arg:"pos" help:"just some pos f32"`
	Pos2sdjfköajsödfas byte    `help:"this is some help text for pos2"`
	Pos3               uint8
}

func main() {
	var opts Options

	args := []string{"-n", "bob", "--age", "23", "123.51", "--verbose", "`", "123"}

	if err := Parse(&opts, args...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if opts.Verbose {
		fmt.Println(opts)
	}
}

func information(t reflect.Type) {
	fmt.Println(t)
	fmt.Println("name", t.Name())
	fmt.Println("kind", t.Kind())

	if hasElem(t) {
		fmt.Println("  elem  ", t.Elem())
		fmt.Println("    kind", t.Elem().Kind())
	}

	fmt.Println("=========================================")
}

func hasElem(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}

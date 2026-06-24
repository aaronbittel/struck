package struck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Parser struct {
	command *Command
	schema  reflect.Value
}

func NewParser(schema any, name ...string) *Parser {
	var cmdName string
	switch len(name) {
	case 0:
		if len(os.Args) > 0 {
			cmdName = filepath.Base(os.Args[0])
		}
	case 1:
		cmdName = name[0]
	default:
		panic("the parser can only have one name")
	}

	t := reflect.TypeOf(schema)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		panic("give type must be a pointer to a struct")
	}

	command := ConstructCommand(t.Elem())
	command.name = cmdName

	return &Parser{
		command: command,
		schema:  reflect.ValueOf(schema).Elem(),
	}
}

//lint:ignore ST1012 sentinel value used for control flow (like io.EOF); intentionally not prefixed with Err
var HelpRequested = errors.New("help requested")

func (p *Parser) ParseArgs(args []string) error {
	for _, arg := range args {
		if arg == "--help" || arg == "-help" || arg == "-h" {
			p.command.PrintHelp(os.Stdout)
			return HelpRequested
		}
	}
	return p.parseArgs(args)
}

func (p *Parser) Parse() error {
	return p.ParseArgs(os.Args[1:])
}

func (p *Parser) parseArgs(args []string) error {
	positionalArgIndex := 0

	i := 0
	for i < len(args) {
		currentArg := args[i]
		flag, ok := p.command.matchesFlag(currentArg)
		if ok {
			if p.ValueByIndex(flag.FieldIndex).Kind() == reflect.Bool {
				p.ValueByIndex(flag.FieldIndex).SetBool(true)
				i++
				continue
			}

			if p.ValueByIndex(flag.FieldIndex).Kind() == reflect.Slice && p.ValueByIndex(flag.FieldIndex).Type().Elem().Kind() == reflect.Bool {
				p.ValueByIndex(flag.FieldIndex).Set(reflect.Append(p.ValueByIndex(flag.FieldIndex), reflect.ValueOf(true)))
				i++
				continue
			}

			if !hasNext(args, i) {
				return fmt.Errorf("TODO: value for flag %q not provided", flag.Name())
			}

			err := SetValue(p.ValueByIndex(flag.FieldIndex), args[i+1])
			if err != nil {
				return err
			}
			i += 2
		} else {
			if positionalArgIndex >= len(p.command.positionals) {
				return fmt.Errorf("TODO: to manny positional arguments, arg=%q", currentArg)
			}

			positionalArg := p.command.positionals[positionalArgIndex]

			err := SetValue(p.ValueByIndex(positionalArg.FieldIndex), currentArg)
			if err != nil {
				return err
			}

			i += 1
			positionalArgIndex += 1
		}
	}

	if positionalArgIndex < len(p.command.positionals) {
		var sb strings.Builder
		sb.WriteString("missing values for the following positionals arguments:\n")
		for ; positionalArgIndex < len(p.command.positionals); positionalArgIndex++ {
			fmt.Fprintf(&sb, "  - %q\n", p.command.positionals[positionalArgIndex].Name)
		}
		return fmt.Errorf("%s", sb.String())
	}

	return nil
}

func (p *Parser) ValueByIndex(index []int) reflect.Value {
	return p.schema.FieldByIndex(index)
}

func hasNext(args []string, i int) bool {
	return i+1 < len(args)
}

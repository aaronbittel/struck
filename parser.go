package struck

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Parser struct {
	Name    string
	command *Command
	value   reflect.Value
}

func NewParser(schema any, name ...string) *Parser {
	var parserName string
	switch len(name) {
	case 0:
		if len(os.Args) > 0 {
			parserName = filepath.Base(os.Args[0])
		}
	case 1:
		parserName = name[0]
	default:
		panic("the parser can only have one name")
	}

	t := reflect.TypeOf(schema)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		panic("give type must be a pointer to a struct")
	}

	return &Parser{
		Name:    parserName,
		command: ConstructCommand(t.Elem()),
		value:   reflect.ValueOf(schema).Elem(),
	}
}

//lint:ignore ST1012 sentinel value used for control flow (like io.EOF); intentionally not prefixed with Err
var HelpRequested = errors.New("help requested")

func (p *Parser) ParseArgs(args []string) error {
	for _, arg := range args {
		if arg == "--help" || arg == "-help" || arg == "-h" {
			p.PrintHelp(os.Stdout)
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
			if flag.Type.Kind() == reflect.Bool {
				p.value.FieldByIndex(flag.FieldIndex).SetBool(true)
				i++
				continue
			}

			if !hasNext(args, i) {
				return fmt.Errorf("TODO: value for flag %q not provided", flag.Name())
			}

			switch flag.Type.Kind() {
			case reflect.String:
				fieldValue := p.value.FieldByIndex(flag.FieldIndex)
				fieldValue.SetString(args[i+1])
				i += 2
				continue
			case reflect.Uint64:
				n, err := strconv.ParseUint(args[i+1], 10, 64)
				if err != nil {
					return fmt.Errorf("TODO: could not parse int, got: %q", args[i+1])
				}
				fieldValue := p.value.FieldByIndex(flag.FieldIndex)
				fieldValue.SetUint(n)
				i += 2
				continue
			default:
				panic(fmt.Sprintf("type %s is not yet supported", flag.Type.Kind()))
			}
		} else {
			if positionalArgIndex >= len(p.command.positionals) {
				return fmt.Errorf("TODO: to manny positional arguments, arg=%q", currentArg)
			}

			positionalArg := p.command.positionals[positionalArgIndex]

			switch positionalArg.Type.Kind() {
			case reflect.Float32:
				f64, err := strconv.ParseFloat(currentArg, 32)
				if err != nil {
					return fmt.Errorf("TODO: could not parse f32: got %q", currentArg)
				}
				p.value.FieldByIndex(positionalArg.FieldIndex).SetFloat(f64)
			case reflect.Uint8:
				if len(currentArg) == 1 && (currentArg[0] < '0' || currentArg[0] > '9') {
					p.value.FieldByIndex(positionalArg.FieldIndex).SetUint(uint64(currentArg[0]))
				} else {
					n, err := strconv.ParseUint(currentArg, 10, 8)
					if err != nil {
						return fmt.Errorf("TODO: could not parse int, got: %q", currentArg)
					}
					p.value.FieldByIndex(positionalArg.FieldIndex).SetUint(n)
				}
			case reflect.Uint64:
				n, err := strconv.ParseUint(currentArg, 10, 64)
				if err != nil {
					return fmt.Errorf("TODO: could not parse int, got: %q", currentArg)
				}
				p.value.FieldByIndex(positionalArg.FieldIndex).SetUint(n)
			default:
				panic(fmt.Sprintf("type %s is not yet supported", positionalArg.Type.Kind()))
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

func hasNext(args []string, i int) bool {
	return i+1 < len(args)
}

func (p *Parser) PrintHelp(w io.Writer) {
	var sb strings.Builder

	fmt.Fprintln(&sb, "Usage:")
	fmt.Fprintf(&sb, "  %s", p.Name)
	for _, arg := range p.command.positionals {
		fmt.Fprintf(&sb, " <%s>", arg.Name)
	}

	if len(p.command.flags) > 0 {
		fmt.Fprint(&sb, " [flags]\n")
	}
	fmt.Fprintln(&sb)

	spaces := func(i int) string {
		if i < 0 {
			return ""
		}
		return strings.Repeat(" ", i)
	}

	if len(p.command.positionals) > 0 {
		fmt.Fprintln(&sb, "Positionals:")
		var maxPosLen int
		for _, arg := range p.command.positionals {
			maxPosLen = max(maxPosLen, utf8.RuneCountInString(arg.Name))
		}

		for _, arg := range p.command.positionals {
			fmt.Fprintf(&sb, "  - %s%s %s\n", arg.Name, spaces(maxPosLen-utf8.RuneCountInString(arg.Name)), arg.Help)
		}
		fmt.Fprintln(&sb)
	}

	if len(p.command.flags) > 0 {
		var (
			maxShortLen = 0
			maxLongLen  = 0
		)
		for _, flag := range p.command.flags {
			var (
				shortLen int
				longLen  int
			)
			if flag.Short != "" {
				shortLen = utf8.RuneCountInString(flag.Short) + 1 // + "-" before short flag
			}
			if flag.Long != "" {
				longLen = utf8.RuneCountInString(flag.Long) + 2 // + "--" before long flag
			}
			maxShortLen = max(maxShortLen, shortLen)
			maxLongLen = max(maxLongLen, longLen)
		}

		fmt.Fprintln(&sb, "Flags:")
		for _, flag := range p.command.flags {
			fmt.Fprintf(&sb, "  ")
			if flag.Short != "" && flag.Long != "" {
				fmt.Fprintf(&sb, "-%s,%s --%s%s",
					flag.Short,
					spaces(maxShortLen-utf8.RuneCountInString(flag.Short)-2), // "-" + ","
					flag.Long,
					spaces(maxLongLen-utf8.RuneCountInString(flag.Long)-2))
			} else if flag.Short != "" {
				fmt.Fprintf(&sb, "-%s %s",
					flag.Short, spaces(maxShortLen-utf8.RuneCountInString(flag.Short)-1+maxLongLen))
			} else {
				fmt.Fprintf(&sb, "%s--%s%s", spaces(maxShortLen), flag.Long, spaces(maxLongLen-utf8.RuneCountInString(flag.Long)-2-1))
			}
			fmt.Fprintf(&sb, "%s\n", flag.Help)
		}
	}

	fmt.Fprint(w, sb.String())
}

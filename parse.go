package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var HelpRequested = errors.New("help requested")

func Parse(opts any, args ...string) error {
	t := reflect.TypeOf(opts)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		panic("give type must be a pointer to a struct")
	}

	command := ConstructCommand(t.Elem())
	for _, arg := range args {
		if arg == "--help" || arg == "-help" || arg == "-h" {
			command.PrintHelp()
			return HelpRequested
		}
	}

	return parseArgs(reflect.ValueOf(opts).Elem(), command, args)
}

func parseArgs(v reflect.Value, command *CommandSpec, args []string) error {
	positionalArgIndex := 0

	i := 0
	for i < len(args) {
		currentArg := args[i]
		flag, ok := command.matchesFlag(currentArg)
		if ok {
			if flag.Type.Kind() == reflect.Bool {
				v.FieldByIndex(flag.FieldIndex).SetBool(true)
				i++
				continue
			}

			if !hasNext(args, i) {
				return fmt.Errorf("TODO: value for flag %q not provided", flag.Name())
			}

			switch flag.Type.Kind() {
			case reflect.String:
				fieldValue := v.FieldByIndex(flag.FieldIndex)
				fieldValue.SetString(args[i+1])
				i += 2
				continue
			case reflect.Uint64:
				n, err := strconv.ParseUint(args[i+1], 10, 64)
				if err != nil {
					return fmt.Errorf("TODO: could not parse int, got: %q", args[i+1])
				}
				fieldValue := v.FieldByIndex(flag.FieldIndex)
				fieldValue.SetUint(n)
				i += 2
				continue
			default:
				panic(fmt.Sprintf("type %s is not yet supported", flag.Type.Kind()))
			}
		} else {
			if positionalArgIndex >= len(command.positionals) {
				return fmt.Errorf("TODO: to manny positional arguments, arg=%q", currentArg)
			}

			positionalArg := command.positionals[positionalArgIndex]

			switch positionalArg.Type.Kind() {
			case reflect.Float32:
				f64, err := strconv.ParseFloat(currentArg, 32)
				if err != nil {
					return fmt.Errorf("TODO: could not parse f32: got %q", currentArg)
				}
				v.FieldByIndex(positionalArg.FieldIndex).SetFloat(f64)
			case reflect.Uint8:
				if len(currentArg) == 1 && (currentArg[0] < '0' || currentArg[0] > '9') {
					v.FieldByIndex(positionalArg.FieldIndex).SetUint(uint64(currentArg[0]))
				} else {
					n, err := strconv.ParseUint(currentArg, 10, 8)
					if err != nil {
						return fmt.Errorf("TODO: could not parse int, got: %q", currentArg)
					}
					v.FieldByIndex(positionalArg.FieldIndex).SetUint(n)
				}
			case reflect.Uint64:
				n, err := strconv.ParseUint(currentArg, 10, 64)
				if err != nil {
					return fmt.Errorf("TODO: could not parse int, got: %q", currentArg)
				}
				v.FieldByIndex(positionalArg.FieldIndex).SetUint(n)
			default:
				panic(fmt.Sprintf("type %s is not yet supported", positionalArg.Type.Kind()))
			}

			i += 1
			positionalArgIndex += 1
		}
	}

	if positionalArgIndex < len(command.positionals) {
		var sb strings.Builder
		sb.WriteString("missing values for the following positionals arugments:\n")
		for ; positionalArgIndex < len(command.positionals); positionalArgIndex++ {
			fmt.Fprintf(&sb, "  - %q\n", command.positionals[positionalArgIndex].Name)
		}
		return fmt.Errorf(sb.String())
	}

	return nil
}

func hasNext(args []string, i int) bool {
	return i+1 < len(args)
}

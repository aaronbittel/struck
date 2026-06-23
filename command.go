package main

import (
	"fmt"
	"reflect"
)

// here are some docs
type CommandSpec struct {
	name string

	flags          []*Flag
	positionalArgs []*Positional
}

func ConstructCommand(t reflect.Type) *CommandSpec {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("ConstructCommand must be called with a \"struct\" got %q", t.Kind()))
	}

	command := new(CommandSpec)

	for field := range t.Fields() {
		if !field.IsExported() {
			continue
		}

		valid, isFlag := check(field)
		if !valid {
			panic(`it is not allowed to mark a field as flag ("long" and/or "short") AND as positional ("arg")`)
		}

		if isFlag {
			command.flags = append(command.flags, FlagFromField(field))
		} else {
			command.positionalArgs = append(command.positionalArgs, NewPositionalFromField(field))
		}
	}

	return command
}

func (cmd *CommandSpec) matchesFlag(arg string) (*Flag, bool) {
	for _, flag := range cmd.flags {
		if flag.Long == arg || flag.Short == arg {
			return flag, true
		}
	}
	return nil, false
}

func check(field reflect.StructField) (isValid, isFlag bool) {
	_, hasLong := field.Tag.Lookup("long")
	_, hasShort := field.Tag.Lookup("short")
	_, isPositional := field.Tag.Lookup("arg")

	isFlag = hasLong || hasShort

	if isFlag && isPositional {
		return false, false
	}

	if isFlag {
		return true, true
	}

	return true, false
}

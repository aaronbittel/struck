package struck

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	TagLong          string = "long"
	TagShort         string = "short"
	TagPositionalArg string = "arg"
	TagHelp          string = "help"
)

type Command struct {
	name string

	flags       []*Flag
	positionals []*Positional
}

func ConstructCommand(t reflect.Type) *Command {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("ConstructCommand must be called with a \"struct\" got %q", t.Kind()))
	}

	command := new(Command)

	for field := range t.Fields() {
		if !field.IsExported() {
			continue
		}

		if hasConflictingFlagAndPositionalTags(field.Tag) {
			panic(fmt.Sprintf(
				"it is not allowed to mark a field as flag (%q and/or %q) AND as positional (%q)",
				TagLong, TagShort, TagPositionalArg))
		}

		if isFlag(field.Tag) {
			command.flags = append(command.flags, FlagFromField(field))
		} else {
			command.positionals = append(command.positionals, NewPositionalFromField(field))
		}
	}

	return command
}

func (cmd *Command) matchesFlag(arg string) (*Flag, bool) {
	for _, flag := range cmd.flags {
		if strings.HasPrefix(arg, "--") && arg[2:] == flag.Long {
			return flag, true
		}
		if strings.HasPrefix(arg, "-") && arg[1:] == flag.Short {
			return flag, true
		}
	}
	return nil, false
}

func hasConflictingFlagAndPositionalTags(tag reflect.StructTag) bool {
	_, hasLong := tag.Lookup(TagLong)
	_, hasShort := tag.Lookup(TagShort)
	_, isPositional := tag.Lookup(TagPositionalArg)
	return (hasLong || hasShort) && isPositional
}

func isFlag(tag reflect.StructTag) bool {
	_, hasLong := tag.Lookup(TagLong)
	_, hasShort := tag.Lookup(TagShort)
	return hasLong || hasShort
}

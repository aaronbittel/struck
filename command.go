package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

var (
	TagLong          string = "long"
	TagShort         string = "short"
	TagPositionalArg string = "arg"
	TagHelp          string = "help"
)

type CommandSpec struct {
	name string

	flags       []*Flag
	positionals []*Positional
}

func ConstructCommand(t reflect.Type) *CommandSpec {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("ConstructCommand must be called with a \"struct\" got %q", t.Kind()))
	}

	command := new(CommandSpec)
	command.name = "struck"

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

func (cmd *CommandSpec) PrintHelp() {
	var sb strings.Builder

	fmt.Fprintln(&sb, "Usage:")
	fmt.Fprintf(&sb, "  %s", cmd.name)
	for _, arg := range cmd.positionals {
		fmt.Fprintf(&sb, " <%s>", arg.Name)
	}

	if len(cmd.flags) > 0 {
		fmt.Fprint(&sb, " [flags]\n")
	}
	fmt.Fprintln(&sb)

	spaces := func(i int) string {
		return strings.Repeat(" ", i)
	}

	if len(cmd.positionals) > 0 {
		fmt.Fprintln(&sb, "Positionals:")
		var maxPosLen int
		for _, arg := range cmd.positionals {
			maxPosLen = max(maxPosLen, utf8.RuneCountInString(arg.Name))
		}

		for _, arg := range cmd.positionals {
			fmt.Fprintf(&sb, "  - %s%s %s\n", arg.Name, spaces(maxPosLen-utf8.RuneCountInString(arg.Name)), arg.Help)
		}
		fmt.Fprintln(&sb)
	}

	if len(cmd.flags) > 0 {
		var (
			maxShortLen = 0
			maxLongLen  = 0
		)
		for _, flag := range cmd.flags {
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
		for _, flag := range cmd.flags {
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
				fmt.Fprintf(&sb, "%s --%s%s ", spaces(maxShortLen), flag.Long, spaces(maxLongLen-utf8.RuneCountInString(flag.Long)-2-1))
			}
			fmt.Fprintf(&sb, " %s\n", flag.Help)
		}
	}

	fmt.Println(sb.String())
}

func (cmd *CommandSpec) matchesFlag(arg string) (*Flag, bool) {
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

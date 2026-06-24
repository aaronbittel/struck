package struck

import (
	"fmt"
	"io"
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

type Command struct {
	Name        string
	Schema      reflect.Value
	Flags       []*Flag
	Positionals []*Positional
}

func NewCommand(name string, schema reflect.Value) *Command {
	if schema.Kind() != reflect.Struct {
		panic(fmt.Sprintf("must be called with a \"struct\" got %q", schema.Kind()))
	}

	command := new(Command)
	command.Name = name
	command.Schema = schema

	for field := range schema.Fields() {
		if !field.IsExported() {
			continue
		}

		if hasConflictingFlagAndPositionalTags(field.Tag) {
			panic(fmt.Sprintf(
				"it is not allowed to mark a field as flag (%q and/or %q) AND as positional (%q)",
				TagLong, TagShort, TagPositionalArg))
		}

		if isFlag(field.Tag) {
			command.Flags = append(command.Flags, FlagFromField(field))
		} else {
			command.Positionals = append(command.Positionals, NewPositionalFromField(field))
		}
	}

	return command
}

func (cmd *Command) matchesFlag(arg string) (*Flag, bool) {
	for _, flag := range cmd.Flags {
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

func (cmd *Command) PrintHelp(w io.Writer) {
	var sb strings.Builder

	fmt.Fprintln(&sb, "Usage:")
	fmt.Fprintf(&sb, "  %s", cmd.Name)
	for _, arg := range cmd.Positionals {
		fmt.Fprintf(&sb, " <%s>", arg.Name)
	}

	if len(cmd.Flags) > 0 {
		fmt.Fprint(&sb, " [flags]")
	}
	fmt.Fprintf(&sb, "\n\n")

	spaces := func(i int) string {
		return strings.Repeat(" ", i)
	}

	if len(cmd.Positionals) > 0 {
		fmt.Fprintln(&sb, "Positionals:")
		var maxPosLen int
		for _, arg := range cmd.Positionals {
			maxPosLen = max(maxPosLen, utf8.RuneCountInString(arg.Name))
		}

		for _, arg := range cmd.Positionals {
			fmt.Fprintf(&sb, "  - %s%s %s\n", arg.Name, spaces(maxPosLen-utf8.RuneCountInString(arg.Name)), arg.Help)
		}
		fmt.Fprintln(&sb)
	}

	shortFlags := false
	for _, flag := range cmd.Flags {
		if flag.Short != "" {
			shortFlags = true
			break
		}
	}
	maxShortLen, maxLongLen := cmd.maxShortAndLongLenths()

	if len(cmd.Flags) > 0 {
		fmt.Fprintln(&sb, "Flags:")
		for _, flag := range cmd.Flags {
			fmt.Fprintf(&sb, "  ")
			switch {
			case flag.OnlyLong():
				shortSpaces := maxShortLen
				if shortFlags {
					shortSpaces += 1
				}
				fmt.Fprintf(&sb, "%s--%s%s",
					spaces(shortSpaces),
					flag.Long,
					spaces(maxLongLen-flag.requiredLongHelpLen()))
			case flag.OnlyShort():
				fmt.Fprintf(&sb, "-%s %s",
					flag.Short, spaces(maxShortLen-flag.requiredShortHelpLen()+maxLongLen))
			case flag.LongAndShort():
				fmt.Fprintf(&sb, "-%s,%s --%s%s",
					flag.Short,
					spaces(maxShortLen-flag.requiredShortHelpLen()),
					flag.Long,
					spaces(maxLongLen-flag.requiredLongHelpLen()))
			}
			// Without any help tags, the output will contain a trailing <space>.
			fmt.Fprintf(&sb, " %s\n", flag.Help)
		}
	}

	fmt.Fprint(w, sb.String())
}

func (cmd *Command) maxShortAndLongLenths() (maxShortLen int, maxLongLen int) {
	for _, flag := range cmd.Flags {
		maxShortLen = max(maxShortLen, flag.requiredShortHelpLen())
		maxLongLen = max(maxLongLen, flag.requiredLongHelpLen())
	}
	return maxShortLen, maxLongLen
}

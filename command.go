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
	name        string
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

func (cmd *Command) PrintHelp(w io.Writer) {
	var sb strings.Builder

	fmt.Fprintln(&sb, "Usage:")
	fmt.Fprintf(&sb, "  %s", cmd.name)
	for _, arg := range cmd.positionals {
		fmt.Fprintf(&sb, " <%s>", arg.Name)
	}

	if len(cmd.flags) > 0 {
		fmt.Fprint(&sb, " [flags]")
	}
	fmt.Fprintf(&sb, "\n\n")

	spaces := func(i int) string {
		if i < 0 {
			return ""
		}
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

	shortFlags := false
	for _, flag := range cmd.flags {
		if flag.Short != "" {
			shortFlags = true
			break
		}
	}
	maxShortLen, maxLongLen := cmd.maxShortAndLongLenths()

	if len(cmd.flags) > 0 {
		fmt.Fprintln(&sb, "Flags:")
		for _, flag := range cmd.flags {
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
	for _, flag := range cmd.flags {
		maxShortLen = max(maxShortLen, flag.requiredShortHelpLen())
		maxLongLen = max(maxLongLen, flag.requiredLongHelpLen())
	}
	return maxShortLen, maxLongLen
}

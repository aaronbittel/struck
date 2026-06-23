package struck

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Flag struct {
	Long       string
	Short      string
	Type       reflect.Type
	FieldIndex []int
	Help       string
}

func FlagFromField(field reflect.StructField) *Flag {
	return &Flag{
		Long:       field.Tag.Get(TagLong),
		Short:      field.Tag.Get(TagShort),
		Type:       field.Type,
		FieldIndex: field.Index,
		Help:       field.Tag.Get(TagHelp),
	}
}

func (f Flag) Name() string {
	if f.Long != "" {
		return f.Long
	}
	return f.Short
}

func (f Flag) String() string {
	var sb strings.Builder

	if f.Long != "" {
		fmt.Fprintf(&sb, "long=%s ", f.Long)
	}
	if f.Short != "" {
		fmt.Fprintf(&sb, "short=%s ", f.Short)
	}
	if f.Type.Name() != "" {
		fmt.Fprintf(&sb, "type=%s ", f.Type.Name())
	} else {
		fmt.Fprintf(&sb, "type=%s ", f.Type.Kind())
	}
	fmt.Fprintf(&sb, "index=%v ", f.FieldIndex)

	return sb.String()
}

func (f Flag) requiredShortHelpLen() int {
	if f.Short == "" {
		return 0
	}

	if f.Long == "" {
		return utf8.RuneCountInString(f.Short) + 1 // for "-"
	}

	return utf8.RuneCountInString(f.Short) + 1 + 1 // for "-", "," (short and long form are comma seperated)
}

func (f Flag) requiredLongHelpLen() int {
	if f.Long == "" {
		return 0
	}

	return utf8.RuneCountInString(f.Long) + 2 // for "--"
}

func (f Flag) OnlyShort() bool {
	return f.Short != "" && f.Long == ""
}

func (f Flag) OnlyLong() bool {
	return f.Short == "" && f.Long != ""
}

func (f Flag) LongAndShort() bool {
	return f.Short != "" && f.Long != ""
}

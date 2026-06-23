package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Flag struct {
	Long       string
	Short      string
	Type       reflect.Type
	FieldIndex []int
}

type emptyStructTag struct {
	name string
}

func (e emptyStructTag) Error() string {
	return fmt.Sprintf("struct tag %q must not be empty", e.name)
}

func FlagFromField(field reflect.StructField) *Flag {
	var long string
	longTag, ok := field.Tag.Lookup("long")
	if ok && longTag != "" {
		long = longTag
	}

	var short string
	shortTag, ok := field.Tag.Lookup("short")
	if ok && shortTag != "" {
		short = shortTag
	}

	return &Flag{
		Long:       long,
		Short:      short,
		Type:       field.Type,
		FieldIndex: field.Index,
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

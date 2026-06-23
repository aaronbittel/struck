package main

import "reflect"

type Positional struct {
	Name       string
	Type       reflect.Type
	FieldIndex []int
	Help       string
}

func NewPositionalFromField(field reflect.StructField) *Positional {
	name := field.Name
	argName, ok := field.Tag.Lookup("arg")
	if ok && argName != "" {
		name = argName
	}
	return &Positional{
		Name:       name,
		Type:       field.Type,
		FieldIndex: field.Index,
		Help:       field.Tag.Get("help"),
	}
}

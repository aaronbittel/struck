package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Parse(opts any, args ...string) error {
	v := reflect.ValueOf(opts).Elem()

	var (
		argIndex int
		posIndex int
		posArgs  []arg
	)

	for field, value := range v.Fields() {
		if isFlag(field) {
			continue
		}
		argTag, ok := field.Tag.Lookup("arg")
		name := field.Name
		if ok && argTag != "" {
			name = argTag
		}
		posArgs = append(posArgs, arg{name: name, field: field, value: value})
	}

outer:
	for argIndex < len(args) {
		for field, value := range v.Fields() {
			_, ok := matches(args[argIndex], field)
			if !ok {
				continue
			}

			if requiresValue(value.Kind()) && argIndex+1 >= len(args) {
				return ParseError{kind: ErrMissingValue, paramKind: KindFlag, name: args[argIndex]}
			}

			switch field.Type.Kind() {
			case reflect.String:
				value.SetString(args[argIndex+1])
				argIndex += 2
				continue outer
			case reflect.Uint64:
				age, err := strconv.ParseUint(args[argIndex+1], 10, 64)
				if err != nil {
					return ParseError{
						kind: ErrInvalidValue, argStr: args[argIndex+1], name: args[argIndex], expectedType: field.Type.Kind(),
					}
				}
				value.SetUint(age)
				argIndex += 2
				continue outer
			}

		}

		if posIndex < len(posArgs) {
			posArg := posArgs[posIndex]
			switch posArg.field.Type.Kind() {
			case reflect.String:
				posArg.value.SetString(args[posIndex])
			case reflect.Float32:
				v, err := strconv.ParseFloat(args[argIndex], 32)
				if err != nil {
					return ParseError{
						kind:         ErrInvalidValue,
						argStr:       args[argIndex],
						expectedType: posArg.field.Type.Kind(),
						paramKind:    KindPositional,
						name:         posArg.name,
						err:          err,
					}
				}
				posArg.value.SetFloat(v)
			default:
				panic(fmt.Sprintf("unsupported type %s for positional argument", posArg.field.Type.Kind()))
			}

			argIndex++
			continue outer
		}

		fmt.Printf("skipping arg[%d] = %s\n", argIndex, args[argIndex])
		argIndex++
	}

	return nil
}

type Arg struct {
	field    reflect.StructField
	value    reflect.Value
	required bool
}

func (p Arg) name() string {
	argTag := p.field.Tag.Get("arg")
	if argTag == "" {
		return p.field.Name
	}
	return argTag
}

type Flag struct {
	long     string
	short    string
	required bool

	value reflect.Value
}

func (f Flag) name() string {
	if f.long == "" {
		return f.short
	}
	return f.long
}

type arg struct {
	name  string
	field reflect.StructField
	value reflect.Value
}

func matches(arg string, field reflect.StructField) (flag string, ok bool) {
	long, ok := field.Tag.Lookup("long")
	if ok && strings.HasPrefix(arg, "--") && arg[2:] == long {
		return long, true
	}

	short, ok := field.Tag.Lookup("short")
	if ok && strings.HasPrefix(arg, "-") && arg[1:] == short {
		return short, true
	}

	return "", false
}

func requiresValue(flagType reflect.Kind) bool {
	return flagType != reflect.Bool
}

func isFlag(field reflect.StructField) bool {
	if _, ok := field.Tag.Lookup("long"); ok {
		return true
	}
	if _, ok := field.Tag.Lookup("short"); ok {
		return true
	}
	return false
}

func isRequired(field reflect.StructField) bool {
	v, ok := field.Tag.Lookup("required")
	if ok && v == "true" {
		return true
	}
	for tag := range strings.FieldsSeq(string(field.Tag)) {
		if tag == "required" {
			return true
		}
	}
	return false
}

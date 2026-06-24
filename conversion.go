package struck

import (
	"fmt"
	"reflect"
	"strconv"
)

// TODO: support hex, octal, binary for numbers

var bitSizes = map[reflect.Kind]int{
	reflect.Int:     strconv.IntSize,
	reflect.Int8:    8,
	reflect.Int16:   16,
	reflect.Int32:   32,
	reflect.Int64:   64,
	reflect.Uint:    strconv.IntSize,
	reflect.Uint8:   8,
	reflect.Uint16:  16,
	reflect.Uint32:  32,
	reflect.Uint64:  64,
	reflect.Float32: 32,
	reflect.Float64: 64,
}

type parseError struct {
	target reflect.Kind
	arg    string
	err    error
}

func (p parseError) Error() string {
	return fmt.Sprintf("could not parse %q into %s", p.arg, p.target)
}

func (p parseError) Unwrap() error {
	return p.err
}

var unsupportedKind = func(kind reflect.Kind) error {
	return fmt.Errorf("unsupported kind: %s", kind)
}

func SetValue(v reflect.Value, arg string) error {
	switch v.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(arg)
		if err != nil {
			return parseError{target: v.Kind(), arg: arg, err: err}
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n64, err := parseInt(v.Kind(), arg)
		if err != nil {
			return err
		}
		v.SetInt(n64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		nU64, err := parseUint(v.Kind(), arg)
		if err != nil {
			return err
		}
		v.SetUint(nU64)
	case reflect.Float32, reflect.Float64:
		f64, err := parseFloat(v.Kind(), arg)
		if err != nil {
			return err
		}
		v.SetFloat(f64)
	case reflect.Array:
		return fmt.Errorf("arrays are not supported, please use slices")
	case reflect.Slice:
		elemType := v.Type().Elem()
		elem := reflect.New(elemType).Elem()
		if err := SetValue(elem, arg); err != nil {
			return fmt.Errorf("could not append to slice: %s", err)
		}
		v.Set(reflect.Append(v, elem))
	case reflect.String:
		v.SetString(arg)
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func,
		reflect.Interface, reflect.UnsafePointer, reflect.Map, reflect.Pointer,
		reflect.Struct, reflect.Uintptr:
		return unsupportedKind(v.Kind())
	default:
		panic(fmt.Sprintf("unknown kind %q", v.Kind()))
	}

	return nil
}

func parseInt(kind reflect.Kind, s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, bitSizes[kind])
	if err != nil {
		return 0, parseError{target: kind, arg: s, err: err}
	}
	return i64, nil
}

func parseUint(kind reflect.Kind, s string) (uint64, error) {
	u64, err := strconv.ParseUint(s, 10, bitSizes[kind])
	if err != nil {
		return 0, parseError{target: kind, arg: s, err: err}
	}
	return u64, nil
}

func parseFloat(kind reflect.Kind, s string) (float64, error) {
	f64, err := strconv.ParseFloat(s, bitSizes[kind])
	if err != nil {
		return 0.0, parseError{target: kind, arg: s, err: err}
	}
	return f64, nil
}

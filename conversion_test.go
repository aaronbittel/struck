package struck

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversionSlice(t *testing.T) {
	validTests := []struct {
		name  string
		input []string
		want  []int
	}{
		{
			name:  "insert one",
			input: []string{"5"},
			want:  []int{5},
		},
		{
			name:  "insert two",
			input: []string{"5", "6"},
			want:  []int{5, 6},
		},
		{
			name:  "insert three",
			input: []string{"5", "6", "7"},
			want:  []int{5, 6, 7},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			var s []int
			v := reflect.ValueOf(&s).Elem()

			for _, arg := range tt.input {
				require.NoError(t, SetValue(v, arg))
			}

			assert.Equal(t, tt.want, v.Interface())
		})
	}

	invalidTests := []struct {
		name    string
		input   []string
		wantErr []bool
	}{
		{
			name:    "wrong type",
			input:   []string{"5", "abc"},
			wantErr: []bool{false, true},
		},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			var s []int
			v := reflect.ValueOf(&s).Elem()

			for i, arg := range tt.input {
				err := SetValue(v, arg)
				if tt.wantErr[i] {
					assert.ErrorContains(t, SetValue(v, arg), "could not append to slice")
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestConversionString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "non-empty", input: "hello", want: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s string
			v := reflect.ValueOf(&s).Elem()
			require.NoError(t, SetValue(v, tt.input))
			assert.Equal(t, tt.want, v.String())
		})
	}
}

func TestConversionBool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{name: "1", input: "1", want: true},
		{name: "t", input: "t", want: true},
		{name: "T", input: "T", want: true},
		{name: "TRUE", input: "TRUE", want: true},
		{name: "true", input: "true", want: true},
		{name: "True", input: "True", want: true},
		{name: "0", input: "0", want: false},
		{name: "f", input: "f", want: false},
		{name: "F", input: "F", want: false},
		{name: "FALSE", input: "FALSE", want: false},
		{name: "false", input: "false", want: false},
		{name: "False", input: "False", want: false},
		{name: "unknown", input: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bool
			v := reflect.ValueOf(&b).Elem()
			err := SetValue(v, tt.input)
			if tt.wantErr {
				assert.ErrorContains(t, err, "could not parse")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, b)
			}
		})
	}
}

var (
	TypeInt   = reflect.TypeFor[int]()
	TypeInt8  = reflect.TypeFor[int8]()
	TypeInt16 = reflect.TypeFor[int16]()
	TypeInt32 = reflect.TypeFor[int32]()
	TypeInt64 = reflect.TypeFor[int64]()
)

func TestConversionInt(t *testing.T) {
	var intTypes = []reflect.Type{TypeInt8, TypeInt16, TypeInt32, TypeInt64}

	tests := []struct {
		name        string
		input       string
		want        int64
		expectError map[reflect.Type]bool
	}{
		{name: "zero", input: "0", want: 0},
		{name: "small positive", input: "123", want: 123},
		{name: "small negative", input: "-123", want: -123},
		{name: "greater maxInt8", input: "128", want: 128, expectError: map[reflect.Type]bool{
			TypeInt8: true,
		}},
		{name: "less minInt8", input: "-129", want: -129, expectError: map[reflect.Type]bool{
			TypeInt8: true,
		}},
		{name: "greater maxInt16", input: "32768", want: 32768, expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
		}},
		{name: "less minInt16", input: "-32769", want: -32769, expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
		}},
		{name: "greater maxInt32", input: "2147483648", want: 2147483648, expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
			TypeInt32: true,
		}},
		{name: "less minInt32", input: "-2147483649", want: -2147483649, expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
			TypeInt32: true,
		}},
		{name: "greater maxInt64", input: "9223372036854775808", expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
			TypeInt32: true,
			TypeInt64: true,
		}},
		{name: "less minInt64", input: "-9223372036854775809", expectError: map[reflect.Type]bool{
			TypeInt8:  true,
			TypeInt16: true,
			TypeInt32: true,
			TypeInt64: true,
		}},
		{name: "not a number", input: "asdf", expectError: map[reflect.Type]bool{
			TypeInt:   true,
			TypeInt8:  true,
			TypeInt16: true,
			TypeInt32: true,
			TypeInt64: true,
		}},
	}

	for _, tt := range tests {
		for _, intType := range intTypes {
			t.Run(fmt.Sprintf("%s/%s", tt.name, intType), func(t *testing.T) {
				v := reflect.New(intType).Elem()
				err := SetValue(v, tt.input)
				if expectErr, ok := tt.expectError[intType]; ok && expectErr {
					assert.ErrorContains(t, err, "could not parse")
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.want, v.Int())
				}
			})
		}

	}

	intTests := []struct {
		name      string
		input     string
		want      int64
		expectErr bool
	}{
		{name: "zero", input: "0", want: 0},
		{name: "small positive", input: "123", want: 123},
		{name: "small negative", input: "-123", want: -123},
		{name: "greater maxInt8", input: "128", want: 128},
		{name: "less minInt8", input: "-129", want: -129},
		{name: "greater maxInt16", input: "32768", want: 32768},
		{name: "less minInt16", input: "-32769", want: -32769},
		{name: "greater maxInt32", input: "2147483648", want: 2147483648, expectErr: strconv.IntSize == 32},
		{name: "less minInt32", input: "-2147483649", want: -2147483649, expectErr: strconv.IntSize == 32},
		{name: "greater maxInt64", input: "9223372036854775808", expectErr: true},
		{name: "less minInt64", input: "-9223372036854775809", expectErr: true},
		{name: "not a number", input: "asdf", expectErr: true},
	}

	for _, tt := range intTests {
		t.Run(fmt.Sprintf("%s/%s", tt.name, TypeInt), func(t *testing.T) {
			v := reflect.New(TypeInt).Elem()
			err := SetValue(v, tt.input)
			if tt.expectErr {
				assert.ErrorContains(t, err, "could not parse")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, v.Int())
			}
		})
	}
}

var (
	TypeUint   = reflect.TypeFor[uint]()
	TypeUint8  = reflect.TypeFor[uint8]()
	TypeUint16 = reflect.TypeFor[uint16]()
	TypeUint32 = reflect.TypeFor[uint32]()
	TypeUint64 = reflect.TypeFor[uint64]()
)

func TestConversionUint(t *testing.T) {
	var uintTypes = []reflect.Type{TypeUint8, TypeUint16, TypeUint32, TypeUint64}

	allUints := map[reflect.Type]bool{
		TypeUint8:  true,
		TypeUint16: true,
		TypeUint32: true,
		TypeUint64: true,
	}

	tests := []struct {
		name        string
		input       string
		want        uint64
		expectError map[reflect.Type]bool
	}{
		{name: "zero", input: "0", want: 0},
		{name: "negative", input: "-123", expectError: allUints},
		{name: "less maxUint8", input: "123", want: 123},
		{name: "greater maxUint8", input: "256", want: 256, expectError: map[reflect.Type]bool{
			TypeUint8: true,
		}},
		{name: "greater maxUint16", input: "65536", want: 65536, expectError: map[reflect.Type]bool{
			TypeUint8:  true,
			TypeUint16: true,
		}},
		{name: "greater maxUint32", input: "4294967296", want: 4294967296, expectError: map[reflect.Type]bool{
			TypeUint8:  true,
			TypeUint16: true,
			TypeUint32: true,
		}},
		{name: "greater maxUint64", input: "18446744073709551616", expectError: allUints},
		{name: "not a number", input: "asdf", expectError: allUints},
	}

	for _, tt := range tests {
		for _, intType := range uintTypes {
			t.Run(fmt.Sprintf("%s/%s", tt.name, intType), func(t *testing.T) {
				v := reflect.New(intType).Elem()
				err := SetValue(v, tt.input)
				if expectErr, ok := tt.expectError[intType]; ok && expectErr {
					assert.ErrorContains(t, err, "could not parse")
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.want, v.Uint())
				}
			})
		}
	}

	uintTests := []struct {
		name      string
		input     string
		want      uint64
		expectErr bool
	}{
		{name: "zero", input: "0", want: 0},
		{name: "negative", input: "-123", expectErr: true},
		{name: "small positive", input: "123", want: 123},
		{name: "greater maxInt8", input: "256", want: 256},
		{name: "greater maxInt16", input: "65536", want: 65536},
		{name: "greater maxInt32", input: "4294967296", want: 4294967296, expectErr: strconv.IntSize == 32},
		{name: "greater maxInt64", input: "18446744073709551616", expectErr: true},
		{name: "not a number", input: "asdf", expectErr: true},
	}

	for _, tt := range uintTests {
		t.Run(fmt.Sprintf("%s/%s", tt.name, TypeUint), func(t *testing.T) {
			v := reflect.New(TypeUint).Elem()
			err := SetValue(v, tt.input)
			if tt.expectErr {
				assert.ErrorContains(t, err, "could not parse")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, v.Uint())
			}
		})
	}
}

var (
	TypeFloat32 = reflect.TypeFor[float32]()
	TypeFloat64 = reflect.TypeFor[float64]()
)

func TestConversionFloat(t *testing.T) {
	var floatTypes = []reflect.Type{TypeFloat32, TypeFloat64}

	tests := []struct {
		name        string
		input       string
		want        float64
		expectError map[reflect.Type]bool
	}{
		{name: "zero", input: "0", want: 0},
		{name: "small positive", input: "123.5", want: 123.5},
		{name: "small negative", input: "-123.5", want: -123.5},
		{name: "integer-like", input: "42", want: 42},
		{name: "invalid", input: "asdf", expectError: map[reflect.Type]bool{
			TypeFloat32: true,
			TypeFloat64: true,
		}},
	}

	for _, tt := range tests {
		for _, floatType := range floatTypes {
			t.Run(fmt.Sprintf("%s/%s", tt.name, floatType), func(t *testing.T) {
				v := reflect.New(floatType).Elem()
				err := SetValue(v, tt.input)

				if expectErr, ok := tt.expectError[floatType]; ok && expectErr {
					assert.ErrorContains(t, err, "could not parse")
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.want, v.Float())
				}
			})
		}
	}
}

package struck

import (
	"io"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructParsing(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want *Parser
	}{
		{
			name: "single string field",
			in: &struct {
				Name string `long:"name"`
			}{},
			want: &Parser{
				Root: &Command{
					Name: "Test",
					Flags: []*Flag{
						{Long: "name", FieldIndex: []int{0}},
					},
				},
				Subcommands: []*Command{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.in, "Test")
			diff := cmp.Diff(tt.want, parser, cmpopts.IgnoreTypes(reflect.Value{}), cmpopts.IgnoreFields(Parser{}, "Output"))
			if diff != "" {
				t.Errorf("diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStructSubcommandParsing(t *testing.T) {
	type subcommand struct {
		name   string
		schema any
	}

	tests := []struct {
		name        string
		in          any
		subcommands []subcommand
		want        *Parser
	}{
		{
			name: "single subcommand",
			in:   &struct{}{},
			subcommands: []subcommand{
				{
					name: "subcommand",
					schema: &struct {
						Name string `long:"name"`
					}{},
				},
			},
			want: &Parser{
				Root: &Command{Name: "Test"},
				Subcommands: []*Command{
					{
						Name: "subcommand",
						Flags: []*Flag{
							{Long: "name", FieldIndex: []int{0}},
						},
					},
				},
			},
		},
		{
			name: "multiple subcommands",
			in: &struct {
				Root       string `long:"root" short:"r" help:"this is root"`
				Positional byte   `arg:"byte"`
			}{},
			subcommands: []subcommand{
				{
					name: "sub1",
					schema: &struct {
						Name string `long:"name" short:"n"`
						Age  int    `long:"age"`
					}{},
				},
				{
					name: "sub2",
					schema: &struct {
						Verbose []bool `long:"verbose" short:"v" help:"adjust verbosity"`
						Value   uint16 `short:"special"`
					}{},
				},
			},
			want: &Parser{
				Root: &Command{
					Name: "Test",
					Flags: []*Flag{
						&Flag{Long: "root", Short: "r", Help: "this is root", FieldIndex: []int{0}},
					},
					Positionals: []*Positional{
						&Positional{Name: "byte", FieldIndex: []int{1}},
					},
				},
				Subcommands: []*Command{
					{
						Name: "sub1",
						Flags: []*Flag{
							{Long: "name", Short: "n", FieldIndex: []int{0}},
							{Long: "age", FieldIndex: []int{1}},
						},
					},
					{
						Name: "sub2",
						Flags: []*Flag{
							{Long: "verbose", Short: "v", Help: "adjust verbosity", FieldIndex: []int{0}},
							{Short: "special", FieldIndex: []int{1}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.in, "Test")

			for _, subcmd := range tt.subcommands {
				parser.AddSubcommand(subcmd.name, subcmd.schema)
			}

			diff := cmp.Diff(tt.want, parser, cmpopts.IgnoreTypes(reflect.Value{}), cmpopts.IgnoreFields(Parser{}, "Output"))
			if diff != "" {
				t.Errorf("diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	happyTests := []struct {
		name     string
		args     []string
		wantName string
	}{
		{
			name:     "no flag",
			args:     []string{},
			wantName: "",
		},
		{
			name:     "set name",
			args:     []string{"--name", "Bob"},
			wantName: "Bob",
		},
		{
			name:     "override flag",
			args:     []string{"--name", "Bob", "-n", "Alice"},
			wantName: "Alice",
		},
	}

	for _, tt := range happyTests {
		t.Run(tt.name, func(t *testing.T) {
			var schema struct {
				Name string `long:"name" short:"n"`
			}

			parser := NewParser(&schema, "Root")
			err := parser.ParseArgs(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, schema.Name)
		})
	}

	errorTests := []struct {
		name       string
		args       []string
		wantErrMsg string
	}{
		{
			name:       "help requested",
			args:       []string{"-help"},
			wantErrMsg: "help requested",
		},
		{
			name:       "unknown flag",
			args:       []string{"--unknown"},
			wantErrMsg: "unknown flag",
		},
		{
			name:       "wrong type",
			args:       []string{"--number", "hello"},
			wantErrMsg: "could not parse",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			var schema struct {
				Number int `long:"number"`
			}

			parser := NewParser(&schema, "Root")
			parser.Output = io.Discard
			assert.ErrorContains(t, parser.parseArgs(tt.args), tt.wantErrMsg)
		})
	}

}

func TestParseArgsBool(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantSingleBool bool
		wantBoolSlice  []bool
	}{
		{
			name:           "no flag",
			args:           []string{},
			wantSingleBool: false,
			wantBoolSlice:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var schema struct {
				SingleBool bool   `long:"dodoit" short:"doit"`
				BoolSlice  []bool `long:"verbose" short:"v"`
			}

			parser := NewParser(&schema, "Root")
			err := parser.ParseArgs(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.wantSingleBool, schema.SingleBool)
			assert.Equal(t, tt.wantBoolSlice, schema.BoolSlice)
		})
	}
}

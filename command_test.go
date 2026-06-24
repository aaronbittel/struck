package struck

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseStructIntoCommand(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want *Command
	}{
		{
			name: "single string flag",
			typ: reflect.TypeFor[struct {
				Name string "long:\"name\""
			}](),
			want: &Command{
				flags: []*Flag{
					{
						Long:       "name",
						Type:       StringType,
						FieldIndex: []int{0},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstructCommand(tt.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

var StringType = reflect.TypeFor[string]()

func TestPrintHelp(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Command
	}{
		{
			name: "single long flag",
			cmd: &Command{
				name:  "test",
				flags: []*Flag{&Flag{Long: "name"}},
			},
		},
		{
			name: "two long flags",
			cmd: &Command{
				name: "test",
				flags: []*Flag{
					&Flag{Long: "name"},
					&Flag{Long: "age", Short: "a"},
				},
			},
		},
		{
			name: "two long flags with help",
			cmd: &Command{
				name: "test",
				flags: []*Flag{
					&Flag{Long: "name", Help: "here is some help message"},
					&Flag{Long: "age", Short: "a", Help: "here as well"},
				},
			},
		},
		{
			name: "long flag name",
			cmd: &Command{
				name: "test",
				flags: []*Flag{
					&Flag{Long: "name", Help: "here is some help message"},
					&Flag{Long: "age", Short: "a", Help: "here as well"},
					&Flag{Long: "this-is-quite-a-long-flag", Help: "which also has some help text"},
				},
			},
		},
		{
			name: "short flag with more than one character",
			cmd: &Command{
				name: "test",
				flags: []*Flag{
					&Flag{Long: "name", Short: "longer", Help: "here is some help message"},
					&Flag{Long: "age", Short: "a", Help: "here as well"},
				},
			},
		},
		{
			name: "different short lengths",
			cmd: &Command{
				name: "test",
				flags: []*Flag{
					&Flag{Short: "name", Help: "here is some help message"},
					&Flag{Short: "age"},
					&Flag{Short: "z", Help: "xxx"},
				},
			},
		},
		{
			name: "positionals",
			cmd: &Command{
				name: "test",
				positionals: []*Positional{
					&Positional{Name: "ArgNotSet", Help: "here was no `arg` tag set"},
					&Positional{Name: "age"},
					&Positional{Name: "justAnotherPositional"},
					&Positional{Name: "z", Help: "xxx"},
				},
			},
		},
		{
			name: "comprehensive list of arguments",
			cmd: &Command{
				name: "deploy",
				flags: []*Flag{
					{Long: "verboseFlaggg", Short: "v", Help: "Enable verbose output"},
					{Long: "config", Help: "Path to configuration file"},
					{Short: "q", Help: "Quiet mode"},
					{Long: "force", Short: "f"},
					{Long: "dry-run", Help: "dry-run"},
					{Short: "x"},
					{Long: "outputoutputoutput", Short: "o", Help: "Write generated artifacts to the specified directory"},
				},
				positionals: []*Positional{
					{Name: "source", Help: "Source directory"},
					{Name: "target"},
					{Name: "environment", Help: "Deployment environment"},
					{Name: "version", Help: ""},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.cmd.PrintHelp(&buf)
			g := goldie.New(t)
			// show spaces visibly as ·
			got := bytes.ReplaceAll(buf.Bytes(), []byte(" "), []byte("·"))
			// normalize file name
			g.Assert(t, strings.ReplaceAll(tt.name, " ", "_"), got)
		})
	}
}

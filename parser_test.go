package struck

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintHelp(t *testing.T) {
	var opt struct {
		Name string `long:"name"`
	}

	parser := NewParser(&opt, "test")

	var sb strings.Builder
	parser.PrintHelp(&sb)

	want := `Usage:
  test [flags]

Flags:
  --name
`

	assert.Equal(t, want, sb.String())
}

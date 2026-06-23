package struck

import (
	"reflect"
	"testing"

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
			typ: reflect.TypeOf(struct {
				Name string `long:"name"`
			}{}),
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

var StringType = reflect.TypeOf("")

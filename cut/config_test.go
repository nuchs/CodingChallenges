package main_test

import (
	"reflect"
	"testing"

	cut "github.com/nuchs/cc/cut"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want cut.Config
	}{
		{
			desc: "Default input from stdin",
			args: []string{},
			want: cut.Config{Field: 0, In: "-"},
		},
		{
			desc: "Input from file",
			args: []string{"main.go"},
			want: cut.Config{Field: 0, In: "main.go"},
		},
		{
			desc: "field selected",
			args: []string{"-f", "7"},
			want: cut.Config{Field: 7, In: "-"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := cut.NewConfig(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Wrong config: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

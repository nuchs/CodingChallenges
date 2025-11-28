package main_test

import (
	"errors"
	"reflect"
	"testing"

	cc "github.com/nuchs/cc/cut"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want cc.Config
	}{
		{
			desc: "Default input from stdin",
			args: []string{"-f", "1"},
			want: cc.Config{Field: 1, In: "-", Delimiter: "\t"},
		},
		{
			desc: "Input from file",
			args: []string{"-f", "1", "main.go"},
			want: cc.Config{Field: 1, In: "main.go", Delimiter: "\t"},
		},
		{
			desc: "field selected",
			args: []string{"-f", "7"},
			want: cc.Config{Field: 7, In: "-", Delimiter: "\t"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := cc.NewConfig(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Wrong config: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func TestBadConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want error
	}{
		{
			desc: "No field",
			args: []string{},
			want: cc.ErrBadSelector,
		},
		{
			desc: "invalid field",
			args: []string{"-f", "0"},
			want: cc.ErrBadSelector,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := cc.NewConfig(tC.args)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !errors.Is(err, cc.ErrBadSelector) {
				t.Fatalf("Wrong error - got %q, want %q", err, cc.ErrBadSelector)
			}
		})
	}
}

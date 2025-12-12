package main_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	be "github.com/nuchs/cc/lb/cmd/be"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want be.Config
	}{
		{
			desc: "set name",
			args: []string{"-n", "bob", "-p", "0"},
			want: be.Config{Name: "bob"},
		},
		{
			desc: "set port",
			args: []string{"-n", "b", "-p", "100"},
			want: be.Config{Name: "b", Port: 100},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := be.NewConfig(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad Config: got: %+v, want: %+v", got, tC.want)
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
			desc: "no name",
			args: []string{"-p", "1"},
			want: be.ErrNoName,
		},
		{
			desc: "no port",
			args: []string{"-n", "a"},
			want: be.ErrInvalidPort,
		},
		{
			desc: "negative port",
			args: []string{"-n", "a", "-p", "-1"},
			want: be.ErrInvalidPort,
		},
		{
			desc: "port too big",
			args: []string{"-n", "a", "-p", "65536"},
			want: be.ErrInvalidPort,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := be.NewConfig(tC.args)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if !errors.Is(err, tC.want) {
				t.Fatalf("Wrong error. got: %s, want %s", err, tC.want)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	args := []string{"-n", "a", "-p", "bacon"}
	want := "parse error"
	_, err := be.NewConfig(args)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Wrong error. got: %s, want %q", err, want)
	}
}

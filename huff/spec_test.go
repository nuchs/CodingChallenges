package main_test

import (
	"errors"
	"reflect"
	"testing"

	hf "github.com/nuchs/cc/huff"
)

func TestFlags(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want hf.Spec
	}{
		{
			desc: "Print header",
			args: []string{"-a", "p", "-i", "someFile.txt"},
			want: hf.Spec{Action: hf.Printing, In: "someFile.txt"},
		},
		{
			desc: "Encode file",
			args: []string{"-a", "e", "-i", "infile.txt", "-o", "outfile.txt"},
			want: hf.Spec{Action: hf.Encoding, In: "infile.txt", Out: "outfile.txt"},
		},
		{
			desc: "Decode file",
			args: []string{"-a", "d", "-i", "infile.txt", "-o", "outfile.txt"},
			want: hf.Spec{Action: hf.Decoding, In: "infile.txt", Out: "outfile.txt"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := hf.NewSpec(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error loading spec: %s", err)
			}
			if !reflect.DeepEqual(*got, tC.want) {
				t.Fatalf("Bad spec: got:%+v, want:%+v", *got, tC.want)
			}
		})
	}
}

func TestFlagErrors(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		err  error
	}{
		{desc: "No action", args: []string{"-i", "in"}, err: hf.ErrActionNotSet},
		{desc: "Print no src", args: []string{"-a", "p"}, err: hf.ErrNoInputFile},
		{
			desc: "Encode no src",
			args: []string{"-a", "p", "-o", "out"},
			err:  hf.ErrNoInputFile,
		},
		{
			desc: "Decode no src",
			args: []string{"-a", "p", "-o", "out"},
			err:  hf.ErrNoInputFile,
		},
		{
			desc: "Encode no sink",
			args: []string{"-a", "e", "-i", "in"},
			err:  hf.ErrNoOutputFile,
		},
		{
			desc: "Decode no sink",
			args: []string{"-a", "d", "-i", "in"},
			err:  hf.ErrNoOutputFile,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s, err := hf.NewSpec(tC.args)
			if err == nil {
				t.Fatalf("Expected error, got spec: %+v", s)
			}
			if !errors.Is(err, tC.err) {
				t.Fatalf("Wrong error, got %s, want %s", err, tC.err)
			}
		})
	}
}

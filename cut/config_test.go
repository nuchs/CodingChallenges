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
			want: cc.Config{
				In: "-", Spread: cc.Spreads{{Min: 1, Max: 1}}, Delimiter: "\t",
			},
		},
		{
			desc: "Input from file",
			args: []string{"-f", "1", "main.go"},
			want: cc.Config{
				In:        "main.go",
				Spread:    cc.Spreads{{Min: 1, Max: 1}},
				Delimiter: "\t",
			},
		},
		{
			desc: "field selection",
			args: []string{"-f", "7"},
			want: cc.Config{
				In:        "-",
				Spread:    cc.Spreads{{Min: 7, Max: 7}},
				Delimiter: "\t",
			},
		},
		{
			desc: "delimiter",
			args: []string{"-f", "1", "-d", ","},
			want: cc.Config{
				In:        "-",
				Spread:    cc.Spreads{{Min: 1, Max: 1}},
				Delimiter: ",",
			},
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
			desc: "zero isn't a valid field",
			args: []string{"-f", "0"},
			want: cc.ErrZeroSpreadBound,
		},
		{
			desc: "bad field spec",
			args: []string{"-f", "1-2-3"},
			want: cc.ErrBadlyFormattedSpread,
		},
		{
			desc: "Non numeric lower bound",
			args: []string{"-f", "a-3"},
			want: cc.ErrNonNumericSpreadBound,
		},
		{
			desc: "Non numeric upper bound",
			args: []string{"-f", "1-x"},
			want: cc.ErrNonNumericSpreadBound,
		},
		{
			desc: "upper bound greater than lower",
			args: []string{"-f", "2-1"},
			want: cc.ErrBadSpreadOrder,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := cc.NewConfig(tC.args)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !errors.Is(err, tC.want) {
				t.Fatalf("Wrong error - got %q, want %q", err, cc.ErrBadSelector)
			}
		})
	}
}

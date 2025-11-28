package main_test

import (
	"strings"
	"testing"

	cc "github.com/nuchs/cc/cut"
)

func TestCut(t *testing.T) {
	testCases := []struct {
		desc string
		cfg  cc.Config
		in   string
		want string
	}{
		{
			desc: "Print all lines with no delimiter",
			cfg:  cc.Config{Field: 1, Delimiter: "\t"},
			in:   "abc\n\tdef\nghi",
			want: "abc\n\nghi\n",
		},
		{
			desc: "Fields counted from 1",
			cfg:  cc.Config{Field: 1, Delimiter: "\t"},
			in:   "111\t222",
			want: "111\n",
		},
		{
			desc: "Cuts are done linewise",
			cfg:  cc.Config{Field: 2, Delimiter: "\t"},
			in:   "111\t222\n111\t222\n111\t",
			want: "222\n222\n\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var out strings.Builder
			in := strings.NewReader(tC.in)
			cc.Cut(in, &out, tC.cfg)
			got := out.String()
			if got != tC.want {
				t.Fatalf("Bad cut - got: %q, want: %q", got, tC.want)
			}
		})
	}
}

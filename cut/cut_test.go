package main_test

import (
	"strings"
	"testing"

	cc "github.com/nuchs/cc/cut"
)

func TestCut(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		in   string
		want string
	}{
		{
			desc: "Print all lines with no delimiter",
			args: []string{"-f", "1"},
			in:   "abc\n\tdef\nghi",
			want: "abc\n\nghi\n",
		},
		{
			desc: "Fields counted from 1",
			args: []string{"-f", "1"},
			in:   "111\t222",
			want: "111\n",
		},
		{
			desc: "Cuts are done linewise",
			args: []string{"-f", "2", "-d", ","},
			in:   "111,222\n111,222\n111,",
			want: "222\n222\n\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var out strings.Builder
			in := strings.NewReader(tC.in)
			cfg, err := cc.NewConfig(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err := cc.Cut(in, &out, cfg); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			got := out.String()
			if got != tC.want {
				t.Fatalf("Bad cut - got: %q, want: %q", got, tC.want)
			}
		})
	}
}

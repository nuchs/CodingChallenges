package main_test

import (
	"testing"

	wc "github.com/nuchs/ccwc"
)

func TestBadFlag(t *testing.T) {
	_, got := wc.LoadSpec([]string{"-bad"})
	if got == nil {
		t.Fatalf("Got nil but wanted error")
	}
}

func TestFlags(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want wc.Spec
	}{
		{
			desc: "count characters",
			args: []string{"-c", "file"},
			want: wc.Spec{Source: "file", Bytes: true},
		},
		{
			desc: "count runes",
			args: []string{"-m", "file"},
			want: wc.Spec{Source: "file", MultiByte: true},
		},
		{
			desc: "count words",
			args: []string{"-w", "file"},
			want: wc.Spec{Source: "file", Words: true},
		},
		{
			desc: "count lines",
			args: []string{"-l", "file"},
			want: wc.Spec{Source: "file", Lines: true},
		},
		{
			desc: "implicit settings",
			args: []string{"file"},
			want: wc.Spec{
				Source: "file",
				Bytes:  true,
				Words:  true,
				Lines:  true,
			},
		},
		{
			desc: "Explicitly everything",
			args: []string{"-c", "-m", "-w", "-l", "file"},
			want: wc.Spec{
				Source:    "file",
				Bytes:     true,
				MultiByte: true,
				Words:     true,
				Lines:     true,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := wc.LoadSpec(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error loading spec: %s", err)
			}
			if got != tC.want {
				t.Fatalf("Bad spec: got:%+v, want:%+v", got, tC.want)
			}
		})
	}
}

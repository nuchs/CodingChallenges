package main_test

import (
	"reflect"
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
			args: []string{"-c"},
			want: wc.Spec{Sources: []string{"stdin"}, Bytes: true},
		},
		{
			desc: "count runes",
			args: []string{"-m"},
			want: wc.Spec{Sources: []string{"stdin"}, MultiByte: true},
		},
		{
			desc: "count words",
			args: []string{"-w"},
			want: wc.Spec{Sources: []string{"stdin"}, Words: true},
		},
		{
			desc: "count lines",
			args: []string{"-l"},
			want: wc.Spec{Sources: []string{"stdin"}, Lines: true},
		},
		{
			desc: "multi file",
			args: []string{"-l", "file1", "file2", "file3"},
			want: wc.Spec{Sources: []string{"file1", "file2", "file3"}, Lines: true},
		},
		{
			desc: "implicit settings",
			args: []string{},
			want: wc.Spec{
				Sources: []string{"stdin"},
				Bytes:   true,
				Words:   true,
				Lines:   true,
			},
		},
		{
			desc: "Explicitly everything",
			args: []string{"-c", "-m", "-w", "-l", "file1", "file2"},
			want: wc.Spec{
				Sources:   []string{"file1", "file2"},
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
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad spec: got:%+v, want:%+v", got, tC.want)
			}
		})
	}
}

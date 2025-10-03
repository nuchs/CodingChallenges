package main_test

import (
	"strings"
	"testing"

	wc "github.com/nuchs/ccwc"
)

func TestBytes(t *testing.T) {
	spec := wc.Spec{Bytes: true}
	testCases := []struct {
		desc string
		file string
		want int
	}{
		{desc: "empty", file: "", want: 0},
		{desc: "single", file: "a", want: 1},
		{desc: "one line", file: "I am big", want: 8},
		{desc: "one line+newline", file: "I am big\n", want: 9},
		{
			desc: "multiline",
			file: "Tests whisper softly,\n`got` drifts far from `want` again,\nred leaves fill the log.",
			want: 82,
		},
		{desc: "multibyte character", file: "§", want: 2},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := wc.Count(spec, strings.NewReader(tC.file))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got.Bytes != tC.want {
				t.Fatalf("Got %v bytes but want %v bytes", got.Bytes, tC.want)
			}
		})
	}
}

func TestRunes(t *testing.T) {
	spec := wc.Spec{MultiByte: true}
	testCases := []struct {
		desc string
		file string
		want int
	}{
		{desc: "empty", file: "", want: 0},
		{desc: "single", file: "a", want: 1},
		{desc: "multibyte character", file: "§", want: 1},
		{desc: "mixed one line", file: "Hello   world\tΣ😊 café\u00A0", want: 22},
		{desc: "mixed multiline", file: "L1 αβ  \nL2😊\r\nL3\tcafé ", want: 21},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := wc.Count(spec, strings.NewReader(tC.file))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got.Runes != tC.want {
				t.Fatalf("Got %v runes but want %v runes", got.Runes, tC.want)
			}
		})
	}
}

func TestWords(t *testing.T) {
	spec := wc.Spec{Words: true}
	testCases := []struct {
		desc string
		file string
		want int
	}{
		{desc: "empty", file: "", want: 0},
		{desc: "all whitespace", file: " \t \r \n  \r\n   ", want: 0},
		{desc: "one word", file: "word", want: 1},
		{desc: "multi word", file: "I am in the industry", want: 5},
		{desc: "alt whitespace", file: "I   am\tin\rthe\nindustry", want: 5},
		{
			desc: "multi word+newline",
			file: "You are in the industry too\n",
			want: 6,
		},
		{
			desc: "multiline",
			file: "I'm ready\nPromotion!\nI am ready? No...",
			want: 7,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := wc.Count(spec, strings.NewReader(tC.file))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got.Words != tC.want {
				t.Fatalf("Got %v words but want %v words", got.Words, tC.want)
			}
		})
	}
}

func TestLines(t *testing.T) {
	spec := wc.Spec{Lines: true}
	testCases := []struct {
		desc string
		file string
		want int
	}{
		{desc: "Empty", file: "", want: 0},
		{desc: "one line", file: "hello\n", want: 1},
		{desc: "multiple lines", file: "hello\ncruel\nworld\n", want: 3},
		{desc: "no trailing newline", file: "I'm still here", want: 1},
		{desc: "blanklines", file: "\n\n\n\n", want: 4},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := wc.Count(spec, strings.NewReader(tC.file))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got.Lines != tC.want {
				t.Fatalf("Got %v lines but want %v lines", got.Lines, tC.want)
			}
		})
	}
}

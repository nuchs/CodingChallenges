package main_test

import (
	"strings"
	"testing"

	jp "github.com/nuchs/ccjp"
)

func TestParseMinimal(t *testing.T) {
	testCases := []struct {
		desc string
		data string
	}{
		{desc: "object", data: "{}"},
		{desc: "array", data: "[]"},
		{desc: "null", data: "null"},
		{desc: "true", data: "true"},
		{desc: "false", data: "false"},
		{desc: "number", data: "0"},
		{desc: "string", data: "\"\""},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := jp.NewParser(strings.NewReader(tC.data))
			if err := p.Parse(); err != nil {
				t.Fatalf("Unexpected parse error: %v", err)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	testCases := []struct {
		desc string
		data string
	}{
		{desc: "single", data: `[1.0e-4]`},
		{desc: "multi", data: `["bacon", "egg", "sausage"]`},
		{desc: "mixed", data: `["a", true, false, null, 1, 2.0, 3e+1, 4E-2]`},
		{desc: "subarray", data: `[[], [5.1e-4], [6.2E+5, 7]]`},
		{desc: "subobject", data: `[{}, {"a":1, "n":"eep"}]`},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := jp.NewDebugParser(strings.NewReader(tC.data), true)
			if err := p.Parse(); err != nil {
				t.Fatalf("Unexpected parse error: %v", err)
			}
		})
	}
}

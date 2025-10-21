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
			p := jp.NewParser(strings.NewReader(tC.data))
			if err := p.Parse(); err != nil {
				t.Fatalf("Unexpected parse error: %v", err)
			}
		})
	}
}

func TestParseObject(t *testing.T) {
	testCases := []struct {
		desc string
		data string
	}{
		{desc: "empty object", data: "{}"},
		{desc: "single key", data: `{"key":"value"}`},
		{desc: "subarrary", data: `{"key":[1, 2, 3]}`},
		{desc: "subobject", data: `{"k1":{ "k2": {} }}`},
		{
			desc: "multi key",
			data: `{"a":"v", "b":1, "c":true, "d":null, "e":false}`,
		},
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

func TestBadJson(t *testing.T) {
	testCases := []struct {
		desc string
		data string
		err  string
	}{
		{
			desc: "Unterminated array",
			data: "[1,2",
			err:  "Parse failure: malformed array, expected ']'",
		},
		{
			desc: "Unterminated object",
			data: `{ "k":"v" `,
			err:  "Parse failure: malformed object, expected '}'",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := jp.NewParser(strings.NewReader(tC.data))
			err := p.Parse()
			if err == nil {
				t.Fatalf("Got no error, wanted: %s", err)
			}
			if !strings.HasPrefix(err.Error(), tC.err) {
				t.Fatalf("Wrong error, got '%s', want '%s'", err, tC.err)
			}
		})
	}
}

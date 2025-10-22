package main_test

import (
	"reflect"
	"strings"
	"testing"

	huff "github.com/nuchs/cc/huff"
)

func TestFrequencyCounting(t *testing.T) {
	testCases := []struct {
		data string
		want map[rune]int
	}{
		{
			data: "",
			want: map[rune]int{},
		},
		{
			data: "aaa",
			want: map[rune]int{'a': 3},
		},
		{
			data: "ababaa",
			want: map[rune]int{'a': 4, 'b': 2},
		},
		{
			data: "😃😃😃😃",
			want: map[rune]int{'😃': 4},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.data, func(t *testing.T) {
			r := strings.NewReader(tC.data)
			got, err := huff.MakeFrequencyTable(r)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Wrong frequencies, got %+v, want %+v", got, tC.want)
			}
		})
	}
}

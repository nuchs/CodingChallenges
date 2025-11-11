package main

import (
	"reflect"
	"testing"
)

func TestMask(t *testing.T) {
	testCases := []struct {
		desc string
		size byte
		want uint16
	}{
		{desc: "one", size: 1, want: 0b1},
		{desc: "two", size: 2, want: 0b11},
		{desc: "five", size: 5, want: 0b11111},
		{desc: "eight", size: 8, want: 0b11111111},
		{desc: "eight", size: 8, want: 0b11111111},
		{desc: "eight", size: 11, want: 0b111_11111111},
		{desc: "sixteen", size: 16, want: 0b11111111_11111111},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := bitmask(tC.size)
			if got != tC.want {
				t.Fatalf("Bad bit mask: got %b, want %b", got, tC.want)
			}
		})
	}
}

func TestWriting(t *testing.T) {
	a := entry{value: 'a', prefix: 0b0, len: 2}
	b := entry{value: 'b', prefix: 0b1, len: 2}
	c := entry{value: 'c', prefix: 0b110, len: 3}
	testCases := []struct {
		desc string
		data []entry
		want []byte
	}{
		{
			desc: "partial byte only",
			data: []entry{a, b, a},
			want: []byte{0b00_01_00_00, 2},
		},
		{
			desc: "full byte on boundry, same lengths",
			data: []entry{b, a, b, a},
			want: []byte{0b01_00_01_00, 0},
		},
		{
			desc: "full byte on boundry, different lengths",
			data: []entry{c, a, c},
			want: []byte{0b110_00_110, 0},
		},
		{
			desc: "across boundry in middle but end on boundry",
			data: []entry{c, c, c, a, c, a},
			want: []byte{0b110_110_11, 0b0_00_110_00, 0},
		},
		{
			desc: "no boundry matching",
			data: []entry{b, a, c, b},
			want: []byte{0b01_00_110_0, 0b1_0000000, 7},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sb := newShortBuffer([]byte{})
			for _, e := range tC.data {
				sb.write(e.prefix, e.len)
			}
			got := sb.close()
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad short buffer: got %+v, want%+v", got, tC.want)
			}
		})
	}
}

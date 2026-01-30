package encode

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func FuzzKraftMcMillan(f *testing.F) {
	const epsilon float64 = 1e-9
	f.Add([]byte{})
	f.Add([]byte("a"))
	f.Add([]byte("aaa"))
	f.Add([]byte("ab"))
	f.Add([]byte("ba"))
	f.Add([]byte("aabbcccc"))
	f.Add([]byte("aabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaeeaaaaadaaaaca"))

	f.Fuzz(func(t *testing.T, bs []byte) {
		if len(bs) == 0 {
			t.Log("Skip empty slice")
			return
		}
		et := newEncodingTable(bs)
		sum := 0.0
		for _, v := range et {
			sum += math.Pow(2, -float64(v.len))
		}
		if sum > (1 + epsilon) {
			t.Fatalf("Kraft McMillan violated: test %q, sum = %g", bs, sum)
		}
	})
}

func FuzzDuplicateCodes(f *testing.F) {
	f.Add([]byte("ab"))
	f.Add([]byte("ba"))
	f.Add([]byte("bac"))
	f.Add([]byte("babcac"))
	f.Add([]byte("abcdefghijklmnopqrstuvwxyz0123456789"))
	f.Fuzz(func(t *testing.T, bs []byte) {
		if len(bs) == 0 {
			t.Log("Skip empty slice")
			return
		}
		et := newEncodingTable(bs)
		for _, b := range bs {
			if _, ok := et[b]; !ok {
				t.Fatalf(
					"No code in table for %v, must have been overwritten by a duplicate",
					b,
				)
			}
		}
	})
}

func FuzzCodeLength(f *testing.F) {
	f.Add([]byte("abc"))
	f.Add([]byte("aaaab"))
	f.Add([]byte("baaaa"))
	f.Add([]byte("babcac"))
	f.Add([]byte("the quick brown fox jumped over the lazy cat. cat cat cat cat"))
	f.Fuzz(func(t *testing.T, bs []byte) {
		if len(bs) == 0 {
			t.Log("Skip empty slice")
			return
		}
		et := newEncodingTable(bs)
		entries := et.sortTable()
		last := entries[0]
		t.Log("\n" + et.String())
		for _, curr := range entries[1:] {
			if curr.freq > last.freq && curr.len != last.len {
				t.Fatalf(
					"prefix length should increase as frequency decreases: curr: %+v, last: %+v",
					curr,
					last,
				)
			}
			last = curr
		}
	})
}

func FuzzSharedPrefix(f *testing.F) {
	f.Add([]byte("ab"))
	f.Add([]byte("ba"))
	f.Add([]byte("bac"))
	f.Add([]byte("babcac"))
	f.Add([]byte("abcdefghijklmnopqrstuvwxyz0123456789"))
	f.Fuzz(func(t *testing.T, bs []byte) {
		if len(bs) == 0 {
			t.Log("Skip empty slice")
			return
		}
		et := newEncodingTable(bs)
		entries := et.sortTable()
		for i := 0; i < len(entries); i++ {
			for j := i + 1; j < len(entries); j++ {
				left := fmt.Sprintf("%0*b", entries[i].len, entries[i].prefix)
				right := fmt.Sprintf("%0*b", entries[j].len, entries[j].prefix)
				if strings.HasPrefix(right, left) {
					t.Fatalf("Shared prefix: %+v - %+v", entries[i], entries[j])
				}
			}
		}
	})
}

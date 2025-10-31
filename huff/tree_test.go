package main

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func FuzzKraftMcMillan(f *testing.F) {
	const epsilon float64 = 1e-9
	f.Add("")
	f.Add("a")
	f.Add("aaa")
	f.Add("ab")
	f.Add("ba")
	f.Add("aabbcccc")
	f.Add("😃😃😃a😃bb")

	f.Fuzz(func(t *testing.T, s string) {
		if s == "" {
			t.Log("Skip empty string")
			return
		}
		pt, err := generateEncodingTable(strings.NewReader(s))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		sum := 0.0
		for _, v := range pt {
			sum += math.Pow(2, float64(-v.len))
		}
		if sum > (1 + epsilon) {
			t.Fatalf("Kraft McMillan violated: test %q, sum = %g", s, sum)
		}
	})
}

func FuzzDuplicateCodes(f *testing.F) {
	f.Add("")
	f.Add("a")
	f.Add("ab")
	f.Add("ba")
	f.Add("bac")
	f.Add("babcac")
	f.Fuzz(func(t *testing.T, s string) {
		if s == "" {
			t.Log("Skip empty string")
			return
		}
		pt, err := generateEncodingTable(strings.NewReader(s))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		for _, c := range s {
			if _, ok := pt[c]; !ok {
				t.Fatalf(
					"No code in table for %v, must have been overwritten by a duplicate",
					c,
				)
			}
		}
	})
}

func TestPrefixTable(t *testing.T) {
	testCases := []struct {
		desc string
		tree node
		want prefixTable
	}{
		{
			desc: "single",
			tree: newLeaf('a', 1),
			want: prefixTable{'a': prefix{value: 0, len: 1}},
		},
		{
			desc: "two",
			tree: newNode(newLeaf('a', 1), newLeaf('b', 2)),
			want: prefixTable{
				'a': prefix{value: 0, len: 1},
				'b': prefix{value: 1, len: 1},
			},
		},
		{
			desc: "many",
			tree: newNode(
				newNode(newLeaf('c', 4), newLeaf('d', 4)),
				newNode(
					newNode(newLeaf('a', 2), newLeaf('b', 3)),
					newLeaf('e', 6),
				),
			),
			want: prefixTable{
				'a': prefix{value: 0b100, len: 3},
				'b': prefix{value: 0b101, len: 3},
				'c': prefix{value: 0b0, len: 2},
				'd': prefix{value: 0b1, len: 2},
				'e': prefix{value: 0b11, len: 2},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.tree.toPrefixTable()
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad prefix table: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func TestTreeBuilding(t *testing.T) {
	testCases := []struct {
		desc string
		freq map[rune]int
		want node
	}{
		{
			desc: "one node",
			freq: map[rune]int{'a': 1},
			want: newLeaf('a', 1),
		},
		{
			desc: "two node",
			freq: map[rune]int{'a': 1, 'b': 2},
			want: newNode(newLeaf('a', 1), newLeaf('b', 2)),
		},
		{
			desc: "more node",
			freq: map[rune]int{'a': 2, 'b': 3, 'd': 4, 'c': 4},
			want: newNode(
				newNode(newLeaf('a', 2), newLeaf('b', 3)),
				newNode(newLeaf('c', 4), newLeaf('d', 4)),
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := buildHuffmanTree(tC.freq)
			if !treeEqual(&got, &tC.want) {
				t.Fatalf("Wrong tree - got %+v, want %+v", got, tC.want)
			}
		})
	}
}

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
			got, err := makeFrequencyTable(r)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Wrong frequencies, got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func TestNodeWeight(t *testing.T) {
	testCases := []struct {
		desc string
		a    node
		b    node
		want int
	}{
		{
			desc: "two leaves",
			a:    newLeaf('a', 1),
			b:    newLeaf('b', 2),
			want: 3,
		},
		{
			desc: "left leaf, right node",
			a:    newLeaf('a', 1),
			b:    newNode(newLeaf('b', 4), newLeaf('c', 5)),
			want: 10,
		},
		{
			desc: "left node, right leaf",
			a:    newNode(newLeaf('b', 2), newLeaf('c', 8)),
			b:    newLeaf('a', 1),
			want: 11,
		},
		{
			desc: "two nodes",
			a:    newNode(newLeaf('a', 2), newLeaf('c', 8)),
			b:    newNode(newLeaf('b', 2), newLeaf('d', 8)),
			want: 20,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := newNode(tC.a, tC.b).weight
			if got != tC.want {
				t.Fatalf(
					"Node weight should be sum of child weights. Got %v, want %v",
					got,
					tC.want,
				)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	testCases := []struct {
		desc  string
		left  node
		right node
		want  int
	}{
		{
			desc:  "left leaf lighter",
			left:  newLeaf('a', 1),
			right: newLeaf('b', 2),
			want:  -1,
		},
		{
			desc:  "right leaf lighter",
			left:  newLeaf('a', 2),
			right: newLeaf('b', 1),
			want:  1,
		},
		{
			desc:  "left node lighter",
			left:  newNode(newLeaf('a', 2), newLeaf('c', 7)),
			right: newNode(newLeaf('b', 2), newLeaf('d', 8)),
			want:  -1,
		},
		{
			desc:  "right node lighter",
			left:  newNode(newLeaf('b', 2), newLeaf('d', 8)),
			right: newNode(newLeaf('a', 2), newLeaf('c', 7)),
			want:  1,
		},
		{
			desc:  "left leaf lighter than right node",
			left:  newLeaf('a', 1),
			right: newNode(newLeaf('b', 2), newLeaf('d', 8)),
			want:  -1,
		},
		{
			desc:  "right leaf lighter than left node",
			left:  newNode(newLeaf('b', 2), newLeaf('d', 8)),
			right: newLeaf('a', 1),
			want:  1,
		},
		{
			desc:  "left node lighter than right leaf",
			left:  newNode(newLeaf('b', 1), newLeaf('d', 1)),
			right: newLeaf('a', 3),
			want:  -1,
		},
		{
			desc:  "right node lighter than left leaf",
			left:  newLeaf('a', 5),
			right: newNode(newLeaf('b', 2), newLeaf('d', 2)),
			want:  1,
		},
		{
			desc:  "ordinal breaks leaf tie",
			left:  newLeaf('a', 1),
			right: newLeaf('b', 1),
			want:  -1,
		},
		{
			desc:  "ordinal breaks leaf/node tie",
			left:  newLeaf('e', 3),
			right: newNode(newLeaf('b', 1), newLeaf('d', 2)),
			want:  1,
		},
		{
			desc:  "ordinal breaks node tie",
			left:  newNode(newLeaf('a', 1), newLeaf('z', 2)),
			right: newNode(newLeaf('b', 1), newLeaf('d', 2)),
			want:  -1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			w := cmpNode(tC.left, tC.right)
			got := normaliseWeight(w)
			if got != tC.want {
				t.Fatalf("left < right? got %v, want %v", got, tC.want)
			}
		})
	}
}

func normaliseWeight(weight int) int {
	if weight < 0 {
		return -1
	}
	if weight > 0 {
		return 1
	}
	return 0
}

func treeEqual(a, b *node) bool {
	switch {
	case a.isLeaf != b.isLeaf,
		a.weight != b.weight,
		a.value != b.value,
		a.minSymbol != b.minSymbol,
		a.left != nil && b.left == nil,
		a.left == nil && b.left != nil,
		a.right != nil && b.right == nil,
		a.right == nil && b.right != nil:
		return false
	case a.left != nil && a.right != nil:
		return treeEqual(a.left, b.left) && treeEqual(a.right, b.right)
	case a.left != nil:
		return treeEqual(a.left, b.left)
	case a.right != nil:
		return treeEqual(a.right, b.right)
	default:
		return true
	}
}

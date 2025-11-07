package main

import (
	"fmt"
	"slices"
	"strings"
)

const maxCodeLength = 16

// Huffman encoding details for a specific byte pattern
type entry struct {
	value  byte   // The value being encoded
	prefix uint16 // The prefix it is encoded to
	len    byte   // The length of the prefix
	freq   int    // How often the value occurred in the input
}

// The huffman encoding table for the input
type encodeTable map[byte]entry

func NewEncodingTable(data []byte) encodeTable {
	freq := newFrequencyTable(data)
	ht := newHuffmanTree(freq)
	entries, lens := generateEntries(&ht)
	capLengths(entries, lens)

	return assignCodes(entries)
}

func (et encodeTable) String() string {
	es := et.sortTable()
	var sb strings.Builder
	fmt.Fprintf(&sb, "  value  |      prefix      | len |  freq\n")
	fmt.Fprintf(&sb, "------------------------------------------\n")
	for _, e := range es {
		fmt.Fprintf(
			&sb,
			"%08b |%*s%0*b | %3d | %5d\n",
			e.value,
			maxCodeLength-e.len+1,
			" ",
			e.len,
			e.prefix,
			e.len,
			e.freq,
		)
	}
	return sb.String()
}

func (et encodeTable) sortTable() []entry {
	entries := make([]entry, 0, len(et))
	for _, v := range et {
		entries = append(entries, v)
	}

	sortEntries(entries)
	return entries
}

func sortEntries(entries []entry) {
	slices.SortFunc(entries, func(a, b entry) int {
		if a.len != b.len {
			return int(a.len) - int(b.len)
		}
		if a.freq != b.freq {
			return b.freq - a.freq
		}
		return int(a.value) - int(b.value)
	})
}

func generateEntries(ht *node) ([]entry, []byte) {
	entries := make([]entry, 0, 256)
	lens := make([]byte, maxCodeLength+2) // 1 based array, + 1 for overflow

	var dfs func(d byte, n *node)
	dfs = func(d byte, n *node) {
		if n.isLeaf {
			e := entry{
				len:   max(1, d), // Need max to handle edge case of only one symbol
				freq:  n.weight,
				value: n.value,
			}
			// record everything over the max length together
			lens[min(e.len, maxCodeLength+1)]++
			entries = append(entries, e)
			return
		}
		if n.left != nil {
			dfs(d+1, n.left)
		}
		if n.right != nil {
			dfs(d+1, n.right)
		}
	}
	dfs(0, ht)
	sortEntries(entries)

	return entries, lens
}

func capLengths(entries []entry, lens []byte) {
	over := lens[maxCodeLength+1]
	if over == 0 {
		return
	}
	for i := maxCodeLength - 1; i > 0; i-- {
		if over == 0 {
			break
		}
		if lens[i] == 0 {
			continue
		}
		split := min(over, lens[i])
		lens[i] -= split
		lens[i+1] += 2 * split
		over -= split
		i = min(i+1, maxCodeLength-1)

	}

	var curr int
	for i := 1; i < maxCodeLength+1; i++ {
		for j := lens[i]; j > 0; j-- {
			entries[curr].len = byte(i)
			curr++
		}
	}
}

func assignCodes(entries []entry) encodeTable {
	et := make(encodeTable, len(entries))
	var code uint16
	var currLen byte

	for _, e := range entries {
		if currLen != e.len {
			code <<= e.len - currLen
			currLen = e.len
		}
		e.prefix = code
		code++
		et[e.value] = e
	}
	return et
}

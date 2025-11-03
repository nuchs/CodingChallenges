package main

import (
	"fmt"
	"slices"
	"strings"
)

// Record how often a particular byte pattern occurs in the input
type frequencyTable map[byte]int

// Huffman encoding details for a specific byte pattern
type entry struct {
	value  byte   // The value being encoded
	prefix uint16 // The prefix it is encoded to
	len    byte   // The length of the prefix
	freq   int    // How often the value occurred in the input
}

// The huffman encoding table for the input
type encodeTable map[byte]entry

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
			17-e.len,
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

	slices.SortFunc(entries, func(a, b entry) int {
		if a.len == b.len {
			return b.freq - a.freq
		}
		return int(a.len) - int(b.len)
	})

	return entries
}

// A node in a huffman tree
type node struct {
	weight      int   // The frequency of all the symbols in the subtree
	minSymbol   byte  // The lexically smallest symbol in the subtree (used for sorting)
	isLeaf      bool  // Is this node a leaf
	value       byte  // For leaves, this is the symbol being encoded
	left, right *node // For non leaves the left and right child nodes
}

// Creates a new leaf node for the specified symbol
func newLeaf(value byte, weight int) node {
	return node{
		weight:    weight,
		minSymbol: value,
		isLeaf:    true,
		value:     value,
	}
}

// Create a new internal node to be the parent of two existing nodes
func newNode(a, b node) node {
	left := a
	right := b
	if cmpNode(a, b) > 0 {
		left = b
		right = a
	}

	return node{
		weight:    a.weight + b.weight,
		minSymbol: minSymbol(a, b),
		isLeaf:    false,
		left:      &left,
		right:     &right,
	}
}

// Convert a huffman tree to an encoding table
func (n *node) toEncodeTable() encodeTable {
	tab := encodeTable{}
	var pLen byte
	var pfx uint16

	if n.isLeaf {
		pLen++
	}

	fillPrefixes(n, tab, pLen, pfx)

	return tab
}

func fillPrefixes(n *node, tab encodeTable, pLen byte, pfx uint16) {
	if n == nil {
		return
	}
	if n.isLeaf {
		tab[n.value] = entry{value: n.value, prefix: pfx, len: pLen, freq: n.weight}
		return
	}
	pLen++
	left := pfx << 1
	right := left + 1
	fillPrefixes(n.left, tab, pLen, left)
	fillPrefixes(n.right, tab, pLen, right)
}

func buildHuffmanTree(ft frequencyTable) node {
	trees := make([]node, 0, len(ft))

	for r, f := range ft {
		trees = append(trees, newLeaf(r, f))
	}

	for len(trees) > 1 {
		slices.SortFunc(trees, cmpNode)
		n := newNode(trees[0], trees[1])
		trees[1] = n
		trees = trees[1:]
	}

	return trees[0]
}

func minSymbol(a, b node) byte {
	if a.minSymbol == b.minSymbol {
		panic(fmt.Sprintf(
			"Two nodes should not contain same symbol: %v",
			a.minSymbol,
		))
	}
	if a.minSymbol < b.minSymbol {
		return a.minSymbol
	}
	return b.minSymbol
}

func cmpNode(a, b node) int {
	if a.weight == b.weight {
		return int(a.minSymbol) - int(b.minSymbol)
	}
	return a.weight - b.weight
}

func generateEncodingTable(data []byte) encodeTable {
	freq := makeFrequencyTable(data)
	t := buildHuffmanTree(freq)

	return t.toEncodeTable()
}

func makeFrequencyTable(data []byte) frequencyTable {
	counts := frequencyTable{}
	for _, b := range data {
		if _, ok := counts[b]; !ok {
			counts[b] = 0
		}
		counts[b]++
	}

	return counts
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

type prefix struct {
	value int
	len   int
}

type prefixTable map[rune]prefix

type node struct {
	weight      int
	minSymbol   rune
	isLeaf      bool
	value       rune
	left, right *node
}

func newLeaf(value rune, weight int) node {
	return node{
		weight:    weight,
		minSymbol: value,
		isLeaf:    true,
		value:     value,
	}
}

func newNode(a, b node) node {
	left := a
	right := b
	if cmpNode(a, b) >= 0 {
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

func (n *node) toPrefixTable() prefixTable {
	tab := make(map[rune]prefix)
	var plen, pfx int

	if n.isLeaf {
		plen++
	}

	fillPrefixes(n, tab, plen, pfx)

	return tab
}

func fillPrefixes(n *node, tab map[rune]prefix, plen, pfx int) {
	if n == nil {
		return
	}
	if n.isLeaf {
		tab[n.value] = prefix{value: pfx, len: plen}
		return
	}
	plen++
	left := pfx << 1
	right := left + 1
	fillPrefixes(n.left, tab, plen, left)
	fillPrefixes(n.right, tab, plen, right)
}

func generateEncodingTable(data io.Reader) (prefixTable, error) {
	freq, err := makeFrequencyTable(data)
	if err != nil {
		return nil, err
	}
	t := buildHuffmanTree(freq)

	return t.toPrefixTable(), nil
}

func makeFrequencyTable(data io.Reader) (map[rune]int, error) {
	counts := make(map[rune]int)
	read := bufio.NewReader(data)
	for {
		r, _, err := read.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to generate frequency table: %w", err)
		}
		if _, ok := counts[r]; !ok {
			counts[r] = 0
		}
		counts[r]++
	}

	return counts, nil
}

func buildHuffmanTree(freq map[rune]int) node {
	trees := make([]node, 0, len(freq))

	for r, f := range freq {
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

func minSymbol(a, b node) rune {
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
		return int(a.minSymbol - b.minSymbol)
	}
	return a.weight - b.weight
}

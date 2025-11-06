package main

import (
	"fmt"
	"slices"
)

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

func newHuffmanTree(ft frequencyTable) node {
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

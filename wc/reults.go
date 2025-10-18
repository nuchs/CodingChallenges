package main

import (
	"fmt"
	"math"
)

const DefaultWidth = 6

type Results struct {
	fw     FieldWidths
	counts []Counts
	totals Counts
}

type FieldWidths struct {
	Byte int
	Rune int
	Word int
	Line int
}

func NewResults(size int) Results {
	return Results{
		fw: FieldWidths{
			Byte: DefaultWidth,
			Rune: DefaultWidth,
			Word: DefaultWidth,
			Line: DefaultWidth,
		},
		counts: make([]Counts, 0, size),
		totals: Counts{Src: "total"},
	}
}

func (r *Results) Add(count Counts) {
	r.counts = append(r.counts, count)
	if digitCount(count.Bytes) >= r.fw.Byte {
		r.fw.Byte = digitCount(count.Bytes) + 1
	}
	if digitCount(count.Runes) >= r.fw.Rune {
		r.fw.Rune = digitCount(count.Runes) + 1
	}
	if digitCount(count.Words) >= r.fw.Word {
		r.fw.Word = digitCount(count.Words) + 1
	}
	if digitCount(count.Lines) >= r.fw.Line {
		r.fw.Line = digitCount(count.Lines) + 1
	}

	r.totals.Bytes += count.Bytes
	r.totals.Runes += count.Runes
	r.totals.Words += count.Words
	r.totals.Lines += count.Lines
}

func (r *Results) Print(spec Spec) {
	for _, res := range r.counts {
		fmt.Println(res.Format(spec, r.fw))
	}

	if len(r.counts) > 1 {
		fmt.Println(r.totals.Format(spec, r.fw))
	}
}

func digitCount(num int) int {
	return int(math.Log10(float64(num))) + 1
}

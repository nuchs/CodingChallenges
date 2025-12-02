package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrBadlyFormattedSpread  = errors.New("spreads must be in the format 'a-b'")
	ErrNonNumericSpreadBound = errors.New("spread boundries must be numeric")
	ErrZeroSpreadBound       = errors.New("spread ranges cannot include zero")
	ErrBadSpreadOrder        = errors.New("start of spread is greater than end")
)

type Spreads []Spread

type Spread struct {
	Min int
	Max int
}

func newSpread(spec string) (Spread, error) {
	var s Spread
	parts := strings.Split(spec, "-")
	if len(parts) > 2 {
		return s, ErrBadlyFormattedSpread
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return Spread{}, ErrNonNumericSpreadBound
	}
	end := start
	if len(parts) > 1 {
		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return Spread{}, ErrNonNumericSpreadBound
		}
	}
	if start == 0 || end == 0 {
		return Spread{}, ErrZeroSpreadBound
	}
	if start > end {
		return Spread{}, ErrBadSpreadOrder
	}

	return Spread{Min: start, Max: end}, nil
}

func NewSpreads(spec string) (Spreads, error) {
	var r Spreads
	for v := range strings.SplitSeq(spec, ",") {
		s, err := newSpread(v)
		if err != nil {
			return r, fmt.Errorf("invalid range specification %q: %w", spec, err)
		}
		r = append(r, s)
	}

	slices.SortFunc(r, func(a, b Spread) int {
		return a.Min - b.Min
	})

	r = mergeSpreads(r)

	return r, nil
}

func (r Spreads) Includes(i int) bool {
	for _, s := range r {
		if i >= s.Min && i <= s.Max {
			return true
		}
	}
	return false
}

func mergeSpreads(r Spreads) Spreads {
	if len(r) == 0 {
		return r
	}
	m := append(Spreads(nil), r[0])
	curr := 0

	for next := 1; next < len(r); next++ {
		if (m[curr].Min >= r[next].Min && m[curr].Min <= r[next].Max) ||
			(m[curr].Max >= r[next].Min && m[curr].Max <= r[next].Max) ||
			(r[next].Min >= m[curr].Min && r[next].Max <= m[curr].Max) ||
			(r[next].Max == m[curr].Min-1) || (r[next].Min == m[curr].Max+1) {
			m[curr].Min = min(m[curr].Min, r[next].Min)
			m[curr].Max = max(m[curr].Max, r[next].Max)
		} else {
			m = append(m, r[next])
			curr++
		}
	}

	return m
}

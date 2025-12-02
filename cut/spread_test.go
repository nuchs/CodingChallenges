package main_test

import (
	"testing"

	cc "github.com/nuchs/cc/cut"
)

func TestRangeInclusion(t *testing.T) {
	testCases := []struct {
		desc   string
		values []int
		in     string
	}{
		{
			desc:   "Singular",
			values: []int{1},
			in:     "1",
		},
		{
			desc:   "Repeated Singular",
			values: []int{1, 2, 5},
			in:     "1,2,5",
		},
		{
			desc:   "Spread includes bounds",
			values: []int{2, 3, 4, 5},
			in:     "2-5",
		},
		{
			desc:   "Mixed",
			values: []int{1, 3, 4, 5, 9, 10, 13},
			in:     "1,3-5,9-10,13",
		},
		{
			desc:   "Overlapping/adjacent",
			values: []int{1, 2, 3, 7, 8, 9, 11, 12, 13, 14},
			in:     "1-3,2,4-8,8-10,13-14,11-13",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := cc.NewSpreads(tC.in)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			for _, v := range tC.values {
				if !r.Includes(v) {
					t.Logf("%+v should include %d but doesn't", r, v)
					t.Fail()
				}
			}
		})
	}
}

func TestRangeExclusion(t *testing.T) {
	testCases := []struct {
		desc   string
		values []int
		in     string
	}{
		{
			desc:   "Singular",
			values: []int{-1, 0, 1, 3},
			in:     "2",
		},
		{
			desc:   "Repeated Singular",
			values: []int{1, 3, 4, 6},
			in:     "2,5",
		},
		{
			desc:   "Spread",
			values: []int{1, 6},
			in:     "2-5",
		},
		{
			desc:   "mixed",
			values: []int{1, 6, 7, 9, 21},
			in:     "2-5,8,10,11-20",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := cc.NewSpreads(tC.in)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			for _, v := range tC.values {
				if r.Includes(v) {
					t.Logf("%+v should not include %d but does", r, v)
					t.Fail()
				}
			}
		})
	}
}

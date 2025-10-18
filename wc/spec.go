package main

import (
	"flag"
	"fmt"
	"strings"
)

type Spec struct {
	Bytes     bool
	MultiByte bool
	Words     bool
	Lines     bool
	Sources   []string
}

func LoadSpec(args []string) (Spec, error) {
	var spec Spec
	parser := flag.NewFlagSet("ccwc", flag.ContinueOnError)
	var usage strings.Builder
	parser.Usage = func() {
		parser.SetOutput(&usage)
		parser.PrintDefaults()
	}
	parser.BoolVar(
		&spec.Bytes,
		"c",
		false,
		"count the number of bytes in a file",
	)
	parser.BoolVar(
		&spec.MultiByte,
		"m",
		false,
		"count the number of runes in a file",
	)
	parser.BoolVar(
		&spec.Words,
		"w",
		false,
		"count the number of word in a file",
	)
	parser.BoolVar(
		&spec.Lines,
		"l",
		false,
		"count the number of lines in a file",
	)
	if err := parser.Parse(args); err != nil {
		return Spec{}, fmt.Errorf(
			"failed to parse arguments: %w\n%s",
			err,
			usage.String(),
		)
	}

	if !spec.Bytes && !spec.MultiByte && !spec.Words && !spec.Lines {
		spec.Bytes = true
		spec.Words = true
		spec.Lines = true
	}

	spec.Sources = []string{"stdin"}
	tail := parser.Args()
	if len(tail) >= 1 {
		spec.Sources = tail
	}

	return spec, nil
}

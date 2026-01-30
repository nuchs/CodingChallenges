package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Actions huff can perform
type Action int

const (
	NotSet Action = iota

	// Generate the encoding table for a file and print it
	EncodePrint

	// Huffman encode a file
	Encoding

	// Read the encoding table from a file and print it
	DecodePrint

	// Decode a huffman encoded file
	Decoding
)

func (a Action) String() string {
	switch a {
	case NotSet:
		return "Not set"
	case EncodePrint:
		return "Encoding Print"
	case Encoding:
		return "Encode"
	case DecodePrint:
		return "Decoding Print"
	case Decoding:
		return "Decode"
	default:
		panic("bad action!")
	}
}

var (
	ErrActionNotSet = errors.New("action has not been set")
	ErrNoInputFile  = errors.New("you must specify a file to process")
	ErrNoOutputFile = errors.New(
		"for encode/decode actions you must specify an output file",
	)
)

// The Spec contains the user provided details on the job to be performed
type Spec struct {
	Action Action
	In     string
	Out    string
}

// Create a new job specification from the provided command line arguments
func NewSpec(args []string) (*Spec, error) {
	var usage strings.Builder
	usage.WriteString("A huffman encoder/decoder\n\nUsage: huff [OPTIONS]\n")
	spec := &Spec{}
	parser := configureParser(&usage, spec)

	err := validateArgs(parser, args, spec)
	if err != nil {
		parser.PrintDefaults()
		return nil, fmt.Errorf(
			"%w\n%s",
			err,
			usage.String(),
		)
	}

	return spec, nil
}

func validateArgs(
	parser *flag.FlagSet,
	args []string,
	spec *Spec,
) error {
	if err := parser.Parse(args); err != nil {
		return err
	}

	if spec.Action == NotSet {
		return ErrActionNotSet
	}
	if spec.In == "" {
		return ErrNoInputFile
	}
	if spec.Out == "" && spec.Action != EncodePrint {
		return ErrNoOutputFile
	}
	return nil
}

func configureParser(usage *strings.Builder, spec *Spec) *flag.FlagSet {
	parser := flag.NewFlagSet("huff", flag.ContinueOnError)

	parser.SetOutput(usage)
	parser.StringVar(&spec.In, "i", "", "File to process")
	parser.StringVar(&spec.Out, "o", "", "File to write output to")
	parser.Func(
		"a",
		"action to perform <e|ep|d|dp> (print header/encode/decode)",
		func(a string) error {
			switch a {
			case "ep":
				spec.Action = EncodePrint
			case "e":
				spec.Action = Encoding
			case "dp":
				spec.Action = DecodePrint
			case "d":
				spec.Action = Decoding
			default:
				return fmt.Errorf("unrecognised option: %s", a)
			}
			return nil
		},
	)

	return parser
}

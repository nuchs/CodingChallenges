// A program for generating length limited huffman codes for an input and using
// them to encode an output, can also be used to decode the files after
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nuchs/cc/huff/decode"
	"github.com/nuchs/cc/huff/encode"
)

func main() {
	spec, err := NewSpec(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	switch spec.Action {
	case EncodePrint:
		err = printEncodingMetadata(spec)
	case Encoding:
		err = encodeFile(spec)
	case DecodePrint:
		err = printDecodingMetadata(spec)
	case Decoding:
		err = decodeFile(spec)
	default:
		err = fmt.Errorf("unrecognised action: %s", spec.Action)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
}

func printEncodingMetadata(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for printing: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)
	return encode.PrintEncodeTable(ifp, os.Stdout)
}

func encodeFile(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for encoding: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)

	ofp, err := os.OpenFile(spec.Out, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %s", spec.Out, err)
	}
	defer closeFile(spec.Out, ofp)

	return encode.Encode(ifp, ofp)
}

func printDecodingMetadata(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for printing: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)
	return decode.PrintDecodeTable(ifp, os.Stdout)
}

func decodeFile(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for decoding: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)

	ofp, err := os.OpenFile(spec.Out, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %s", spec.Out, err)
	}
	defer closeFile(spec.Out, ofp)

	return decode.Decode(ifp, ofp)
}

func closeFile(filename string, c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close %s: %s", filename, err)
	}
}

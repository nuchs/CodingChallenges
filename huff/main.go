package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	spec, err := NewSpec(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	switch spec.Action {
	case Printing:
		err = printMetadata(spec)
	case Encoding:
		err = encodeFile(spec)
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

func printMetadata(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for printing: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)
	return PrintHeaders(ifp, os.Stdout)
}

func encodeFile(spec *Spec) error {
	ifp, err := os.Open(spec.In)
	if err != nil {
		return fmt.Errorf("failed to open %q for encoding: %s", spec.In, err)
	}
	defer closeFile(spec.In, ifp)

	ofp, err := os.OpenFile(spec.Out, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("Failed to open %q for writing: %s", spec.Out, err)
	}
	defer closeFile(spec.Out, ofp)

	return Encode(ifp, ofp)
}

func decodeFile(_ *Spec) error {
	return nil
}

func closeFile(filename string, c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close %s: %s", filename, err)
	}
}

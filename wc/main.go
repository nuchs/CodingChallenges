// Code Challenge: Build a wc clone
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	spec, err := LoadSpec(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rd, err := openSource(&spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open source for reading: %v", err)
		os.Exit(2)
	}
	defer func() {
		if err := rd.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close source %+v: %s", spec, err)
		}
	}()

	counts, err := Count(spec, rd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to complete request %+v: %v", spec, err)
		os.Exit(3)
	}

	fmt.Println(counts.Format(spec))
}

func openSource(spec *Spec) (io.ReadCloser, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		spec.Source = "stdin"
		return os.Stdin, nil
	}

	f, err := os.Open(spec.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", spec.Source, err)
	}
	return f, nil
}

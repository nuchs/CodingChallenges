package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	cfg, err := NewConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	in, err := openInputStream(cfg.In)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open input for reading: %s\n", err)
		os.Exit(2)
	}
	defer CloseStream(in)

	fmt.Println(cfg)
}

func openInputStream(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func CloseStream(src io.Closer) {
	if err := src.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close input stream: %s\n", err)
	}
}

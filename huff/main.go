package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	infile := flag.String("i", "", "The file to encode")
	outfile := flag.String("o", "", "The encoded file")
	flag.Parse()

	file, err := os.Open(*infile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s for reading: %s\n", *infile, err)
		os.Exit(1)
	}
	defer closeFile(file)

	fmt.Printf("Compressing %s and writing output to %s\n", *infile, *outfile)
}

func closeFile(file io.Closer) {
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close file: %s\n", err)
	}
}

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

	counts, err := makeFrequencyTable(file)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"The following error occurred while processing %s: %s",
			*infile,
			err,
		)
		os.Exit(2)
	}

	fmt.Printf("Compressing %s and writing output to %s\n", *infile, *outfile)
	fmt.Printf("%+v\n", counts)
}

func closeFile(file io.Closer) {
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close file: %s\n", err)
	}
}

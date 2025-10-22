package main

import (
	"flag"
	"fmt"
)

func main() {
	infile := flag.String("i", "", "The file to encode")
	outfile := flag.String("o", "", "The encoded file")
	flag.Parse()

	fmt.Printf("Compressing %s and writing output to %s\n", *infile, *outfile)
}

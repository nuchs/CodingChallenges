package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
)

const chunkSize = 64 * 1024 * 1024

// Huffman encode the input stream and write to the output stream
func Encode(in io.Reader, out io.Writer) error {
	for c, err := range chunks(in) {
		if err != nil && err != io.EOF {
			return fmt.Errorf("encoder reading chunk: %w", err)
		}
		if len(c) == 0 {
			break
		}

		enc := encodeChunk(c)
		n, err := out.Write(enc)
		if err != nil {
			return fmt.Errorf("writing encoded chunk: %w", err)
		}
		if n != len(enc) {
			return io.ErrShortWrite
		}
	}

	return nil
}

// Splits the input into chunks and prints the huffman encoding table for each
func PrintHeaders(in io.Reader, out io.Writer) error {
	buf := bufio.NewWriter(out)
	var i int
	for c, err := range chunks(in) {
		if err != nil && err != io.EOF {
			return fmt.Errorf("printer reading chunk: %w", err)
		}
		if len(c) == 0 {
			break
		}

		i++
		et := generateEncodingTable(c)
		_, err := fmt.Fprintf(buf, "Encoding Table for chunk %d\n%s\n", i, et)
		if err != nil {
			return fmt.Errorf("writing header: %w", err)
		}
	}

	return nil
}

// Splits data from 'in' into chunkSize byte slices. This allows the main
// processing logic to be separated from any io handling
func chunks(in io.Reader) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for {
			buf := make([]byte, chunkSize)
			_, err := in.Read(buf)
			if !yield(buf, err) || err != nil {
				break
			}
		}
	}
}

func encodeChunk(in []byte) []byte {
	et := generateEncodingTable(in)
	out := make([]byte, 0, 1)
	writeHeader(out, et)
	writeData(in, out, et)

	return out
}

func writeHeader(out []byte, et encodeTable) {
}

func writeData(in, out []byte, et encodeTable) {
}

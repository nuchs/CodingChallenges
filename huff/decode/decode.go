// Package decode is responsible for reading huffman encoded files and producing
// the original content.
package decode

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"iter"
)

var (
	ErrInsufficientData = errors.New("failed to read required number of bytes")
	ErrHeaderSize       = errors.New("invalid header size")
)

type chunk struct {
	header     []byte
	headerSize uint16
	body       []byte
	bodySize   uint16
}

// PrintDecodeTable iterates over the chunks in the encoded input and prints out
// their encoding tables.
func PrintDecodeTable(in io.Reader, out io.Writer) error {
	var i int
	for c, err := range chunks(in) {
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("printing decode table: %w", err)
		}
		i++
		dt := newDecodeTable(c.header)
		_, err := fmt.Fprintf(out, "Decoding Table for chunk %d\n%s\n", i, dt)
		if err != nil {
			return fmt.Errorf("writing decode table: %w", err)
		}
	}
	return nil
}

// Decode the huffman encoded data from the input and writes it to the output.
func Decode(in io.Reader, out io.Writer) error {
	return nil
}

func chunks(in io.Reader) iter.Seq2[chunk, error] {
	return func(yield func(chunk, error) bool) {
		for {
			var c chunk

			chunkSize, err := readSize(in)
			if err != nil {
				yield(c, fmt.Errorf("getting chunk size: %w", err))
				break
			}

			headerSize, err := readSize(in)
			if err != nil {
				yield(c, fmt.Errorf("getting header size: %w", err))
				break
			}
			if headerSize%4 != 0 {
				yield(c, ErrHeaderSize)
				break
			}

			c.headerSize = headerSize - 2 // don't include the header size itself
			c.header = make([]byte, c.headerSize)
			if err := read(in, c.header, c.headerSize); err != nil {
				yield(c, fmt.Errorf("reading header: %w", err))
				break
			}

			c.bodySize = chunkSize - (headerSize + 2) // don't include the header or chunk size
			c.body = make([]byte, c.bodySize)
			if err := read(in, c.body, c.bodySize); err != nil {
				yield(c, fmt.Errorf("reading body: %w", err))
				break
			}

			if !yield(c, nil) {
				break
			}
		}
	}
}

func read(in io.Reader, buf []byte, size uint16) error {
	n, err := in.Read(buf)
	if err != nil {
		return err
	}
	if n != int(size) {
		return ErrInsufficientData
	}

	return nil
}

func readSize(in io.Reader) (uint16, error) {
	var sizeBuf [2]byte
	n, err := in.Read(sizeBuf[:])
	if err != nil {
		return 0, fmt.Errorf("reading size: %w", err)
	}
	if n != 2 {
		return 0, ErrInsufficientData
	}
	size := binary.BigEndian.Uint16(sizeBuf[:])
	return size, nil
}

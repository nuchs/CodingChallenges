package encode

import (
	"fmt"
	"io"
	"iter"
)

const chunkSize = 64 * 1024

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
func PrintEncodeTable(in io.Reader, out io.Writer) error {
	var i int
	for c, err := range chunks(in) {
		if err != nil && err != io.EOF {
			return fmt.Errorf("printer reading chunk: %w", err)
		}
		if len(c) == 0 {
			break
		}

		i++
		et := newEncodingTable(c)
		_, err := fmt.Fprintf(out, "Encoding Table for chunk %d\n%s\n", i, et)
		if err != nil {
			return fmt.Errorf("writing header: %w", err)
		}
	}
	fmt.Fprintln(out)

	return nil
}

// Splits data from 'in' into chunkSize byte slices. This allows the main
// processing logic to be separated from any io handling
func chunks(in io.Reader) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for {
			buf := make([]byte, chunkSize)
			n, err := in.Read(buf)
			if n < chunkSize {
				buf = buf[:n]
			}
			if !yield(buf, err) || err != nil {
				break
			}
		}
	}
}

func encodeChunk(in []byte) []byte {
	et := newEncodingTable(in)
	out := make([]byte, 2, len(in)) // first two bytes are to record chunk size
	out = writeHeader(out, et)
	out = writeData(in, out, et)

	size := uint16(len(out))
	out[0] = byte(size >> 8)
	out[1] = byte(0x00FF & size)

	return out
}

func writeHeader(out []byte, et encodeTable) []byte {
	headerSize := uint16(2 + 4*len(et))
	out = append(out, byte(headerSize>>8))
	out = append(out, byte(0x00FF&headerSize))

	for _, e := range et {
		out = append(out, e.value)
		out = append(out, e.len)
		out = append(out, byte(e.prefix>>8))
		out = append(out, byte(0x00FF&e.prefix))
	}

	return out
}

func writeData(in, out []byte, et encodeTable) []byte {
	bb := newShortBuffer(out)
	for _, c := range in {
		e := et[c]
		bb.write(e.prefix, e.len)
	}
	return bb.close()
}

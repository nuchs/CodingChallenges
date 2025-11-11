package main

const BufLen = 16

type ShortBuffer struct {
	buf uint16
	len byte
	out []byte
}

func newShortBuffer(out []byte) *ShortBuffer {
	return &ShortBuffer{out: out}
}

func (bb *ShortBuffer) write(val uint16, vlen byte) {
	shift := vlen
	var rem byte
	if shift+bb.len > BufLen {
		shift = BufLen - bb.len
		rem = vlen - shift
	}
	bb.buf = (bb.buf << shift) + (val >> rem)
	bb.len += shift

	if bb.len != BufLen {
		return
	}
	bb.out = append(bb.out, byte(bb.buf>>8))
	bb.out = append(bb.out, byte(0x00FF&bb.buf))
	bb.len = rem
	bb.buf = val & bitmask(rem)
}

func (bb *ShortBuffer) close() []byte {
	shift := bb.flush()
	padding := shift % 8
	bb.out = append(bb.out, padding)

	return bb.out
}

func (bb *ShortBuffer) flush() byte {
	if bb.len == 0 {
		return 0
	}
	shift := BufLen - bb.len
	bb.buf <<= shift
	bb.out = append(bb.out, byte(bb.buf>>8))
	if bb.len > 8 {
		bb.out = append(bb.out, byte(0x00FF&bb.buf))
	}
	return shift
}

func bitmask(size byte) uint16 {
	return (1 << size) - 1
}

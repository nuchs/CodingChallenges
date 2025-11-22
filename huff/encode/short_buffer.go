package encode

const bufLen = 16

type shortBuffer struct {
	buf uint16
	len byte
	out []byte
}

func newShortBuffer(out []byte) *shortBuffer {
	return &shortBuffer{out: out}
}

func (bb *shortBuffer) write(val uint16, vlen byte) {
	shift := vlen
	var rem byte
	if shift+bb.len > bufLen {
		shift = bufLen - bb.len
		rem = vlen - shift
	}
	bb.buf = (bb.buf << shift) + (val >> rem)
	bb.len += shift

	if bb.len != bufLen {
		return
	}
	bb.out = append(bb.out, byte(bb.buf>>8))
	bb.out = append(bb.out, byte(0x00FF&bb.buf))
	bb.len = rem
	bb.buf = val & bitmask(rem)
}

func (bb *shortBuffer) close() []byte {
	shift := bb.flush()
	padding := shift % 8
	bb.out = append(bb.out, padding)

	return bb.out
}

func (bb *shortBuffer) flush() byte {
	if bb.len == 0 {
		return 0
	}
	shift := bufLen - bb.len
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

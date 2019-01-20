package parse

var (
	lookup = [...]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0}
)

func decode(a, b, c, d byte) (byte, byte) {
	r := int(lookup[a]) << 12
	r += int(lookup[b]) << 8
	r += int(lookup[c]) << 4
	r += int(lookup[d])

	b0 := byte(r >> 6)
	b0 |= 0xc0
	b1 := byte(r) & 0x3f
	b1 |= 0x80
	return b0, b1
}

func UnquoteInplace(b []byte) []byte {
	n := len(b)
	i := 0
	j := 0
	for i < n {
		if (i+5) < n && b[i] == 0x5C && b[i+1] == 0x75 {
			b[j], b[j+1] = decode(b[i+2], b[i+3], b[i+4], b[i+5])
			i += 6
			j += 2
		} else {
			b[j] = b[i]
			i++
			j++
		}
	}
	return b[:j]
}

func ParseSpaces(b []byte) []byte {
	for i, c := range b {
		if c > 0x20 {
			return b[i:]
		}
	}
	return nil
}

func ParseSymbol(b []byte, x byte) ([]byte, bool) {
	for i, c := range b {
		if c > 0x20 {
			if c != x {
				return b, false
			}
			return b[i+1:], true
		}
	}
	return b, false
}

func ParseNumbers(b []byte) ([]byte, []byte, bool) {
	for i, c := range b {
		if c > 0x20 {
			if c < 0x30 || c > 0x39 {
				return b, nil, false
			}
			for j, c := range b[i:] {
				if c < 0x30 || c > 0x39 {
					return b[i+j:], b[i : j+1], true
				}
			}
			return nil, b[i:], true
		}
	}
	return b, nil, false
}

func ParseQuoted(b []byte) ([]byte, []byte, bool) {
	for i, c := range b {
		if c > 0x20 {
			if c != 0x22 {
				return b, nil, false
			}
			i++
			for j, c := range b[i:] {
				if c == 0x22 {
					return b[i+j+1:], b[i : i+j], true
				}
			}
			return b, nil, false
		}
	}
	return b, nil, false
}

func AtoiNocheck(b []byte) int {
	x := 0
	for _, c := range b {
		x *= 10
		x += int(c) - 0x30
	}
	return x
}
package parse

func ParseSpaces(b []byte) []byte {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			return b[i:]
		}
	}
	return b[n:]
}

func ParseSymbol(b []byte, c byte) ([]byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] == c {
				return b[i+1:], true
			}
			return b, false
		}
	}
	return b, false
}

func ParseInt(b []byte) ([]byte, int, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] < 0x30 || b[i] > 0x39 {
				return b, 0, false
			}
			x := int(b[i]) - 0x30
			i++
			for i < n {
				if b[i] < 0x30 || b[i] > 0x39 {
					break
				}
				x *= 10
				x += int(b[i]) - 0x30
				i++
			}
			return b[i:], x, true
		}
	}
	return b, 0, false
}

func ParseQuoted(b []byte) ([]byte, []byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] != 0x22 {
				return b, nil, false
			}
			p := i + 1
			i++
			for i < n {
				j := i + 1
				switch b[i] {
				case 0x5C:
					if (i+5) < n && b[j] == 0x75 {
						i += 5
					} else {
						i += 1
					}
				case 0x22:
					return b[j:], b[p:i], true
				}
				i++
			}
			return b, nil, false
		}
	}
	return b, nil, false
}

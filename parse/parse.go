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

func ParseInt(b []byte) (int, []byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] < 0x30 || b[i] > 0x39 {
				return 0, b, false
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
			return x, b[i:], true
		}
	}
	return 0, b, false
}

func ParsePhrase(b []byte, c byte) ([]byte, []byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] != c {
				return nil, b, false
			}
			p := i
			i++
			for i < n {
				switch b[i] {
				case '\\':
					if i+5 < n && b[i+1] == 'u' {
						i += 5
					} else {
						i += 1
					}
				case c:
					return b[p : i+1], b[i+1:], true
				}
				i++
			}
			return nil, b, false
		}
	}
	return nil, b, false
}

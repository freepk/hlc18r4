package parse

func SkipSpaces(b []byte) []byte {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			return b[i:]
		}
	}
	return nil
}

func SkipSymbol(b []byte, x byte) ([]byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] != x {
				return b, false
			}
			return b[i+1:], true
		}
	}
	return b, false
}

func ParseNumber(b []byte) ([]byte, []byte, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] < 0x30 || b[i] > 0x39 {
				return b, nil, false
			}
			for j := i + 1; j < n; j++ {
				if b[j] < 0x30 || b[j] > 0x39 {
					return b[j:], b[i:j], true
				}
			}
			return nil, b[i:], true
		}
	}
	return b, nil, false
}

func ParseInt(b []byte) ([]byte, int, bool) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] > 0x20 {
			if b[i] < 0x30 || b[i] > 0x39 {
				return b, 0, false
			}
			x := int(b[i]) - 0x30
			for j := i + 1; j < n; j++ {
				if b[j] < 0x30 || b[j] > 0x39 {
					return b[j:], x, true
				}
				x *= 10
				x += int(b[j]) - 0x30
			}
			return nil, x, true
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
			for j := i + 1; j < n; j++ {
				switch b[j] {
				case 0x5C:
					j++
					if j < n && b[j] == 0x75 {
						j += 4
					}
				case 0x22:
					return b[j+1:], b[i+1 : j], true
				}
			}
			return b, nil, false
		}
	}
	return b, nil, false
}

func ScanSymbol(b []byte, c byte) ([]byte, []byte) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] == c {
			return b[i+1:], b[:i]
		}
	}
	return nil, b
}

package main

func parseSymbol(b []byte, c byte) ([]byte, bool) {
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

func parseInt(b []byte) (int, []byte, bool) {
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

package main

func utf8Unhex(b byte) (int, bool) {
	c := int(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

func utf8UnquoteChar(a, b, c, d byte) (byte, byte, bool) {
	va, ok := utf8Unhex(a)
	if !ok {
		return 0, 0, false
	}
	va <<= 12
	vb, ok := utf8Unhex(b)
	if !ok {
		return 0, 0, false
	}
	vb <<= 8
	vc, ok := utf8Unhex(c)
	if !ok {
		return 0, 0, false
	}
	vc <<= 4
	vd, ok := utf8Unhex(d)
	if !ok {
		return 0, 0, false
	}
	v := va + vb + vc + vd
	x := []byte(string(v))
	return x[0], x[1], true
}

func utf8Unquote(dst, src []byte) int {
	n := len(src)
	i := 0
	j := 0
	for i < n {
		if i+6 <= n && src[i] == '\\' && src[i+1] == 'u' {
			a, b, ok := utf8UnquoteChar(src[i+2], src[i+3], src[i+4], src[i+5])
			if ok {
				dst[j] = a
				j++
				dst[j] = b
				j++
				i += 6
				continue
			}
		}
		dst[j] = src[i]
		i++
		j++
	}
	return j
}

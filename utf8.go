package main

func utf8Lookup(c byte) int {
	if c >= 48 && c <= 57 {
		return int(c) - 48
	} else if c >= 97 && c <= 102 {
		return int(c) - 87
	}
	return -1
}
func lookup(c byte) int {
	if c >= 48 && c <= 57 {
		return int(c) - 48
	} else if c >= 97 && c <= 102 {
		return int(c) - 87
	}
	return -1
}

func decode(a, b, c, d byte) (byte, byte, bool) {
	va := lookup(a)
	if va < 0 {
		return 0, 0, false
	}
	va <<= 12
	vb := lookup(b)
	if vb < 0 {
		return 0, 0, false
	}
	vb <<= 8
	vc := lookup(c)
	if vc < 0 {
		return 0, 0, false
	}
	vc <<= 4
	vd := lookup(d)
	if vd < 0 {
		return 0, 0, false
	}
	v := int(va + vb + vc + vd)
	x := []byte(string(v))
	return x[0], x[1], true
}

func unquote(dst, src []byte) int {
	n := len(src)
	i := 0
	j := 0
	for i < n {
		if (i+5) < n && src[i] == '\\' && src[i+1] == 'u' {
			a, b, ok := decode(src[i+2], src[i+3], src[i+4], src[i+5])
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

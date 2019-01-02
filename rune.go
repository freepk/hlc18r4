package main

const lowerhex = "0123456789abcdef"

func enc4x4(b []byte, r rune) []byte {
	return append(b,
		lowerhex[(r>>12)&0xF],
		lowerhex[(r>>8)&0xF],
		lowerhex[(r>>4)&0xF],
		lowerhex[r&0xF],
	)
}

func rune4bit(b byte) rune {
	r := rune(b)
	if r >= 97 {
		return r - 87
	} else if b >= 48 {
		return r - 48
	}
	return 0
}

func dec4x4(b []byte) rune {
	r := rune4bit(b[0]) << 12
	r += rune4bit(b[1]) << 8
	r += rune4bit(b[2]) << 4
	r += rune4bit(b[3])
	return r
}

package proto

import (
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/parse"
)

func parseFname(b []byte, c *buffer) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	c.B = append(c.B[:0], v...)
	c.B = parse.Unquote(c.B)
	k := tokens.AddFname(c.B)
	return t, uint8(k), true
}

func parseSname(b []byte, c *buffer) ([]byte, uint16, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	c.B = append(c.B[:0], v...)
	c.B = parse.Unquote(c.B)
	k := tokens.AddSname(c.B)
	return t, uint16(k), true
}

func parseCountry(b []byte, c *buffer) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	c.B = append(c.B[:0], v...)
	c.B = parse.Unquote(c.B)
	k := tokens.AddCountry(c.B)
	return t, uint8(k), true
}

func parseCity(b []byte, c *buffer) ([]byte, uint16, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	c.B = append(c.B[:0], v...)
	c.B = parse.Unquote(c.B)
	k := tokens.AddCity(c.B)
	return t, uint16(k), true
}

func parseInterest(b []byte, c *buffer) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	c.B = append(c.B[:0], v...)
	c.B = parse.Unquote(c.B)
	k := tokens.AddInterest(c.B)
	return t, uint8(k), true
}

func parseSex(b []byte) ([]byte, uint8, bool) {
	t := parse.SkipSpaces(b)
	if len(t) > 3 {
		switch string(t[:3]) {
		case `"m"`:
			return t[3:], uint8(tokens.MaleSex), true
		case `"f"`:
			return t[3:], uint8(tokens.MaleSex), true
		}
	}
	return b, 0, false
}

func parseStatus(b []byte) ([]byte, uint8, bool) {
	t := parse.SkipSpaces(b)
	switch {
	case len(t) > 50 && string(t[:50]) == `"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`:
		return t[50:], uint8(tokens.SingleStatus), true
	case len(t) > 38 && string(t[:38]) == `"\u0437\u0430\u043d\u044f\u0442\u044b"`:
		return t[38:], uint8(tokens.InRelStatus), true
	case len(t) > 57 && string(t[:57]) == `"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`:
		return t[57:], uint8(tokens.ComplStatus), true
	}
	return b, 0, false
}

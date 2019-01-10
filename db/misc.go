package db

import (
	"bytes"
)

func splitEmail(b []byte) ([]byte, []byte, bool) {
	p := bytes.IndexByte(b, '@')
	if p == -1 {
		return nil, b, false
	}
	return b[:p], b[p:], true
}

package main

import (
	"bytes"
)

func skipSpace(b []byte) []byte {
	for i := range b {
		if b[i] > 0x20 {
			return b[i:]
		}
	}
	return nil
}

func parseAccount(b []byte) ([]byte, bool) {
	// SkipSpace and check '{'
	b = skipSpace(b)
	if len(b) == 0 {
		return b, false
	}
	if b[0] != '{' {
		return b, false
	}
	b = b[1:]

	// Parse key
	var k []byte
	// SkipSpace and check '"'
	b = skipSpace(b)
	if len(b) == 0 {
		return b, false
	}
	if b[0] != '"' {
		return b, false
	}
	b = b[1:]
	// Lookup double quote
	for i, c := range b {
		if c == '"' {
			k = b[:i]
			b = b[i+1:]
			break
		}
		if c == '\\' {
			return b, false
		}
	}
	// SkipSpace and check ':'
	b = skipSpace(b)
	if len(b) == 0 {
		return b, false
	}
	if b[0] != ':' {
		return b, false
	}
	b = b[1:]

	if len(k) < 2 {
		return b, false
	}
	// Validate key
	//	i	id interests
	//	e	email
	//	f	finish fname
	//	p	phone premium
	//	b	birth
	//	c	city country
	//	j	joined
	//	s	sex sname status
	//	l	likes
	switch k[0] {
	case 'i':
		switch k[1] {
		case 'd':
			if !bytes.Equal(k, []byte(`id`)) {
				return b, false
			}
		case 'n':
			if !bytes.Equal(k, []byte(`interests`)) {
				return b, false
			}
		}
	case 'e':
		if !bytes.Equal(k, []byte(`email`)) {
			return b, false
		}
	case 'f':
		switch k[1] {
		case 'i':
			if !bytes.Equal(k, []byte(`finish`)) {
				return b, false
			}
		case 'n':
			if !bytes.Equal(k, []byte(`fname`)) {
				return b, false
			}
		}
	case 'p':
		switch k[1] {
		case 'h':
			if !bytes.Equal(k, []byte(`phone`)) {
				return b, false
			}
		case 'r':
			if !bytes.Equal(k, []byte(`premium`)) {
				return b, false
			}
		}
	case 'b':
		if !bytes.Equal(k, []byte(`birth`)) {
			return b, false
		}
	case 'c':
		switch k[1] {
		case 'i':
			if !bytes.Equal(k, []byte(`city`)) {
				return b, false
			}
		case 'o':
			if !bytes.Equal(k, []byte(`country`)) {
				return b, false
			}
		}
	case 'j':
		if !bytes.Equal(k, []byte(`joined`)) {
			return b, false
		}
	case 's':
		switch k[1] {
		case 'e':
			if !bytes.Equal(k, []byte(`sex`)) {
				return b, false
			}
		case 'n':
			if !bytes.Equal(k, []byte(`sname`)) {
				return b, false
			}
		case 't':
			if !bytes.Equal(k, []byte(`status`)) {
				return b, false
			}
		}
	case 'l':
		if !bytes.Equal(k, []byte(`likes`)) {
			return b, false
		}
	}

	return b, true
}

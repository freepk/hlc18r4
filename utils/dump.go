package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	xbytes "github.com/freepk/bytes"
	"io"
	"log"
	"os"
	"unicode/utf8"
)

func utf8Unhex(b byte) (rune, bool) {
	c := rune(b)
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

func utf8UnquoteChar(b []byte) (rune, bool) {
	if len(b) < 6 {
		return 0, false
	}
	if b[0] != '\\' {
		return 0, false
	}
	if b[1] != 'u' {
		return 0, false
	}
	r := rune(0)
	for i := 2; i < 6; i++ {
		x, ok := utf8Unhex(b[i])
		if !ok {
			return 0, false
		}
		r = (r << 4) | x
	}
	return r, true
}

func utf8Unquote(d, s []byte) int {
	n := len(s)
	if n > len(d) {
		return 0
	}
	i := 0
	j := 0
	for i < n {
		c, ok := utf8UnquoteChar(s[i:])
		if ok {
			z := utf8.EncodeRune(d[j:], c)
			j += z
			i += 6
			continue
		}
		d[j] = s[i]
		i++
		j++
	}
	return j
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("no input file")
	}
	arch, err := zip.OpenReader(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer arch.Close()
	n := len(arch.File)
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		src, err := arch.File[i].Open()
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
		io.Copy(buf, src)
		payload := buf.Bytes()
		//payload = payload[14 : len(payload)-2]
		payload = payload[14 : 8192*4]

		for {
			a, b, ok := xbytes.IndexScoped(payload, '\\', '"', '{', '}')
			if !ok {
				break
			}
			json := payload[a : b+1]
			fmt.Println(string(json))
			payload = payload[b+1:]
		}

		src.Close()
	}
}

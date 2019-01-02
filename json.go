package main

import (
	"io"
)

func jsonFindWrapped(b []byte, x, y byte) (int, int, bool) {
	q := 0
	i := 0
	for j, c := range b {
		switch c {
		case x:
			q++
			if q == 1 {
				i = j
			}
		case y:
			q--
			if q == 0 {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func jsonFindObj(b []byte) (int, int, bool) {
	return jsonFindWrapped(b, 123, 125)
}

func jsonFindArr(b []byte) (int, int, bool) {
	return jsonFindWrapped(b, 91, 93)
}

func jsonFindQuoted(b []byte) (int, int, bool) {
	q := 0
	i := 0
	for j, c := range b {
		if c == 34 {
			q++
			if q == 2 {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func jsonReadObj(r io.Reader, b []byte, x int, f func([]byte) error) error {
	p := 0
	for {
		if n, err := r.Read(b[p:]); n > 0 {
			n += p
			t := b[x:n]
			for {
				i, j, ok := jsonFindObj(t)
				if !ok {
					break
				}
				j++
				err = f(t[i:j])
				if err != nil {
					return err
				}
				t = t[j:]
			}
			x = 0
			p = copy(b, t)
		} else if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

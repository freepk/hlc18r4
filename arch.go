package main

import (
	"archive/zip"
	"io"
	"log"
	"sync"
)

func readFile(wait *sync.WaitGroup, file *zip.File) {
	if wait != nil {
		defer wait.Done()
	}
	r, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	p := 0
	x := 14
	b := make([]byte, 8192)
	for {
		if n, err := r.Read(b[p:]); n > 0 {
			n += p
			c := 0
			i := 0
			j := x
			for j < n {
				switch b[j] {
				case '{':
					c++
					if c == 1 {
						i = j
					}
				case '}':
					c--
					if c == 0 {
						j++
					}
				}
				j++
			}
			x = 0
			p = copy(b, b[i:n])
		} else if err == io.EOF {
			return
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func readArchive(path string) {
	z, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	n := len(z.File)
	w := new(sync.WaitGroup)
	w.Add(n)
	for i := 0; i < n; i++ {
		f := z.File[i]
		go readFile(w, f)
	}
	w.Wait()
}

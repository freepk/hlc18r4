package main

import (
	"archive/zip"
	"io"
	"log"
	"sync"
)

func readData(r io.Reader) {
	acc := &Account{}
	b := make([]byte, 8192)
	p := 0
	x := 14
	for {
		if n, err := r.Read(b[p:]); n > 0 {
			n += p
			t, ok := b[x:n], true
			for {
				t, ok = parseSymbol(t, ',')
				t, ok = acc.Parse(t)
				if !ok {
					break
				}
			}
			p = copy(b, t)
			x = 0
		} else if err == io.EOF {
			return
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func readBackup(path string) {
	a, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	n := len(a.File)
	w := new(sync.WaitGroup)
	w.Add(n)
	for i := 0; i < n; i++ {
		f := a.File[i]
		log.Println("File", f.Name)
		r, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer r.Close()
			defer w.Done()
			readData(r)
		}()
	}
	w.Wait()
}

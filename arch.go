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
	b := make([]byte, 4096)
	for {
		n, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < n; i++ {
			c := b[i]
			if c > 0 {
			}
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

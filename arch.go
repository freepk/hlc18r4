package main

import (
	"archive/zip"
	"log"
	"sync"
)

func readFile(wait *sync.WaitGroup, file *zip.File) {
	if wait != nil {
		defer wait.Done()
	}
	println(file.Name)
	r, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	b := make([]byte, 8192)
	jsonReadObj(r, b, 14, func(b []byte) error {
		//println(string(b))
		return nil
	})
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

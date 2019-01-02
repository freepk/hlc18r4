package main

import (
	"archive/zip"
	//"fmt"
	"sync"
	"testing"
)

func TestReadArchive(t *testing.T) {
	z, err := zip.OpenReader("data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	n := len(z.File)
	w := new(sync.WaitGroup)
	w.Add(n)
	for i := 0; i < n; i++ {
		r, err := z.File[i].Open()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			defer w.Done()
			defer r.Close()
			b := make([]byte, 8192)
			jsonReadObj(r, b, 14, func(b []byte) error {
				//fmt.Println(string(b))
				return nil
			})
		}()
	}
	w.Wait()
}

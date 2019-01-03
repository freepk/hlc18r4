package main

import (
	"archive/zip"
	"sync"
	"testing"
)

func TestReadArchive(t *testing.T) {
	arch, err := zip.OpenReader("data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer arch.Close()
	n := 1 //len(arch.File)
	group := new(sync.WaitGroup)
	group.Add(n)
	for i := 0; i < n; i++ {
		src, err := arch.File[i].Open()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			defer group.Done()
			defer src.Close()
		}()
	}
	group.Wait()
}

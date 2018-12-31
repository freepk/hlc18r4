package main

import (
	"archive/zip"
	"io"
	"log"
	"sync"
	"testing"

	"github.com/freepk/hashtab"
	"github.com/freepk/workerpool"
)

type Account struct {
	f0  int
	f1  int
	f2  int
	f3  int
	f4  int
	f5  int
	f6  int
	f7  int
	f10 int
	f11 int
	f12 int
	f13 int
	f14 int
	f15 int
	f16 int
	f17 int
}

var (
	pool            *workerpool.Pool
	accounts        []Account
	accountsHashTab *hashtab.HashTab
)

func init() {
	pool = workerpool.NewPool(4)
	go pool.Start()
	power := uint32(21)
	accounts = make([]Account, (1 << power))
	accountsHashTab = hashtab.NewHashTab(power)
}

func readFile(file *zip.File) {
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
		pool.Run(func() {
			readFile(f)
			w.Done()
		})
	}
	w.Wait()
}

func TestReadArchive(t *testing.T) {
	readArchive("data/data.zip")
}

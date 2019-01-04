package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	xbytes "github.com/freepk/bytes"
	"io"
	"sync"
	"testing"
)

func TestReadArchive(t *testing.T) {
	return
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
			defer src.Close()
			defer group.Done()
		}()
	}
	group.Wait()
}

func TestReadArchiveX(t *testing.T) {
	arch, err := zip.OpenReader("data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer arch.Close()
	n := 1 //len(arch.File)
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		src, err := arch.File[i].Open()
		if err != nil {
			t.Fatal(err)
		}
		buf.Reset()
		io.Copy(buf, src)
		payload := buf.Bytes()
		payload = payload[14 : len(payload)-2]

		for {
			a, b, ok := xbytes.IndexScoped(payload, '\\', '"', '{', '}')
			if !ok {
				break
			}
			fmt.Println(string(payload[a:b+1]))
			payload = payload[b+1:]
		}

		src.Close()
	}
}

package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

func readData(r io.Reader) {
	buf := make([]byte, 8192)
	//p := 0
	//x := 14
	for {
		n, err := r.Read(buf[p:])
		if n > 0 {
			n += p
		} else if err == io.EOF {
			return
		} else if err != nil {
			log.Fatal(err)
		}
	}
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
	for i := 0; i < n; i++ {
		log.Println("File", arch.File[i].Name)
		src, err := arch.File[i].Open()
		if err != nil {
			log.Fatal(err)
		}
		readData(src)
		src.Close()
	}
}

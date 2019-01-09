package db

import (
	//"archive/zip"
	"io"
	"log"
	"sync"

<<<<<<< HEAD
	"github.com/klauspost/compress/zip"
=======
	"gitlab.com/freepk/hlc18r4/json"
>>>>>>> 1080997c262264826bb2bda81788695301410c77
	"gitlab.com/freepk/hlc18r4/parse"
)

type account struct {
	country byte
}

type DB struct {
	a []account
}

func NewDB() *DB {
	a := make([]account, 1400000)
	return &DB{a: a}
}

func (db *DB) readData(r io.Reader) {
	a := &json.Account{}
	b := make([]byte, 8192)
	p := 0
	x := 14
	for {
		if n, err := r.Read(b[p:]); n > 0 {
			n += p
			t, ok := b[x:n], true
			for {
				t, ok = parse.ParseSymbol(t, ',')
				a.Reset()
				t, ok = a.UnmarshalJSON(t)
				if !ok {
					break
				}
				x := db.a[a.ID]
				x.country = 10
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

func (db *DB) Restore(path string) {
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
		r, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer r.Close()
			defer w.Done()
			db.readData(r)
		}()
	}
	w.Wait()
}

package db

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/freepk/hashtab"
	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/hash"
	"gitlab.com/freepk/hlc18r4/lookup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	DefaultError = errors.New("Error")
)

var (
	domainLookup    = lookup.NewLookup(16)
	emailHashTab, _ = hashtab.NewHashTab(21)
)

type DB struct {
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) insertAccount(src *proto.Account) error {
	login, domain, ok := splitEmail(src.Email)
	if !ok {
		return DefaultError
	}
	emailHash := hash.Hash64(src.Email, 1234)
	id, ok := emailHashTab.GetOrSet(uint64(emailHash), uint64(src.ID))
	if ok {
		log.Println(login, domain, id, src.ID)
	}
	return nil
}

func (db *DB) readData(r io.Reader) {
	a := &proto.Account{}
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
				db.insertAccount(a)
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

package db

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/freepk/hashtab"
	"github.com/klauspost/compress/zip"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/lookup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	DefaultError = errors.New("Error")
)

var (
	emailHashTab, _   = hashtab.NewHashTab(21)
	domainLookup, _   = lookup.NewLookup(4)
	fnameLookup, _    = lookup.NewLookup(8)
	snameLookup, _    = lookup.NewLookup(12)
	sexLookup, _      = lookup.NewLookup(4)
	countryLookup, _  = lookup.NewLookup(8)
	cityLookup, _     = lookup.NewLookup(10)
	statusLookup, _   = lookup.NewLookup(4)
	interestLookup, _ = lookup.NewLookup(8)
)

func Print() {
	fmt.Println("\n\nid domain")
	domainLookup.Print()
	fmt.Println("\n\nid fname")
	fnameLookup.Print()
	fmt.Println("\n\nid sname")
	snameLookup.Print()
	fmt.Println("\n\nid sex")
	sexLookup.Print()
	fmt.Println("\n\nid country")
	countryLookup.Print()
	fmt.Println("\n\nid city")
	cityLookup.Print()
	fmt.Println("\n\nid status")
	statusLookup.Print()
	fmt.Println("\n\nid interest")
	interestLookup.Print()
}

type DB struct {
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) insertAccount(src *proto.Account) error {
	emailHash := murmur3.Sum64(src.Email)
	id, ok := emailHashTab.GetOrSet(uint64(emailHash), uint64(src.ID))
	if ok {
		log.Println("Email duplicate", src.Email, src.ID, id)
	}
	login, domain, ok := splitEmail(src.Email)
	if !ok {
		return DefaultError
	}
	_ = login
	domainLookup.GetOrGen(domain)
	fnameLookup.GetOrGen(src.Fname)
	snameLookup.GetOrGen(src.Sname)
	sexLookup.GetOrGen(src.Sex)
	countryLookup.GetOrGen(src.Country)
	cityLookup.GetOrGen(src.City)
	statusLookup.GetOrGen(src.Status)
	n := len(src.Interests)
	for i := 0; i < n; i++ {
		interestLookup.GetOrGen(src.Interests[i])
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

package db

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/json"
	"gitlab.com/freepk/hlc18r4/lookup"
	"gitlab.com/freepk/hlc18r4/parse"
)

type account struct {
}

type DB struct {
	// email
	// phone
	fname    *lookup.Lookup
	sname    *lookup.Lookup
	sex      *lookup.Lookup
	country  *lookup.Lookup
	city     *lookup.Lookup
	status   *lookup.Lookup
	interest *lookup.Lookup
	accounts []account
}

func NewDB() *DB {
	return &DB{
		fname:    lookup.NewLookup(128),
		sname:    lookup.NewLookup(1024),
		sex:      lookup.NewLookup(2),
		country:  lookup.NewLookup(128),
		city:     lookup.NewLookup(1024),
		status:   lookup.NewLookup(4),
		interest: lookup.NewLookup(128),
		accounts: make([]account, 1400000)}
}

func (db *DB) insertAccount(a *json.Account) {
	if a.ID > 0 && len(a.Email) > 0 {
		if len(a.Fname) > 0 {
			db.fname.GetKeyOrSet(a.Fname)
		}
		if len(a.Sname) > 0 {
			db.sname.GetKeyOrSet(a.Sname)
		}
		if len(a.Sex) > 0 {
			db.sex.GetKeyOrSet(a.Sex)
		}
		if len(a.Country) > 0 {
			db.country.GetKeyOrSet(a.Country)
		}
		if len(a.City) > 0 {
			db.city.GetKeyOrSet(a.City)
		}
		if len(a.Status) > 0 {
			db.status.GetKeyOrSet(a.Status)
		}
		n := len(a.Interests)
		for i := 0; i < n; i++ {
			db.interest.GetKeyOrSet(a.Interests[i])
		}
	}
}

func (db *DB) dump() {
	fmt.Printf("\n\nFname: %#v", db.fname)
	fmt.Printf("\n\nSname: %#v", db.sname)
	fmt.Printf("\n\nSex: %#v", db.sex)
	fmt.Printf("\n\nCountry: %#v", db.country)
	fmt.Printf("\n\nCity: %#v", db.city)
	fmt.Printf("\n\nStatus: %#v", db.status)
	fmt.Printf("\n\nInterest: %#v", db.interest)
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
	db.dump()
}

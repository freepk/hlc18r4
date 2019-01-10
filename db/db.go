package db

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/lookup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

type account struct {
	fname   uint8
	sname   uint16
	sex     bool
	country uint8
	city    uint8
	status  uint8
}

type DB struct {
	email    *lookup.Lookup
	fname    *lookup.Lookup
	sname    *lookup.Lookup
	phone    *lookup.Lookup
	sex      *lookup.Lookup
	country  *lookup.Lookup
	city     *lookup.Lookup
	status   *lookup.Lookup
	interest *lookup.Lookup
	accounts []account
}

func NewDB() *DB {
	return &DB{
		email:    lookup.NewLookup(1400000),
		fname:    lookup.NewLookup(128),
		sname:    lookup.NewLookup(2048),
		phone:    lookup.NewLookup(1048576),
		sex:      lookup.NewLookup(4),
		country:  lookup.NewLookup(128),
		city:     lookup.NewLookup(1024),
		status:   lookup.NewLookup(8),
		interest: lookup.NewLookup(128),
		accounts: make([]account, 1500000)}
}

func (db *DB) insertAccount(a *proto.Account) {
	db.email.GetKeyOrSet(a.Email)
	db.fname.GetKeyOrSet(a.Fname)
	db.sname.GetKeyOrSet(a.Sname)
	db.phone.GetKeyOrSet(a.Phone)
	db.sex.GetKeyOrSet(a.Sex)
	db.country.GetKeyOrSet(a.Country)
	db.city.GetKeyOrSet(a.City)
	db.status.GetKeyOrSet(a.Status)
	n := len(a.Interests)
	for i := 0; i < n; i++ {
		db.interest.GetKeyOrSet(a.Interests[i])
	}
}

func (db *DB) printStats() {
	fmt.Println("email", db.email.LastKey())
	fmt.Println("fname", db.fname.LastKey())
	fmt.Println("sname", db.sname.LastKey())
	fmt.Println("phone", db.phone.LastKey())
	fmt.Println("sex", db.sex.LastKey())
	fmt.Println("country", db.country.LastKey())
	fmt.Println("city", db.city.LastKey())
	fmt.Println("status", db.status.LastKey())
	fmt.Println("interest", db.interest.LastKey())
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

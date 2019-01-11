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
	emailHashTab   = hashtab.NewHashTab(21)
	domainLookup   = lookup.NewLookup(4)
	fnameLookup    = lookup.NewLookup(8)
	snameLookup    = lookup.NewLookup(12)
	sexLookup      = lookup.NewLookup(4)
	countryLookup  = lookup.NewLookup(8)
	cityLookup     = lookup.NewLookup(10)
	statusLookup   = lookup.NewLookup(4)
	interestLookup = lookup.NewLookup(8)
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

type account struct {
	domain    uint8
	fname     uint8
	sname     uint16
	sex       uint8
	country   uint8
	city      uint16
	status    uint8
	loginSize uint8
	login     [24]byte
}

type DB struct {
	a []account
}

func NewDB() *DB {
	return &DB{a: make([]account, 1400000)}
}

func (db *DB) insertAccount(src *proto.Account) error {
	emailHash := murmur3.Sum64(src.Email)
	id, ok := emailHashTab.GetOrSet(uint64(emailHash), uint64(src.ID))
	if ok {
		log.Println("Email duplicate", src.Email, src.ID, id)
		return DefaultError
	}
	login, domain, ok := splitEmail(src.Email)
	if !ok {
		return DefaultError
	}
	n := len(login)
	if n > 24 {
		return DefaultError
	}
	dst := &db.a[src.ID]
	dst.loginSize = uint8(n)
	for i := 0; i < n; i++ {
		dst.login[i] = login[i]
	}
	k := 0
	k, _ = domainLookup.GetOrGen(domain)
	dst.domain = uint8(k)
	k, _ = fnameLookup.GetOrGen(src.Fname)
	dst.fname = uint8(k)
	k, _ = snameLookup.GetOrGen(src.Sname)
	dst.sname = uint16(k)
	k, _ = sexLookup.GetOrGen(src.Sex)
	dst.sex = uint8(k)
	k, _ = countryLookup.GetOrGen(src.Country)
	dst.country = uint8(k)
	k, _ = cityLookup.GetOrGen(src.City)
	dst.city = uint16(k)
	k, _ = statusLookup.GetOrGen(src.Status)
	dst.status = uint8(k)

	//n := len(src.Interests)
	//for i := 0; i < n; i++ {
	//	interestLookup.GetOrGen(src.Interests[i])
	//}
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

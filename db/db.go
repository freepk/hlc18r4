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
	likeCounters   [20]int
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
	fmt.Println("# Page size", likePageSize, "like struct size", likeStructSize)
	fmt.Printf("%8s %8s %10s %10s\n", "pages", "accounts", "part", "total")
	totalSize := 0
	for i := 0; i < len(likeCounters); i++ {
		pages := i + 1
		accounts := likeCounters[i]
		partSize := accounts * pages * likePageSize
		totalSize += partSize
		fmt.Printf("%8d %8d %10d %10d\n", pages, accounts, partSize, totalSize)
	}
	fmt.Println("\n# Domain", domainLookup.LastKey())
	fmt.Println("id qty domain")
	domainLookup.Print()
	fmt.Println("\n# Fname", fnameLookup.LastKey())
	fmt.Println("id qty fname")
	fnameLookup.Print()
	//fmt.Println("\n# Sname", snameLookup.LastKey())
	//fmt.Println("id qty sname")
	//snameLookup.Print()
	fmt.Println("\n# Sex", sexLookup.LastKey())
	fmt.Println("id qty sex")
	sexLookup.Print()
	fmt.Println("\n# Country", countryLookup.LastKey())
	fmt.Println("id qty country")
	countryLookup.Print()
	//fmt.Println("\n# City", cityLookup.LastKey())
	//fmt.Println("id qty city")
	//cityLookup.Print()
	fmt.Println("\n# Status", statusLookup.LastKey())
	fmt.Println("id qty status")
	statusLookup.Print()
	fmt.Println("\n# Interest", interestLookup.LastKey())
	fmt.Println("id qty interest")
	interestLookup.Print()
}

const likeStructSize = 8
const likePageSize = 256

type likePage [likePageSize]byte

type account struct {
	domain       uint8
	fname        uint8
	sname        uint16
	sex          uint8
	country      uint8
	city         uint16
	status       uint8
	interestSize uint8
	interest     [10]uint8
	loginSize    uint8
	login        [24]byte
	likePages    [3]uint32
	likePage     likePage
}

type DB struct {
	a []account
	p []likePage
}

func NewDB() *DB {
	return &DB{
		a: make([]account, 1400000),
		p: make([]likePage, 900000),
	}
}

func (db *DB) insertAccount(src *proto.Account, checkLikes bool) error {
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
	n = len(src.Interests)
	if n > 10 {
		log.Println("Too much interests", n)
		return DefaultError
	}
	dst.interestSize = uint8(n)
	for i := 0; i < n; i++ {
		k, _ = interestLookup.GetOrGen(src.Interests[i])
		dst.interest[i] = uint8(k)
	}
	n = len(src.Likes)
	likeCounters[n/(likePageSize/likeStructSize)]++
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
				a.Reset()
				t, ok = parse.ParseSymbol(t, ',')
				t, ok = a.UnmarshalJSON(t)
				if !ok {
					break
				}
				err := db.insertAccount(a, false)
				if err != nil {
					log.Fatal(err)
				}
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

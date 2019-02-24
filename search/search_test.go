package search

import (
	"log"
	"testing"

	"github.com/freepk/hlc18r4/accounts"
	"github.com/freepk/hlc18r4/backup"
	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/iterator"
)

var (
	accountsSvc *accounts.AccountsService
	searchSvc   *SearchService
)

func init() {
	log.Println("Restoring")
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts service")
	accountsSvc = accounts.NewAccountsService(rep)
	log.Println("Search service")
	searchSvc = NewSearchService(rep)
	searchSvc.Rebuild()
}

func TestSearch(t *testing.T) {
	countryName := []byte(`Геризия`)
	countryKey, ok := tokens.Country(countryName)
	if !ok {
		t.Fatal("Invalid country")
	}
	interestName := []byte(`Симпсоны`)
	interestKey, ok := tokens.Interest(interestName)
	if !ok {
		t.Fatal("Invalid interest")
	}
	t.Log(string(countryName), countryKey, string(interestName), interestKey)
	country := searchSvc.Countries(countryKey)
	iter := country.Interest(interestKey)
	count := 0
	for {
		id, ok := iter.Next()
		if !ok {
			break
		}
		_ = id
		count++
	}
	t.Log("count", count)
}

func BenchmarkTest0(b *testing.B) {
	countryKey, ok := tokens.Country([]byte(`Испезия`))
	if !ok {
		log.Fatal("Invalid country")
	}
	country := searchSvc.Countries(countryKey)
	if country == nil {
		log.Fatal("Country index is null")
	}
	interestKey, ok := tokens.Interest([]byte(`Обнимашки`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter := iterator.Iterator(country.Interest(interestKey))
	interestKey, ok = tokens.Interest([]byte(`YouTube`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter = iterator.NewUnionIter(iter, country.Interest(interestKey))
	interestKey, ok = tokens.Interest([]byte(`Солнце`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter = iterator.NewUnionIter(iter, country.Interest(interestKey))
	acc := &proto.Account{}
	buf := make([]byte, 0, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limit := 22
		iter.Reset()
		buf = buf[:0]
		for {
			if limit == 0 {
				break
			}
			id, ok := iter.Next()
			if !ok {
				break
			}
			*acc = *accountsSvc.Get(2000000 - id)
			if acc.Sex != tokens.FemaleSex {
				continue
			}
			buf = acc.MarshalJSON((proto.IDField | proto.EmailField | proto.SexField | proto.CountryField), buf)
			limit--
		}
	}
}

func BenchmarkTest1(b *testing.B) {
	interestKey, ok := tokens.Interest([]byte(`South Park`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter := iterator.Iterator(searchSvc.Common().Interest(interestKey))
	acc := &proto.Account{}
	buf := make([]byte, 0, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limit := 24
		iter.Reset()
		buf = buf[:0]
		for {
			if limit == 0 {
				break
			}
			id, ok := iter.Next()
			if !ok {
				break
			}
			*acc = *accountsSvc.Get(2000000 - id)
			if acc.Status != tokens.ComplStatus {
				continue
			}
			buf = acc.MarshalJSON((proto.IDField | proto.EmailField | proto.StatusField), buf)
			limit--
		}
	}
}

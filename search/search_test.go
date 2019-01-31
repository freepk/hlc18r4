package search

import (
	"testing"

	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/tokens"
)

func TestSearch(t *testing.T) {
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	svc := NewSearchService(rep)
	svc.Rebuild()
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
	country := svc.Countries(countryKey)
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
	t.Log(count)
}

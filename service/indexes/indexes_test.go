package indexes

import (
	"testing"

	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/backup"
)

func TestDefaultIndex(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Create default index")
	index := NewDefaultIndex(rep)
	t.Log("Rebuild default")
	index.Rebuild()

	country, ok := GetCountryToken([]byte(`Мализия`))
	if !ok {
		t.Fail()
	}
	it := iterator.Iterator(index.CountryIter(country))
	it = iterator.NewInterIter(it, index.SexIter(MaleToken))
	limit := 32
	for limit > 0 {
		limit--
		pseudo, ok := it.Next()
		if !ok {
			break
		}
		id := 2000000 - pseudo
		t.Log(limit, id)
	}
}

func BenchmarkDefaultIndex(b *testing.B) {
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		b.Fatal(err)
	}
	defaultIdx := NewDefaultIndex(rep)
	defaultIdx.Rebuild()
	country, ok := GetCountryToken([]byte(`Мализия`))
	if !ok {
		b.Fail()
	}
	it := iterator.Iterator(defaultIdx.CountryIter(country))
	it = iterator.NewInterIter(it, defaultIdx.SexIter(MaleToken))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it.Reset()
		limit := 32
		for limit > 0 {
			limit--
			pseudo, ok := it.Next()
			if !ok {
				break
			}
			id := 2000000 - pseudo
			_ = id
		}
	}
}

func BenchmarkCountryIndex(b *testing.B) {
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		b.Fatal(err)
	}

	defaultIdx := NewDefaultIndex(rep)
	defaultIdx.Rebuild()

	countryIdx := NewCountryIndex(rep)
	countryIdx.Rebuild()

	country, ok := GetCountryToken([]byte(`Мализия`))
	if !ok {
		b.Fail()
	}
	it := iterator.Iterator(defaultIdx.CountryIter(country))
	it = iterator.NewInterIter(it, countryIdx.SexIter(country, MaleToken))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it.Reset()
		limit := 32
		for limit > 0 {
			limit--
			pseudo, ok := it.Next()
			if !ok {
				break
			}
			id := 2000000 - pseudo
			_ = id
		}
	}
}

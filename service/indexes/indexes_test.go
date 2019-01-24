package indexes

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/proto"
	"testing"
)

func TestDefaultIndex(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Create index")
	index := NewDefaultIndex(rep)
	t.Log("Rebuild")
	index.Rebuild()

	country, ok := proto.CountryToken([]byte(`Мализия`))
	if !ok {
		t.Fail()
	}
	it := iterator.Iterator(index.Country(country))
	it = iterator.NewInterIter(it, index.Sex(MaleToken))
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

func TestCountryIndex(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Create index")
	defaultIdx := NewDefaultIndex(rep)
	t.Log("Rebuild")
	defaultIdx.Rebuild()

	t.Log("Create index")
	countryIdx := NewCountryIndex(rep)
	t.Log("Rebuild")
	countryIdx.Rebuild()

	country, ok := proto.CountryToken([]byte(`Мализия`))
	if !ok {
		t.Fail()
	}
	it := iterator.Iterator(defaultIdx.Country(country))
	it = iterator.NewInterIter(it, countryIdx.Sex(country, MaleToken))

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
	country, ok := proto.CountryToken([]byte(`Мализия`))
	if !ok {
		b.Fail()
	}
	it := iterator.Iterator(defaultIdx.Country(country))
	it = iterator.NewInterIter(it, defaultIdx.Sex(MaleToken))

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

	country, ok := proto.CountryToken([]byte(`Мализия`))
	if !ok {
		b.Fail()
	}
	it := iterator.Iterator(defaultIdx.Country(country))
	it = iterator.NewInterIter(it, countryIdx.Sex(country, MaleToken))

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
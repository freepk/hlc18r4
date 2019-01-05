package main

import (
	"testing"
)

func TestCountry(t *testing.T) {
	for i := 0; i < 70; i++ {
		item := countryLookup.GetItem(i)
		if item == nil {
			t.Fail()
		}
		index := countryLookup.GetIndex(item)
		if index != i {
			t.Fail()
		}
	}
}

func BenchmarkCountryGetItemNoLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countryLookup.GetItemNoLock(55)
	}
}

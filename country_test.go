package main

import (
	"testing"
)

func TestCountry(t *testing.T) {
	for i := 1; i <= 70; i++ {
		name := getCountryName(i)
		if getCountryId(name) != i {
			t.Log(name)
			t.Fail()
		}
	}
}

func BenchmarkGetCountryName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryName(55)
	}
}

func BenchmarkGetCountryId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryId("Алания")
	}
}

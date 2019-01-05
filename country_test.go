package main

import (
	"testing"
)

func BenchmarkGetCountryName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryName(55)
	}
}

func BenchmarkGetCountryId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryId("Росция")
	}
}

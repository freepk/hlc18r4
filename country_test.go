package main

import (
	"testing"
)

var country = ""

func BenchmarkGetCountry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		country = getCountry(55)
	}
}

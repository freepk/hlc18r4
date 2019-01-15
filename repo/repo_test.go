package repo

import (
	"testing"
)

func TestAccountsRepo(t *testing.T) {
	rep := NewAccountsRepo(1024)
	if rep.Exists(512) {
		t.Fail()
	}
	if rep.Exists(4096) {
		t.Fail()
	}
}

func BenchmarkAccountsRepoExists(b *testing.B) {
	rep := NewAccountsRepo(1024)
	for i := 0; i < b.N; i++ {
		rep.Exists(512)
	}
}

func BenchmarkAccountsGet(b *testing.B) {
	rep := NewAccountsRepo(1024)
	for i := 0; i < b.N; i++ {
		_, _ = rep.Get(512)
	}
}

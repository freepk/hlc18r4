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
	if ok := rep.Add(101, &Account{Joined: 1, Birth: 1, Status: BusyStatus, Email: "test@mail.ru"}); !ok {
		t.Fail()
	}
	if ok := rep.Add(102, &Account{Joined: 1, Birth: 1, Status: BusyStatus, Email: "test@mail.ru"}); ok {
		t.Fail()
	}
	if ok := rep.Add(101, &Account{Joined: 1, Birth: 1, Status: BusyStatus, Email: "test1@mail.ru"}); ok {
		t.Fail()
	}
}

func BenchmarkAccountsRepoExists(b *testing.B) {
	rep := NewAccountsRepo(1024)
	rep.Add(101, &Account{Joined: 1, Birth: 1, Status: BusyStatus, Email: "test@mail.ru"})
	for i := 0; i < b.N; i++ {
		_ = rep.Exists(101)
	}
}

func BenchmarkAccountsGet(b *testing.B) {
	rep := NewAccountsRepo(1024)
	rep.Add(101, &Account{Joined: 1, Birth: 1, Status: BusyStatus, Email: "test@mail.ru"})
	for i := 0; i < b.N; i++ {
		rep.Get(101)
	}
}

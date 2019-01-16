package repo

import (
	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/proto"
)

type AccountsRepo struct {
	emails   *hashtab.HashTab
	accounts []proto.Account
}

func NewAccountsRepo(size int) *AccountsRepo {
	emails := hashtab.NewHashTab(size)
	accounts := make([]proto.Account, size)
	return &AccountsRepo{emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) Exists(id int) bool {
	if id < 1 || id > len(rep.accounts)-1 {
		return false
	}
	return (rep.accounts[id].Joined > 0)
}

func (rep *AccountsRepo) Get(id int) (*proto.Account, bool) {
	if id < 1 || id > len(rep.accounts)-1 {
		return nil, false
	}
	acc := rep.accounts[id]
	if acc.Joined == 0 {
		return nil, false
	}
	return &acc, true
}

func (rep *AccountsRepo) validate(acc *proto.Account) bool {
	if acc.Email == "" {
		return false
	}
	if acc.Joined == 0 {
		return false
	}
	if acc.Birth == 0 {
		return false
	}
	return true
}

func (rep *AccountsRepo) set(id int, acc *proto.Account, checkExists bool) bool {
	if id < 1 || id > len(rep.accounts)-1 {
		return false
	}
	if valid := rep.validate(acc); !valid {
		return false
	}
	if checkExists && rep.accounts[id].Joined > 0 {
		return false
	}
	hash := murmur3.Sum64([]byte(acc.Email))
	owner, ok := rep.emails.GetOrSet(hash, uint64(id))
	if ok && owner != uint64(id) {
		return false
	}
	rep.accounts[id] = *acc
	return true
}

func (rep *AccountsRepo) Add(id int, acc *proto.Account) bool {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *proto.Account) bool {
	return rep.set(id, acc, false)
}

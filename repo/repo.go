package repo

import (
	"gitlab.com/freepk/hlc18r4/proto"
)

type AccountsRepo struct {
	accounts []proto.Account
}

type ForEachFunc func(id int, acc *proto.Account)

func NewAccountsRepo(num int) *AccountsRepo {
	accounts := make([]proto.Account, num)
	return &AccountsRepo{accounts: accounts}
}

func (rep *AccountsRepo) Get(id int) *proto.Account {
	if id > 0 && id < len(rep.accounts) {
		acc := rep.accounts[id]
		return &acc
	}
	return nil
}

func (rep *AccountsRepo) Set(id int, acc *proto.Account) bool {
	if id > 0 && id < len(rep.accounts) {
		rep.accounts[id] = *acc
		return true
	}
	return false
}

func (rep *AccountsRepo) Forward(handler ForEachFunc) {
	acc := &proto.Account{}
	for id := range rep.accounts {
		*acc = rep.accounts[id]
		if acc.Email.Len > 0 {
			handler(id, acc)
		}
	}
}

func (rep *AccountsRepo) Reverse(handler ForEachFunc) {
	acc := &proto.Account{}
	last := len(rep.accounts) - 1
	for id := range rep.accounts {
		*acc = rep.accounts[last-id]
		if acc.Email.Len > 0 {
			handler(last-id, acc)
		}
	}
}

func (rep *AccountsRepo) Len() int {
	return len(rep.accounts)
}

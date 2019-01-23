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
	if id < len(rep.accounts) {
		acc := rep.accounts[id]
		return &acc
	}
	return nil
}

func (rep *AccountsRepo) Set(id int, acc *proto.Account) bool {
	if id < len(rep.accounts) {
		rep.accounts[id] = *acc
		return true
	}
	return false
}

func (rep *AccountsRepo) Len() int {
	return len(rep.accounts)
}

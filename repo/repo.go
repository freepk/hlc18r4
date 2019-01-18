package repo

import (
	"gitlab.com/freepk/hlc18r4/proto"
)

type AccountsRepo struct {
	accounts []proto.Account
}

type ForEachFunc func(acc *proto.Account)

func NewAccountsRepo(num int) *AccountsRepo {
	accounts := make([]proto.Account, num)
	return &AccountsRepo{accounts: accounts}
}

func (rep *AccountsRepo) Get(id int) *proto.Account {
	if id > 0 && id < len(rep.accounts) {
		return &rep.accounts[id]
	}
	return nil
}

func (rep *AccountsRepo) Add(acc *proto.Account) {
	if dst := rep.Get(int(acc.ID)); dst != nil {
		*dst = *acc
	}
}

func (rep *AccountsRepo) ForEach(handler ForEachFunc) {
	for id := range rep.accounts {
		acc := &rep.accounts[id]
		if acc.ID > 0 {
			handler(acc)
		}
	}
}

func (rep *AccountsRepo) Size() int {
	return len(rep.accounts)
}

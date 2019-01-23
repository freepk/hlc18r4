package service

import (
	"gitlab.com/freepk/hlc18r4/index"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsIndexer struct {
	last int
	rep  *repo.AccountsRepo
}

func NewAccountsIndexer(rep *repo.AccountsRepo) *AccountsIndexer {
	return &AccountsIndexer{rep: rep}
}

func (ix *AccountsIndexer) Reset() {
	ix.last = 0
}

func (ix *AccountsIndexer) Next() (*index.Item, bool) {
	n := ix.rep.Len()
	for i := ix.last; i < n; i++ {
		acc := ix.rep.Get(n - i - 1)
		if acc.Email.Len > 0 {
		}
	}
	return nil, false
}

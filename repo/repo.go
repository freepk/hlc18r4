package repo

import (
	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

type AccountsRepo struct {
	size     uint32
	emails   *hashtab.HashTab
	accounts []Account
}

func NewAccountsRepo(size uint32) *AccountsRepo {
	emails := hashtab.NewHashTab(size)
	accounts := make([]Account, size)
	return &AccountsRepo{size: size, emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) Exists(id int) bool {
	if id >= int(rep.size) {
		return false
	}
	return (rep.accounts[id].Joined > 0)
}

func (rep *AccountsRepo) Get(id int) (*Account, bool) {
	if id >= int(rep.size) {
		return nil, false
	}
	acc := rep.accounts[id]
	if acc.Joined == 0 {
		return nil, false
	}
	return &acc, true
}

func (rep *AccountsRepo) validate(acc *Account) bool {
	if acc.Email == "" {
		return false
	}
	if acc.Joined == 0 {
		return false
	}
	if acc.Birth == 0 {
		return false
	}
	if acc.Status != FreeStatus && acc.Status != BusyStatus && acc.Status != ComplicatedStatus {
		return false
	}
	if acc.PremiumFinish > 0 && acc.PremiumPeriod != MonthPeriod && acc.PremiumPeriod != QuarterPeriod && acc.PremiumPeriod != HalfYearPeriod {
		return false
	}
	return true
}

func (rep *AccountsRepo) set(id int, acc *Account, checkExists bool) bool {
	if id >= int(rep.size) {
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

func (rep *AccountsRepo) Add(id int, acc *Account) bool {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *Account) bool {
	return rep.set(id, acc, false)
}

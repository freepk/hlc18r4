package repo

import (
	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/proto"
)

type AccountsRepo struct {
	size     uint32
	emails   *hashtab.HashTab
	accounts []proto.Account
}

func NewAccountsRepo(size uint32) *AccountsRepo {
	emails := hashtab.NewHashTab(size)
	accounts := make([]proto.Account, size)
	return &AccountsRepo{size: size, emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) Exists(id int) bool {
	if id >= int(rep.size) {
		return false
	}
	return (rep.accounts[id].Joined > 0)
}

func (rep *AccountsRepo) Get(id int) (*proto.Account, bool) {
	if id >= int(rep.size) {
		return nil, false
	}
	acc := rep.accounts[id]
	if acc.Joined == 0 {
		return nil, false
	}
	return &acc, true
}

func (rep *AccountsRepo) validate(acc *proto.Account) bool {
	/*
		if acc.Email == "" {
			return false
		}
		if acc.Joined == 0 {
			return false
		}
		if acc.Birth == 0 {
			return false
		}
		if acc.Status != proto.FreeStatus &&
			acc.Status != proto.BusyStatus &&
			acc.Status != proto.ComplicatedStatus {
			return false
		}
		if acc.PremiumFinish > 0 &&
			acc.PremiumPeriod != proto.MonthPeriod &&
			acc.PremiumPeriod != proto.QuarterPeriod &&
			acc.PremiumPeriod != proto.HalfYearPeriod {
			return false
		}
	*/
	return true
}

func (rep *AccountsRepo) set(id int, acc *proto.Account, checkExists bool) bool {
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

func (rep *AccountsRepo) Add(id int, acc *proto.Account) bool {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *proto.Account) bool {
	return rep.set(id, acc, false)
}

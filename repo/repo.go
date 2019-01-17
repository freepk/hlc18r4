package repo

import (
	"bytes"
	"log"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/proto"
)

type AccountsRepo struct {
	emails   *hashtab.HashTab
	accounts []proto.Account
}

func NewAccountsRepo(num int) *AccountsRepo {
	emails := hashtab.NewHashTab(num * 120 / 100)
	log.Println("Emails", emails.Size())
	accounts := make([]proto.Account, num)
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
	if len(acc.Email) == 0 || bytes.IndexByte(acc.Email, '@') == -1 {
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
	hash := murmur3.Sum64(acc.Email)
	owner, ok := rep.emails.GetOrSet(hash, uint64(id))
	if ok && owner != uint64(id) {
		return false
	}
	dst := &rep.accounts[id]
	dst.ID = acc.ID
	dst.Birth = acc.Birth
	dst.Joined = acc.Joined
	dst.Email = append(dst.Email[:0], acc.Email...)
	dst.Fname = acc.Fname
	dst.Sname = acc.Sname
	dst.Phone = append(dst.Phone[:0], acc.Phone...)
	dst.Sex = acc.Sex
	dst.Country = acc.Country
	dst.City = acc.City
	dst.Status = acc.Status
	dst.Interests = append(dst.Interests[:0], acc.Interests...)
	return true
}

func (rep *AccountsRepo) Add(id int, acc *proto.Account) bool {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *proto.Account) bool {
	return rep.set(id, acc, false)
}

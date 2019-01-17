package service

import (
	"bytes"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep    *repo.AccountsRepo
	emails *hashtab.HashTab
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	emails := hashtab.NewHashTab(rep.Size())
	rep.ForEach(func(acc *proto.Account) {
		if len(acc.Email) > 0 {
			hash := murmur3.Sum64(acc.Email)
			emails.Set(hash, uint64(acc.ID))
		}
	})
	return &AccountsService{rep: rep, emails: emails}
}

func (svc *AccountsService) Exists(id int) bool {
	acc := svc.rep.Get(id)
	if acc == nil {
		return false
	}
	return (acc.Joined > 0)
}

func (svc *AccountsService) Create(acc *proto.Account) bool {
	id := int(acc.ID)
	if id == 0 || svc.Exists(id) {
		return false
	}
	if len(acc.Email) == 0 || bytes.IndexByte(acc.Email, 0x40) == -1 {
		return false
	}
	// hold new
	hash := murmur3.Sum64(acc.Email)
	if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
		return false
	}
	dst := acc.Clone()
	svc.rep.Add(dst)
	return true
}

func (svc *AccountsService) Update(id int, acc *proto.Account) bool {
	if len(acc.Email) > 0 && bytes.IndexByte(acc.Email, 0x40) == -1 {
		return false
	}
	dst := svc.rep.Get(id)
	if dst == nil {
		return false
	}
	if len(acc.Email) > 0 {
		// hold new
		hash := murmur3.Sum64(acc.Email)
		if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
			return false
		}
		// release old
		// hash = murmur3.Sum64(dst.Email)
		// svc.emails.Del(hash)
		dst.Email = append(dst.Email[:0], acc.Email...)
	}
	if len(acc.Phone) > 0 {
		dst.Phone = append(dst.Phone[:0], acc.Phone...)
	}
	if len(acc.Interests) > 0 {
		dst.Interests = append(dst.Interests[:0], acc.Interests...)
	}
	// etc...
	return true
}

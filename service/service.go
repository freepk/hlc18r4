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
		if acc.EmailSize > 0 {
			hash := murmur3.Sum64(acc.Email[:acc.EmailSize])
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
	if acc.EmailSize == 0 || bytes.IndexByte(acc.Email[:acc.EmailSize], 0x40) == -1 {
		return false
	}
	// hold new
	hash := murmur3.Sum64(acc.Email[:acc.EmailSize])
	if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
		return false
	}
	dst := acc.Clone()
	svc.rep.Add(dst)
	return true
}

func (svc *AccountsService) Update(id int, acc *proto.Account) bool {
	if acc.EmailSize > 0 && bytes.IndexByte(acc.Email[:acc.EmailSize], 0x40) == -1 {
		return false
	}
	dst := svc.rep.Get(id)
	if dst == nil {
		return false
	}
	if acc.EmailSize > 0 {
		// hold new
		hash := murmur3.Sum64(acc.Email[:acc.EmailSize])
		if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
			return false
		}
		dst.EmailSize = acc.EmailSize
		dst.Email = acc.Email
	}
	if acc.PhoneSize > 0 {
		dst.PhoneSize = acc.PhoneSize
		dst.Phone = acc.Phone
	}
	if len(acc.Interests) > 0 {
		dst.Interests = append(dst.Interests[:0], acc.Interests...)
	}
	if len(acc.LikesTo) > 0 {
		dst.LikesTo = append(dst.LikesTo[:0], acc.LikesTo...)
	}
	// etc...
	return true
}

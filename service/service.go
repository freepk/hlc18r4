package service

import (
	"bytes"

	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep *repo.AccountsRepo
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	return &AccountsService{rep: rep}
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
		dst.Email = append(dst.Email, acc.Email...)
	}
	if len(acc.Phone) > 0 {
		dst.Phone = append(dst.Phone, acc.Phone...)
	}
	if len(acc.Interests) > 0 {
		dst.Interests = append(dst.Interests, acc.Interests...)
	}
	// etc...
	return true
}

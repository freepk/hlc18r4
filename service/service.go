package service

import (
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
	return svc.rep.Exists(id)
}

func (svc *AccountsService) Create(id int, src *proto.Account) bool {
	return svc.rep.Add(id, src)
}

func (svc *AccountsService) Update(id int, src *proto.Account) bool {
	dst, ok := svc.rep.Get(id)
	if !ok {
		return false
	}
	if len(src.Email) > 0 {
		dst.Email = src.Email
	}
	if len(src.Phone) > 0 {
		dst.Phone = src.Phone
	}
	// etc...
	return svc.rep.Set(id, dst)
}

package service

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep *repo.AccountsRepo
}

func NewAccountsService() *AccountsService {
	rep := repo.NewAccountsRepo()
	return &AccountsService{rep: rep}
}

func (svc *AccountsService) Add(acc *proto.Account) error {
	tmp := &repo.Account{}
	tmp.Email = string(acc.Email)
	tmp.Birth = uint32(acc.Birth)
	tmp.Joined = uint32(acc.Joined)
	tmp.Status = repo.FreeStatus
	return svc.rep.Add(acc.ID, tmp)
}

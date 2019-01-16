package service

import (
	"bytes"
	"errors"

	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

var (
	AccountsServiceEmailError = errors.New("Account Email error")
	AccountsServiceSexError   = errors.New("Account Sex error")
)

type AccountsService struct {
	rep *repo.AccountsRepo
}

func NewAccountsService() *AccountsService {
	rep := repo.NewAccountsRepo()
	return &AccountsService{rep: rep}
}

func (svc *AccountsService) Create(acc *proto.Account) error {
	if len(acc.Email) == 0 || bytes.IndexByte(acc.Email, '@') == -1 {
		return AccountsServiceEmailError
	}
	if len(acc.Sex) > 1 || (len(acc.Sex) == 1 && acc.Sex[0] != 'm' && acc.Sex[0] != 'f') {
		return AccountsServiceSexError
	}
	tmp := &repo.Account{}
	tmp.Email = string(acc.Email)
	tmp.Birth = uint32(acc.Birth)
	tmp.Joined = uint32(acc.Joined)
	tmp.Status = repo.FreeStatus
	return svc.rep.Add(acc.ID, tmp)
}

func (svc *AccountsService) Exists(id int) bool {
	return svc.rep.Exists(id)
}

func (svc *AccountsService) Update(id int, acc *proto.Account) error {
	if len(acc.Email) > 0 && bytes.IndexByte(acc.Email, '@') == -1 {
		return AccountsServiceEmailError
	}
	if len(acc.Sex) > 1 || (len(acc.Sex) == 1 && acc.Sex[0] != 'm' && acc.Sex[0] != 'f') {
		return AccountsServiceSexError
	}
	tmp, err := svc.rep.Get(id)
	if err != nil {
		return err
	}
	if len(acc.Email) > 0 {
		tmp.Email = string(acc.Email)
	}
	return svc.rep.Set(id, tmp)
}

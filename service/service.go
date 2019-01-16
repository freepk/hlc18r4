package service

import (
	"bytes"

	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep *repo.AccountsRepo
}

func NewAccountsService(size uint32) *AccountsService {
	rep := repo.NewAccountsRepo(size)
	return &AccountsService{rep: rep}
}

func (svc *AccountsService) Exists(id int) bool {
	return svc.rep.Exists(id)
}

func (svc *AccountsService) Create(acc *proto.Account) bool {
	if len(acc.Email) == 0 || bytes.IndexByte(acc.Email, '@') == -1 {
		return false
	}
	if len(acc.Sex) > 1 || (len(acc.Sex) == 1 && acc.Sex[0] != 'm' && acc.Sex[0] != 'f') {
		return false
	}
	if len(acc.Status) > 0 &&
		string(acc.Status) != `\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e` &&
		string(acc.Status) != `\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b` &&
		string(acc.Status) != `\u0437\u0430\u043d\u044f\u0442\u044b` {
		return false
	}
	tmp := &repo.Account{}
	tmp.Email = string(acc.Email)
	tmp.Birth = uint32(acc.Birth)
	tmp.Joined = uint32(acc.Joined)
	tmp.Status = repo.FreeStatus
	return svc.rep.Add(acc.ID, tmp)
}

func (svc *AccountsService) Update(id int, acc *proto.Account) bool {
	if len(acc.Email) > 0 && bytes.IndexByte(acc.Email, '@') == -1 {
		return false
	}
	if len(acc.Sex) > 1 || (len(acc.Sex) == 1 && acc.Sex[0] != 'm' && acc.Sex[0] != 'f') {
		return false
	}
	if len(acc.Status) > 0 &&
		string(acc.Status) != `\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e` &&
		string(acc.Status) != `\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b` &&
		string(acc.Status) != `\u0437\u0430\u043d\u044f\u0442\u044b` {
		return false
	}
	tmp, ok := svc.rep.Get(id)
	if !ok {
		return false
	}
	if len(acc.Email) > 0 {
		tmp.Email = string(acc.Email)
	}
	return svc.rep.Set(id, tmp)
}

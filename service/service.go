package service

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	repo *repo.AccountsRepo
}

func (svc *AccountsService) Insert(src *proto.Account) error {
	return nil
}

func (svc *AccountsService) Update(src *proto.Account) error {
	return nil
}

func (svc *AccountsService) Exists(id int) bool {
	return true
}

func (svc *AccountsService) Filter() {
}

func (svc *AccountsService) Group() {
}

func (svc *AccountsService) Suggest() {
}

func (svc *AccountsService) Recommend() {
}

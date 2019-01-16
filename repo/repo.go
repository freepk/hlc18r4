package repo

import (
	"errors"
)

var (
	AccountsRepoExistsError    = errors.New("Account exists")
	AccountsRepoNotExistsError = errors.New("Account not exists")
	AccountsRepoEmailError     = errors.New("Account Email error")
	AccountsRepoJonedError     = errors.New("Account Joined error")
	AccountsRepoBirthError     = errors.New("Account Birth error")
	AccountsRepoStatusError    = errors.New("Account Status error")
	AccountsRepoPremiumError   = errors.New("Account Premium error")
)

type AccountsRepo struct {
	emails   map[string]int
	accounts map[int]Account
}

func NewAccountsRepo() *AccountsRepo {
	emails := make(map[string]int)
	accounts := make(map[int]Account)
	return &AccountsRepo{emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) exists(id int) bool {
	_, ok := rep.accounts[id]
	return ok
}

func (rep *AccountsRepo) Exists(id int) bool {
	return rep.exists(id)
}

func (rep *AccountsRepo) Get(id int) (*Account, error) {
	if !rep.exists(id) {
		return nil, AccountsRepoNotExistsError
	}
	acc := rep.accounts[id]
	return &acc, nil
}

func (rep *AccountsRepo) validate(acc *Account) error {
	if acc.Email == "" {
		return AccountsRepoEmailError
	}
	if acc.Joined == 0 {
		return AccountsRepoJonedError
	}
	if acc.Birth == 0 {
		return AccountsRepoBirthError
	}
	if acc.Status != FreeStatus && acc.Status != BusyStatus && acc.Status != ComplicatedStatus {
		return AccountsRepoStatusError
	}
	if acc.PremiumFinish > 0 && acc.PremiumPeriod != MonthPeriod && acc.PremiumPeriod != QuarterPeriod && acc.PremiumPeriod != HalfYearPeriod {
		return AccountsRepoPremiumError
	}
	return nil
}

func (rep *AccountsRepo) set(id int, acc *Account, checkExists bool) error {
	if err := rep.validate(acc); err != nil {
		return err
	}
	if eid, ok := rep.emails[acc.Email]; ok && eid != id {
		return AccountsRepoEmailError
	}
	if checkExists && rep.exists(id) {
		return AccountsRepoExistsError
	}
	rep.emails[acc.Email] = id
	rep.accounts[id] = *acc
	return nil
}

func (rep *AccountsRepo) Add(id int, acc *Account) error {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *Account) error {
	return rep.set(id, acc, false)
}

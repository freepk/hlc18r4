package repo

import (
	"errors"
	"sync"
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
	sync.RWMutex
	emails   map[string]int
	accounts map[int]Account
}

func NewAccountsRepo() *AccountsRepo {
	emails := make(map[string]int)
	accounts := make(map[int]Account)
	return &AccountsRepo{emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) Exists(id int) bool {
	rep.RLock()
	_, ok := rep.accounts[id]
	rep.RUnlock()
	return ok
}

func (rep *AccountsRepo) Get(id int) (*Account, error) {
	rep.RLock()
	acc, ok := rep.accounts[id]
	rep.RUnlock()
	if !ok {
		return nil, AccountsRepoNotExistsError
	}
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
	rep.Lock()
	if checkExists {
		if _, ok := rep.accounts[id]; ok {
			rep.Unlock()
			return AccountsRepoExistsError
		}
	}
	if emailID, ok := rep.emails[acc.Email]; ok && emailID != id {
		rep.Unlock()
		return AccountsRepoEmailError
	}
	rep.emails[acc.Email] = id
	rep.accounts[id] = *acc
	rep.Unlock()
	return nil
}

func (rep *AccountsRepo) Add(id int, acc *Account) error {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *Account) error {
	return rep.set(id, acc, false)
}

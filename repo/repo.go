package repo

import (
	"errors"
)

const (
	accountsBucketSize = 10000
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
	emails   map[string]uint32
	accounts map[uint32]Account
}

func NewAccountsRepo() *AccountsRepo {
	emails := make(map[string]uint32)
	accounts := make(map[uint32]Account)
	return &AccountsRepo{emails: emails, accounts: accounts}
}

func (rep *AccountsRepo) exists(id uint32) bool {
	_, ok := rep.accounts[id]
	return ok
}

func (rep *AccountsRepo) Exists(id uint32) bool {
	return rep.exists(id)
}

func (rep *AccountsRepo) Get(id uint32) (*Account, error) {
	if !rep.exists(id) {
		return nil, AccountsRepoNotExistsError
	}
	account := rep.accounts[id]
	return &account, nil
}

func (rep *AccountsRepo) validate(account *Account) error {
	if account.Email == "" {
		return AccountsRepoEmailError
	}
	if account.Joined == 0 {
		return AccountsRepoJonedError
	}
	if account.Birth == 0 {
		return AccountsRepoBirthError
	}
	if account.Status != FreeStatus && account.Status != BusyStatus && account.Status != ComplicatedStatus {
		return AccountsRepoStatusError
	}
	if account.PremiumFinish > 0 && account.PremiumPeriod != MonthPeriod && account.PremiumPeriod != QuarterPeriod && account.PremiumPeriod != HalfYearPeriod {
		return AccountsRepoPremiumError
	}
	return nil
}

func (rep *AccountsRepo) set(id uint32, account *Account, checkExists bool) error {
	if err := rep.validate(account); err != nil {
		return err
	}
	if eid, ok := rep.emails[account.Email]; ok && eid != id {
		return AccountsRepoEmailError
	}
	if checkExists && rep.exists(id) {
		return AccountsRepoExistsError
	}
	rep.accounts[id] = *account
	return nil
}

func (rep *AccountsRepo) Add(id uint32, account Account) error {
	return rep.set(id, &account, true)
}

func (rep *AccountsRepo) Set(id uint32, account Account) error {
	return rep.set(id, &account, false)
}

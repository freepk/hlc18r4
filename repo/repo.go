package repo

import (
	"errors"
	"sync"
)

const (
	accountsBucketSize = 10000
)

var (
	AccountsRepoOverflowError      = errors.New("Account ID overflow")
	AccountsRepoAlreadyExistsError = errors.New("Account already exists")
	AccountsRepoNotExistsError     = errors.New("Account not exists")
	AccountsRepoValidationError    = errors.New("Account validation error")
)

type AccountsRepo struct {
	sync.RWMutex
	accounts []Account
	length   uint32
	lastID   uint32
}

func NewAccountsRepo(size int) *AccountsRepo {
	accounts := make([]Account, size)
	return &AccountsRepo{accounts: accounts}
}

func (rep *AccountsRepo) exists(id uint32) bool {
	if id > rep.lastID {
		return false
	}
	return (rep.accounts[id].Joined > 0)
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
		return AccountsRepoValidationError
	}
	if account.Joined == 0 {
		return AccountsRepoValidationError
	}
	if account.Birth == 0 {
		return AccountsRepoValidationError
	}
	if account.Status != FreeStatus && account.Status != BusyStatus && account.Status != ComplicatedStatus {
		return AccountsRepoValidationError
	}
	if account.PremiumFinish > 0 && account.PremiumPeriod != MonthPeriod && account.PremiumPeriod != QuarterPeriod && account.PremiumPeriod != HalfYearPeriod {
		return AccountsRepoValidationError
	}
	return nil
}

func (rep *AccountsRepo) setLastID(id uint32) {
	if id > rep.lastID {
		rep.Lock()
		if id > rep.lastID {
			rep.lastID = id
		}
		rep.Unlock()
	}
}

func (rep *AccountsRepo) set(id uint32, account *Account, checkExists bool) error {
	if id > uint32(len(rep.accounts)) {
		return AccountsRepoOverflowError
	}
	if err := rep.validate(account); err != nil {
		return err
	}
	if checkExists && rep.exists(id) {
		return AccountsRepoAlreadyExistsError
	}
	rep.accounts[id] = *account
	rep.setLastID(id)
	return nil
}

func (rep *AccountsRepo) Add(id uint32, account Account) error {
	return rep.set(id, &account, true)
}

func (rep *AccountsRepo) Set(id uint32, account Account) error {
	return rep.set(id, &account, false)
}

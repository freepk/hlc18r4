package repo

import (
	"errors"
	"sync"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

var (
	AccountsRepoSizeOverflow   = errors.New("Size overflow")
	AccountsRepoExistsError    = errors.New("Account exists")
	AccountsRepoNotExistsError = errors.New("Account not exists")
	AccountsRepoEmailError     = errors.New("Account Email error")
	AccountsRepoJonedError     = errors.New("Account Joined error")
	AccountsRepoBirthError     = errors.New("Account Birth error")
	AccountsRepoStatusError    = errors.New("Account Status error")
	AccountsRepoPremiumError   = errors.New("Account Premium error")
)

const (
	accountsPerBucket = 10000
)

type AccountsRepo struct {
	sync.RWMutex
	emails   *hashtab.HashTab
	accounts []Account
	locks    []sync.RWMutex
}

func NewAccountsRepo(size uint32) *AccountsRepo {
	emails := hashtab.NewHashTab(size)
	accounts := make([]Account, size)
	locksSize := (size / accountsPerBucket) + 1
	locks := make([]sync.RWMutex, locksSize)
	return &AccountsRepo{emails: emails, accounts: accounts, locks: locks}
}

func (rep *AccountsRepo) lock(id int) *sync.RWMutex {
	bucket := id / accountsPerBucket
	return &rep.locks[bucket]
}

func (rep *AccountsRepo) Exists(id int) bool {
	if id >= len(rep.accounts) {
		return false
	}
	lock := rep.lock(id)
	lock.RLock()
	ok := (rep.accounts[id].Joined > 0)
	lock.RUnlock()
	return ok
}

func (rep *AccountsRepo) Get(id int) *Account {
	if id >= len(rep.accounts) {
		return nil
	}
	lock := rep.lock(id)
	lock.RLock()
	acc := &rep.accounts[id]
	lock.RUnlock()
	if acc.Joined > 0 {
		return acc
	}
	return nil
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
	if id >= len(rep.accounts) {
		return AccountsRepoSizeOverflow
	}
	if err := rep.validate(acc); err != nil {
		return err
	}
	hash := murmur3.Sum64([]byte(acc.Email))
	lock := rep.lock(id)
	lock.Lock()
	original := &rep.accounts[id]
	if checkExists && original.Joined > 0 {
		lock.Unlock()
		return AccountsRepoExistsError
	}
	owner, ok := rep.emails.Get(hash)
	if ok && owner != uint64(id) {
		lock.Unlock()
		return AccountsRepoEmailError
	}
	rep.emails.Set(hash, uint64(id))
	rep.accounts[id] = *acc
	lock.Unlock()
	return nil
}

func (rep *AccountsRepo) Add(id int, acc *Account) error {
	return rep.set(id, acc, true)
}

func (rep *AccountsRepo) Set(id int, acc *Account) error {
	return rep.set(id, acc, false)
}

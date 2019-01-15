package repo

import (
	"errors"
	"sync"
)

const (
	accountsBucketSize = 10000
)

var (
	AccountExistsError     = errors.New("Account already exists")
	AccountValidationError = errors.New("Account validation error")
)

type AccountsRepo struct {
	accounts    []Account
	length      uint32
	lastID      uint32
	bucketLocks []sync.RWMutex
}

func NewAccountsRepo(size int) *AccountsRepo {
	bucketLocksSize := (size / accountsBucketSize) + 1
	bucketLocks := make([]sync.RWMutex, bucketLocksSize)
	accounts := make([]Account, size)
	return &AccountsRepo{accounts: accounts, bucketLocks: bucketLocks}
}

func (rep *AccountsRepo) bucketLock(id uint32) *sync.RWMutex {
	bucket := id / accountsBucketSize
	return &rep.bucketLocks[bucket]
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

func (rep *AccountsRepo) Get(id uint32) (*Account, bool) {
	if rep.exists(id) {
		lock := rep.bucketLock(id)
		lock.RLock()
		account := rep.accounts[id]
		lock.RUnlock()
		return &account, true
	}
	return nil, false
}

func (rep *AccountsRepo) Add(id uint32, account *Account) error {
	if rep.exists(id) {
		return AccountExistsError
	}
	if account.Joined == 0 {
		return AccountValidationError
	}
	lock := rep.bucketLock(id)
	lock.Lock()
	rep.accounts[id] = *account
	lock.Unlock()
	return nil
}

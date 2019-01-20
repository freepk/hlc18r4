package service

import (
	"bytes"
	"log"
	"sync"

	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep        *repo.AccountsRepo
	emailsLock *sync.Mutex
	emails     map[uint64]int
	indexes    []*inverted.InvertedIndex
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	emailsLock := &sync.Mutex{}
	emails := make(map[uint64]int, rep.Len())
	rep.ForEach(func(id int, acc *proto.Account) {
		if acc.Email.Len > 0 {
			email := acc.Email.Buf[:acc.Email.Len]
			hash := murmur3.Sum64(email)
			emails[hash] = id
		}
	})
	return &AccountsService{rep: rep, emailsLock: emailsLock, emails: emails}
}

func (svc *AccountsService) AddInvertedIndex(parts inverted.PartsFunc, tokens inverted.TokensFunc) {
	index := inverted.NewInvertedIndex(svc.rep, parts, tokens)
	svc.indexes = append(svc.indexes, index)
}

func (svc *AccountsService) RebuildIndexes() {
	num := len(svc.indexes)
	wait := &sync.WaitGroup{}
	wait.Add(num)
	for i := range svc.indexes {
		index := svc.indexes[i]
		go func() {
			defer wait.Done()
			index.Rebuild()
		}()
	}
	wait.Wait()
}

func (svc *AccountsService) Exists(id int) bool {
	acc := svc.rep.Get(id)
	if acc == nil {
		return false
	}
	return (acc.Email.Len > 0)
}

func (svc *AccountsService) Create(id int, acc *proto.Account) bool {
	if id == 0 || svc.Exists(id) {
		return false
	}
	if acc.Email.Len == 0 || bytes.IndexByte(acc.Email.Buf[:acc.Email.Len], 0x40) == -1 {
		return false
	}
	// hold new
	hash := murmur3.Sum64(acc.Email.Buf[:acc.Email.Len])
	svc.emailsLock.Lock()
	if _, ok := svc.emails[hash]; ok {
		svc.emailsLock.Unlock()
		return false
	}
	svc.emails[hash] = id
	svc.emailsLock.Unlock()
	tmp := *acc
	tmp.LikesTo = make([]proto.Like, len(acc.LikesTo))
	copy(tmp.LikesTo, acc.LikesTo)
	svc.rep.Set(id, &tmp)
	return true
}

func (svc *AccountsService) Update(id int, acc *proto.Account) bool {
	if acc.Email.Len > 0 && bytes.IndexByte(acc.Email.Buf[:acc.Email.Len], 0x40) == -1 {
		return false
	}
	tmp := svc.rep.Get(id)
	if tmp == nil {
		return false
	}
	// ID
	// Birth
	// Joined
	if acc.ID[0] > 0 || acc.Birth[0] > 0 || acc.Joined[0] > 0 {
		log.Fatal("Update ID, Birth, Joined")
	}
	if acc.Email.Len > 0 {
		// hold new
		hash := murmur3.Sum64(acc.Email.Buf[:acc.Email.Len])
		svc.emailsLock.Lock()
		if _, ok := svc.emails[hash]; ok {
			svc.emailsLock.Unlock()
			return false
		}
		svc.emails[hash] = id
		svc.emailsLock.Unlock()
		tmp.Email = acc.Email
	}
	if acc.Fname > 0 {
		tmp.Fname = acc.Fname
	}
	if acc.Sname > 0 {
		tmp.Sname = acc.Sname
	}
	if acc.Phone[0] > 0 {
		tmp.Phone = acc.Phone
	}
	// Sex
	if acc.Sex > 0 {
		log.Fatal("Update Sex")
	}
	if acc.Country > 0 {
		tmp.Country = acc.Country
	}
	if acc.City > 0 {
		tmp.City = acc.City
	}
	if acc.Status > 0 {
		tmp.Status = acc.Status
	}
	if acc.PremiumStart[0] > 0 {
		tmp.PremiumStart = acc.PremiumStart
	}
	if acc.PremiumFinish[0] > 0 {
		tmp.PremiumFinish = acc.PremiumFinish
	}
	if acc.Interests[0] > 0 {
		tmp.Interests = acc.Interests
	}
	if len(acc.LikesTo) > 0 {
		tmp.LikesTo = append(tmp.LikesTo[:0], acc.LikesTo...)
	}
	svc.rep.Set(id, tmp)
	return true
}

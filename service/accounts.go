package service

import (
	"bytes"
	"sync"

	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep          *repo.AccountsRepo
	accountsPool *sync.Pool
	emailsLock   *sync.Mutex
	emails       map[uint64]int
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	accountsPool := &sync.Pool{New: func() interface{} {
		return &proto.Account{}
	}}
	emailsLock := &sync.Mutex{}
	emails := make(map[uint64]int, rep.Len())
	acc := proto.Account{}
	for id := 0; id < rep.Len(); id++ {
		acc = *rep.Get(id)
		if acc.Email.Len == 0 {
			continue
		}
		email := acc.Email.Buf[:acc.Email.Len]
		hash := murmur3.Sum64(email)
		emails[hash] = id
	}
	return &AccountsService{
		rep:          rep,
		accountsPool: accountsPool,
		emailsLock:   emailsLock,
		emails:       emails,
	}
}

func (svc *AccountsService) assignEmail(id int, email []byte) (int, bool) {
	hash := murmur3.Sum64(email)
	svc.emailsLock.Lock()
	if owner, ok := svc.emails[hash]; ok {
		svc.emailsLock.Unlock()
		return owner, false
	}
	svc.emails[hash] = id
	svc.emailsLock.Unlock()
	return id, true
}

func (svc *AccountsService) Get(id int) *proto.Account {
	return svc.rep.Get(id)
}

func (svc *AccountsService) Exists(id int) bool {
	acc := svc.rep.Get(id)
	if acc == nil || acc.Email.Len == 0 {
		return false
	}
	return true
}

func (svc *AccountsService) Create(data []byte) bool {
	src := svc.accountsPool.Get().(*proto.Account)
	if _, ok := src.UnmarshalJSON(data); !ok {
		return false
	}
	_, id, ok := parse.ParseInt(src.ID[:])
	if !ok || svc.Exists(id) {
		return false
	}
	email := src.Email.Buf[:src.Email.Len]
	if len(email) == 0 || bytes.IndexByte(email, 0x40) == -1 {
		return false
	}
	if _, ok := svc.assignEmail(id, email); !ok {
		return false
	}
	dst := *src
	dst.LikesTo = make([]proto.Like, len(src.LikesTo))
	copy(dst.LikesTo, src.LikesTo)
	svc.rep.Set(id, &dst)
	svc.accountsPool.Put(src)
	return true
}

func (svc *AccountsService) Update(id int, buf []byte) bool {
	dst := svc.rep.Get(id)
	if dst == nil || dst.Email.Len == 0 {
		return false
	}
	src := svc.accountsPool.Get().(*proto.Account)
	if _, ok := src.UnmarshalJSON(buf); !ok {
		return false
	}
	email := src.Email.Buf[:src.Email.Len]
	if len(email) > 0 {
		if bytes.IndexByte(email, 0x40) == -1 {
			return false
		}
		if _, ok := svc.assignEmail(id, email); !ok {
			return false
		}
		dst.Email = src.Email
	}
	if src.Fname > 0 {
		dst.Fname = src.Fname
	}
	if src.Sname > 0 {
		dst.Sname = src.Sname
	}
	if src.Phone[0] > 0 {
		dst.Phone = src.Phone
	}
	if src.Country > 0 {
		dst.Country = src.Country
	}
	if src.City > 0 {
		dst.City = src.City
	}
	if src.Status > 0 {
		dst.Status = src.Status
	}
	if src.PremiumStart[0] > 0 {
		dst.PremiumStart = src.PremiumStart
	}
	if src.PremiumFinish[0] > 0 {
		dst.PremiumFinish = src.PremiumFinish
	}
	if src.Interests[0] > 0 {
		dst.Interests = src.Interests
	}
	if len(src.LikesTo) > 0 {
		dst.LikesTo = append(dst.LikesTo[:0], src.LikesTo...)
	}
	svc.rep.Set(id, dst)
	svc.accountsPool.Put(src)
	return true
}

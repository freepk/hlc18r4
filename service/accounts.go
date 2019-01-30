package service

import (
	"bytes"
	"errors"
	"sync"

	"github.com/freepk/parse"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

var (
	NotFoundError   = errors.New("NotFound")
	BadRequestError = errors.New("BadRequest")
)

type AccountsService struct {
	rep          *repo.AccountsRepo
	accountsPool *sync.Pool
	likesPool    *sync.Pool
	emailsLock   *sync.Mutex
	emails       map[uint64]int
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	accountsPool := &sync.Pool{New: func() interface{} {
		return &proto.Account{}
	}}
	likesPool := &sync.Pool{New: func() interface{} {
		return &proto.Likes{}
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
		likesPool:    likesPool,
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

func (svc *AccountsService) AddLikes(data []byte) error {
	likes := svc.likesPool.Get().(*proto.Likes)
	defer svc.likesPool.Put(likes)
	if _, ok := likes.UnmarshalJSON(data); !ok {

		return BadRequestError
	}
	for _, like := range likes.Likes {
		if !svc.Exists(int(like.Liker)) || !svc.Exists(int(like.Likee)) {
			return BadRequestError
		}
	}
	return nil
}

func (svc *AccountsService) Create(data []byte) error {
	src := svc.accountsPool.Get().(*proto.Account)
	defer svc.accountsPool.Put(src)
	if _, ok := src.UnmarshalJSON(data); !ok {
		return BadRequestError
	}
	_, id, ok := parse.ParseInt(src.ID[:])
	if !ok || svc.Exists(id) {
		return BadRequestError
	}
	email := src.Email.Buf[:src.Email.Len]
	if len(email) == 0 || bytes.IndexByte(email, 0x40) == -1 {
		return BadRequestError
	}
	if _, ok := svc.assignEmail(id, email); !ok {
		return BadRequestError
	}
	dst := *src
	dst.LikesTo = make([]proto.Like, len(src.LikesTo))
	copy(dst.LikesTo, src.LikesTo)
	svc.rep.Set(id, &dst)
	return nil
}

func (svc *AccountsService) Update(id int, buf []byte) error {
	dst := svc.rep.Get(id)
	if dst == nil || dst.Email.Len == 0 {
		return NotFoundError
	}
	src := svc.accountsPool.Get().(*proto.Account)
	defer svc.accountsPool.Put(src)
	if _, ok := src.UnmarshalJSON(buf); !ok {
		return BadRequestError
	}
	email := src.Email.Buf[:src.Email.Len]
	if len(email) > 0 {
		if bytes.IndexByte(email, 0x40) == -1 {
			return BadRequestError
		}
		if _, ok := svc.assignEmail(id, email); !ok {
			return BadRequestError
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
	return nil
}

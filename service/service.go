package service

import (
	"bytes"
	"log"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsService struct {
	rep       *repo.AccountsRepo
	emails    *hashtab.HashTab
	interests *inverted.InvertedIndex
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	emails := hashtab.NewHashTab(rep.Len())
	rep.ForEach(func(id int, acc *proto.Account) {
		if acc.Email.Len > 0 {
			email := acc.Email.Buf[:acc.Email.Len]
			hash := murmur3.Sum64(email)
			emails.Set(hash, uint64(id))
		}
	})
	interests := inverted.NewInvertedIndex(rep, inverted.DefaultParts, inverted.InterestsTokens)
	return &AccountsService{rep: rep, emails: emails, interests: interests}
}

func (svc *AccountsService) Reindex() {
	log.Println("Reindex")
	total, grow := svc.interests.Rebuild()
	log.Println("Interests", total, grow)
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
	if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
		return false
	}
	tmp := *acc
	tmp.LikesTo = make([]proto.Like, len(acc.LikesTo))
	copy(tmp.LikesTo, acc.LikesTo)
	svc.rep.Add(id, &tmp)
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
		if _, ok := svc.emails.GetOrSet(hash, uint64(id)); ok {
			return false
		}
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
	return true
}

package service

import (
	"bytes"
	"sync"

	"github.com/freepk/iterator"
	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
	"gitlab.com/freepk/hlc18r4/service/indexes"
)

type AccountsService struct {
	rep           *repo.AccountsRepo
	accountsPool  *sync.Pool
	emailsLock    *sync.Mutex
	emails        map[uint64]int
	defaultIndex  *indexes.DefaultIndex
	countryIndex  *indexes.CountryIndex
	likeFromIndex *indexes.LikeIndex
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
	accountsPool := &sync.Pool{New: func() interface{} {
		return &proto.Account{}
	}}
	emailsLock := &sync.Mutex{}
	emails := make(map[uint64]int, rep.Len())
	acc := &proto.Account{}
	for id := 0; id < rep.Len(); id++ {
		*acc = *rep.Get(id)
		if acc.Email.Len > 0 {
			email := acc.Email.Buf[:acc.Email.Len]
			hash := murmur3.Sum64(email)
			emails[hash] = id
		}
	}
	defaultIndex := indexes.NewDefaultIndex(rep)
	defaultIndex.Rebuild()
	countryIndex := indexes.NewCountryIndex(rep)
	countryIndex.Rebuild()
	likeFromIndex := indexes.NewLikeIndex(rep)
	likeFromIndex.Rebuild()
	return &AccountsService{
		rep:           rep,
		accountsPool:  accountsPool,
		emailsLock:    emailsLock,
		emails:        emails,
		defaultIndex:  defaultIndex,
		countryIndex:  countryIndex,
		likeFromIndex: likeFromIndex}
}

func (svc *AccountsService) RebuildIndexes() {
	svc.defaultIndex.Rebuild()
	svc.countryIndex.Rebuild()
	svc.likeFromIndex.Rebuild()
}

func (svc *AccountsService) Exists(id int) bool {
	acc := svc.rep.Get(id)
	if acc == nil || acc.Email.Len == 0 {
		return false
	}
	return true
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

func (svc *AccountsService) BySexEq(sex []byte) iterator.Iterator {
	if token, ok := indexes.GetSexToken(sex); ok {
		return svc.defaultIndex.Sex(token)
	}
	return nil
}

func (svc *AccountsService) ByStatusEq(status []byte) iterator.Iterator {
	if token, ok := indexes.GetStatusToken(status); ok {
		return svc.defaultIndex.Status(token)
	}
	return nil
}

func (svc *AccountsService) ByStatusNeq(status []byte) iterator.Iterator {
	if token, ok := indexes.GetNotStatusToken(status); ok {
		return svc.defaultIndex.Status(token)
	}
	return nil
}

func (svc *AccountsService) ByFnameEq(fname []byte) iterator.Iterator {
	if token, ok := indexes.GetFnameToken(fname); ok {
		return svc.defaultIndex.Fname(token)
	}
	return nil
}

func (svc *AccountsService) ByFnameNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.Fname(token)
	}
	return nil
}

func (svc *AccountsService) ByFnameAny(fnames []byte) iterator.Iterator {
	var iter iterator.Iterator
	fnames, fname := parse.ScanSymbol(fnames, 0x2C)
	for len(fname) > 0 {
		if token, ok := indexes.GetFnameToken(fname); ok {
			next := svc.defaultIndex.Fname(token)
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewUnionIter(iter, next)
			}
		}
		fnames, fname = parse.ScanSymbol(fnames, 0x2C)
	}
	return iter
}

func (svc *AccountsService) BySnameEq(sname []byte) iterator.Iterator {
	if token, ok := indexes.GetSnameToken(sname); ok {
		return svc.defaultIndex.Sname(token)
	}
	return nil
}

func (svc *AccountsService) BySnameNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.Sname(token)
	}
	return nil
}

func (svc *AccountsService) ByCountryEq(country []byte) iterator.Iterator {
	if token, ok := indexes.GetCountryToken(country); ok {
		return svc.defaultIndex.Country(token)
	}
	return nil
}

func (svc *AccountsService) ByCountryNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.Country(token)
	}
	return nil
}

func (svc *AccountsService) ByCountryEqSexEq(country, sex []byte) iterator.Iterator {
	if country, ok := indexes.GetCountryToken(country); ok {
		if sex, ok := indexes.GetSexToken(sex); ok {
			return svc.countryIndex.Sex(country, sex)
		}
	}
	return nil
}

func (svc *AccountsService) ByCountryEqStatusEq(country, status []byte) iterator.Iterator {
	if country, ok := indexes.GetCountryToken(country); ok {
		if status, ok := indexes.GetStatusToken(status); ok {
			return svc.countryIndex.Status(country, status)
		}
	}
	return nil
}

func (svc *AccountsService) ByCountryEqStatusNeq(country, status []byte) iterator.Iterator {
	if country, ok := indexes.GetCountryToken(country); ok {
		if status, ok := indexes.GetNotStatusToken(status); ok {
			return svc.countryIndex.Status(country, status)
		}
	}
	return nil
}

func (svc *AccountsService) ByCityEq(city []byte) iterator.Iterator {
	if token, ok := indexes.GetCityToken(city); ok {
		return svc.defaultIndex.City(token)
	}
	return nil
}

func (svc *AccountsService) ByCityNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.City(token)
	}
	return nil
}

func (svc *AccountsService) ByCityAny(cities []byte) iterator.Iterator {
	var iter iterator.Iterator
	cities, city := parse.ScanSymbol(cities, 0x2C)
	for len(city) > 0 {
		if token, ok := indexes.GetCityToken(city); ok {
			next := svc.defaultIndex.City(token)
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewUnionIter(iter, next)
			}
		}
		cities, city = parse.ScanSymbol(cities, 0x2C)
	}
	return iter
}

func (svc *AccountsService) ByInterestsAny(interests []byte) iterator.Iterator {
	var iter iterator.Iterator
	interests, interest := parse.ScanSymbol(interests, 0x2C)
	for len(interest) > 0 {
		if token, ok := indexes.GetInterestToken(interest); ok {
			next := svc.defaultIndex.Interest(token)
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewUnionIter(iter, next)
			}
		}
		interests, interest = parse.ScanSymbol(interests, 0x2C)
	}
	return iter
}

func (svc *AccountsService) ByInterestsContains(interests []byte) iterator.Iterator {
	var iter iterator.Iterator
	interests, interest := parse.ScanSymbol(interests, 0x2C)
	for len(interest) > 0 {
		if token, ok := indexes.GetInterestToken(interest); ok {
			next := svc.defaultIndex.Interest(token)
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewInterIter(iter, next)
			}
		}
		interests, interest = parse.ScanSymbol(interests, 0x2C)
	}
	return iter
}

func (svc *AccountsService) ByBirthYear(year []byte) iterator.Iterator {
	if token, ok := indexes.GetBirthYearToken(year); ok {
		if iter := svc.defaultIndex.BirthYear(token); iter != nil {
			return iter
		}
	}
	return nil
}

func (svc *AccountsService) ByPremiumNow() iterator.Iterator {
	return svc.defaultIndex.Premium(indexes.PremiumNowToken)
}

func (svc *AccountsService) ByPremiumNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.Premium(token)
	}
	return nil
}

func (svc *AccountsService) ByPhoneNull(null []byte) iterator.Iterator {
	if token, ok := indexes.GetNullToken(null); ok {
		return svc.defaultIndex.PhoneCode(token)
	}
	return nil
}

func (svc *AccountsService) ByPhoneCode(code []byte) iterator.Iterator {
	if token, ok := indexes.GetPhoneCodeToken(code); ok {
		return svc.defaultIndex.PhoneCode(token)
	}
	return nil
}

func (svc *AccountsService) ByEmailDomain(domain []byte) iterator.Iterator {
	if token, ok := indexes.GetEmailDomainToken(domain); ok {
		return svc.defaultIndex.EmailDomain(token)
	}
	return nil
}

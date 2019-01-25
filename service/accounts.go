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
	rep          *repo.AccountsRepo
	emailsLock   *sync.Mutex
	emails       map[uint64]int
	defaultIndex *indexes.DefaultIndex
	countryIndex *indexes.CountryIndex
}

func NewAccountsService(rep *repo.AccountsRepo) *AccountsService {
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
	return &AccountsService{
		rep:          rep,
		emailsLock:   emailsLock,
		emails:       emails,
		defaultIndex: defaultIndex,
		countryIndex: countryIndex}
}

func (svc *AccountsService) RebuildIndexes() {
	svc.defaultIndex.Rebuild()
	svc.countryIndex.Rebuild()
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
	src := &proto.Account{}
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
	return true
}

func (svc *AccountsService) Update(id int, buf []byte) bool {
	dst := svc.rep.Get(id)
	if dst == nil || dst.Email.Len == 0 {
		return false
	}
	src := &proto.Account{}
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
	return true
}

func (svc *AccountsService) BySexEq(sex []byte) iterator.Iterator {
	switch string(sex) {
	case `m`:
		return svc.defaultIndex.Sex(indexes.MaleToken)
	case `f`:
		return svc.defaultIndex.Sex(indexes.FemaleToken)
	}
	return nil
}

func (svc *AccountsService) ByStatusEq(status []byte) iterator.Iterator {
	switch string(status) {
	case `свободны`:
		return svc.defaultIndex.Status(indexes.SingleToken)
	case `заняты`:
		return svc.defaultIndex.Status(indexes.InRelToken)
	case `всё сложно`:
		return svc.defaultIndex.Status(indexes.ComplToken)
	}
	return nil
}

func (svc *AccountsService) ByStatusNeq(status []byte) iterator.Iterator {
	switch string(status) {
	case `свободны`:
		return svc.defaultIndex.Status(indexes.NotSingleToken)
	case `заняты`:
		return svc.defaultIndex.Status(indexes.NotInRelToken)
	case `всё сложно`:
		return svc.defaultIndex.Status(indexes.NotComplToken)
	}
	return nil
}

func (svc *AccountsService) ByCountryEq(country []byte) iterator.Iterator {
	if token, ok := proto.CountryToken(country); ok {
		return svc.defaultIndex.Country(token)
	}
	return nil
}

func (svc *AccountsService) ByCountryNull(null []byte) iterator.Iterator {
	switch string(null) {
	case `0`:
		return svc.defaultIndex.Country(indexes.NotNullToken)
	case `1`:
		return svc.defaultIndex.Country(indexes.NullToken)
	}
	return nil
}

func (svc *AccountsService) ByCityEq(city []byte) iterator.Iterator {
	if token, ok := proto.CityToken(city); ok {
		return svc.defaultIndex.City(token)
	}
	return nil
}

func (svc *AccountsService) ByCityNull(null []byte) iterator.Iterator {
	switch string(null) {
	case `0`:
		return svc.defaultIndex.City(indexes.NotNullToken)
	case `1`:
		return svc.defaultIndex.City(indexes.NullToken)
	}
	return nil
}

func (svc *AccountsService) ByCityAny(cities []byte) iterator.Iterator {
	var iter iterator.Iterator
	cities, city := parse.ScanSymbol(cities, 0x2C)
	for len(city) > 0 {
		if token, ok := proto.CityToken(city); ok {
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
		if token, ok := proto.InterestToken(interest); ok {
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
		if token, ok := proto.InterestToken(interest); ok {
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
	if _, year, ok := parse.ParseInt(year); ok && year > 1975 {
		if iter := svc.defaultIndex.BirthYear(year - 1975); iter != nil {
			return iter
		}
	}
	return nil
}

func (svc *AccountsService) ByFnameEq(fname []byte) iterator.Iterator {
	if token, ok := proto.FnameToken(fname); ok {
		return svc.defaultIndex.Fname(token)
	}
	return nil
}

func (svc *AccountsService) ByFnameNull(null []byte) iterator.Iterator {
	switch string(null) {
	case `0`:
		return svc.defaultIndex.Fname(indexes.NotNullToken)
	case `1`:
		return svc.defaultIndex.Fname(indexes.NullToken)
	}
	return nil
}

func (svc *AccountsService) ByFnameAny(fnames []byte) iterator.Iterator {
	var iter iterator.Iterator
	fnames, fname := parse.ScanSymbol(fnames, 0x2C)
	for len(fname) > 0 {
		if token, ok := proto.FnameToken(fname); ok {
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
	if token, ok := proto.SnameToken(sname); ok {
		return svc.defaultIndex.Sname(token)
	}
	return nil
}

func (svc *AccountsService) BySnameNull(null []byte) iterator.Iterator {
	switch string(null) {
	case `0`:
		return svc.defaultIndex.Sname(indexes.NotNullToken)
	case `1`:
		return svc.defaultIndex.Sname(indexes.NullToken)
	}
	return nil
}

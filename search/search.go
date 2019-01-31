package search

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	CommonPartition = 0
)

const (
	SexField = iota
	StatusField
	CountryField
	CityField
	InterestField
	FnameField
	SnameField
	BirthYearField
	PremiumField
	PhoneCodeField
	EmailDomainField
)

type SearchService struct {
	rep    *repo.AccountsRepo
	likes  *LikeIndex
	common *CommonIndex
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	likes := NewLikeIndex(rep)
	common := NewCommonIndex(rep)
	return &SearchService{rep: rep, likes: likes, common: common}
}

func (svc *SearchService) Rebuild() {
	svc.likes.Rebuild()
}

func (svc *SearchService) Likes(id int) iterator.Iterator {
	if svc.rep.IsSet(id) {
		return svc.likes.Iterator(id)
	}
	return nil
}

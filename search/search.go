package search

import (
	"github.com/freepk/inverted"
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
	rep   *repo.AccountsRepo
	likes *LikeIndex
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	likes := NewLikeIndex(rep)
	return &SearchService{rep: rep, likes: likes}
}

func (svc *SearchService) Rebuild() {
	svc.likes.Rebuild()
}

func (svc *SearchService) Likes(t int) *LikeIter {
	return svc.likes.Likes(t)
}

func (svc *SearchService) Sex(t int) *inverted.TokenIter {
	return nil
}

func (svc *SearchService) Status(t int) *inverted.TokenIter {
	return nil
}

func (svc *SearchService) Country(t int) *inverted.TokenIter {
	return nil
}

func (svc *SearchService) Countries(country int) *CountryIndex {
	return nil
}

type CountryIndex struct {
}

func (idx *CountryIndex) Sex(t int) *inverted.TokenIter {
	return nil
}

func (idx *CountryIndex) Status(t int) *inverted.TokenIter {
	return nil
}

func (idx *CountryIndex) Interest(t int) *inverted.TokenIter {
	return nil
}

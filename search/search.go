package search

import (
	"sync"

	"github.com/freepk/inverted"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	CommonPart = 0
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
	rep       *repo.AccountsRepo
	likes     *LikeIndex
	countries *inverted.Inverted
}

func newCountryIndex(rep *repo.AccountsRepo) *inverted.Inverted {
	proc := NewAccountsProcessor(rep, countryProc, 80, 5)
	return inverted.NewInverted(proc)
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	likes := NewLikeIndex(rep)
	countries := newCountryIndex(rep)
	return &SearchService{rep: rep, likes: likes, countries: countries}
}

func (svc *SearchService) Rebuild() {
	gr := &sync.WaitGroup{}
	gr.Add(2)
	go func() { defer gr.Done(); svc.likes.Rebuild() }()
	go func() { defer gr.Done(); svc.countries.Rebuild() }()
	gr.Wait()
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

func (svc *SearchService) Countries(t int) *CountryIndex {
	return &CountryIndex{part: svc.countries.Part(t)}
}

type CountryIndex struct {
	part *inverted.Part
}

func (idx *CountryIndex) Sex(t int) *inverted.TokenIter {
	return idx.part.Field(SexField).Token(t).Iter()
}

func (idx *CountryIndex) Status(t int) *inverted.TokenIter {
	return idx.part.Field(StatusField).Token(t).Iter()
}

func (idx *CountryIndex) Interest(t int) *inverted.TokenIter {
	return idx.part.Field(InterestField).Token(t).Iter()
}

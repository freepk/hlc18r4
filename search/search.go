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
	cities    *inverted.Inverted
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	likes := NewLikeIndex(rep)
	countries := inverted.NewInverted(NewAccountsProcessor(rep, countryProc, 80, 5))
	cities := inverted.NewInverted(NewAccountsProcessor(rep, cityProc, 800, 5))
	return &SearchService{rep: rep, likes: likes, countries: countries, cities: cities}
}

func (svc *SearchService) Rebuild() {
	gr := &sync.WaitGroup{}
	gr.Add(2)
	go func() {
		defer gr.Done()
		svc.likes.Rebuild()
	}()
	go func() {
		defer gr.Done()
		svc.countries.Rebuild()
		svc.cities.Rebuild()
	}()
	gr.Wait()
}

func (svc *SearchService) Likes(t int) *LikeIter {
	return svc.likes.Likes(t)
}

// Common

//func (svc *SearchService) Sex(t int) *inverted.TokenIter {
//	return nil
//}

//func (svc *SearchService) Status(t int) *inverted.TokenIter {
//	return nil
//}

//func (svc *SearchService) Country(t int) *inverted.TokenIter {
//	return nil
//}

// Countries

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

// Cities

func (svc *SearchService) Cities(t int) *CityIndex {
	return &CityIndex{part: svc.cities.Part(t)}
}

type CityIndex struct {
	part *inverted.Part
}

func (idx *CityIndex) Sex(t int) *inverted.TokenIter {
	return idx.part.Field(SexField).Token(t).Iter()
}

func (idx *CityIndex) Status(t int) *inverted.TokenIter {
	return idx.part.Field(StatusField).Token(t).Iter()
}

func (idx *CityIndex) Interest(t int) *inverted.TokenIter {
	return idx.part.Field(InterestField).Token(t).Iter()
}

package service

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type IndexEnum int

const (
	SexIndex = IndexEnum(iota)
	StatusIndex
	InterestsIndex
)

type PartEnum int

const (
	CommonPart = PartEnum(iota)
	MalePart
	FemalePart
	FreePart
	BusyPart
	ComplicatedPart
	NotFreePart
	NotBusyPart
	NotComplicatedPart
	MaleFreePart
	MaleBusyPart
	MaleComplicatedPart
	MaleNotFreePart
	MaleNotBusyPart
	MaleNotComplicatedPart
	FemaleFreePart
	FemaleBusyPart
	FemaleComplicatedPart
	FemaleNotFreePart
	FemaleNotBusyPart
	FemaleNotComplicatedPart
	CountryPart
)

type FiltersService struct {
	rep     *repo.AccountsRepo
	indexes map[IndexEnum]*inverted.InvertedIndex
}

func NewFiltersService(rep *repo.AccountsRepo) *FiltersService {
	indexes := make(map[IndexEnum]*inverted.InvertedIndex, 32)
	return &FiltersService{rep: rep, indexes: indexes}
}

func (svc *FiltersService) RebuildIndexes() {
}

func (svc *FiltersService) InterestsAny(part PartEnum, buf []byte) iterator.Iterator {
	index, _ := svc.indexes[InterestsIndex]
	_ = index
	return nil
}

func (svc *FiltersService) InterestsContains(part PartEnum, buf []byte) iterator.Iterator {
	return nil
}

func interestsTokenizer(name []byte) inverted.Token {
	token, ok := proto.InterestDict.Id(name)
	if !ok {
		return inverted.Token(token)
	}
	return inverted.NullToken
}

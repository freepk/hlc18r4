package service

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type FilterIndexEnum int

const (
	SexIndex = FilterIndexEnum(iota)
	StatusIndex
	InterestsIndex
)

type FilterPartEnum int

const (
	CommonPart = FilterPartEnum(iota)
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
	indexes map[FilterIndexEnum]*inverted.InvertedIndex
}

func NewFiltersService(rep *repo.AccountsRepo) *FiltersService {
	indexes := make(map[FilterIndexEnum]*inverted.InvertedIndex, 32)
	return &FiltersService{rep: rep, indexes: indexes}
}

func (svc *FiltersService) RebuildIndexes() {
}

func (svc *FiltersService) SexEq(part FilterPartEnum, sex proto.SexEnum) iterator.Iterator {
	return nil
}

func (svc *FiltersService) StatusEq(part FilterPartEnum, status proto.StatusEnum) iterator.Iterator {
	return nil
}

func (svc *FiltersService) StatusNeq(part FilterPartEnum, status proto.StatusEnum) iterator.Iterator {
	return nil
}

func (svc *FiltersService) InterestsAny(part FilterPartEnum, buf []byte) iterator.Iterator {
	return nil
}

func (svc *FiltersService) InterestsContains(part FilterPartEnum, buf []byte) iterator.Iterator {
	return nil
}

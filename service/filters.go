package service

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/repo"
	"gitlab.com/freepk/hlc18r4/proto"
)

const (
	CommonPart = inverted.Partition(0)
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
	indexes []*inverted.InvertedIndex
}

func NewFiltersService(rep *repo.AccountsRepo) *FiltersService {
	indexes := []*inverted.InvertedIndex{
		inverted.NewInvertedIndex(rep, DefaultParts, InterestsTokens)}
	return &FiltersService{rep: rep, indexes: indexes}
}

func (svc *FiltersService) RebuildIndexes() {
	svc.indexes[0].Rebuild()
}

func (svc *FiltersService) InterestsAny(part inverted.Partition, buf []byte) iterator.Iterator {
	index := svc.indexes[0]
	buf, interest := parse.ScanSymbol(buf, 0x2C)
	for len(interest) > 0 {
		println(interest)
		buf, interest = parse.ScanSymbol(buf, 0x2C)
	}
	return iterator.NewUnionIter(index.Iterator(part, 10), index.Iterator(part, 15))
}

func (svc *FiltersService) InterestsContains(part inverted.Partition, buf []byte) iterator.Iterator {
	index := svc.indexes[0]
	return iterator.NewInterIter(index.Iterator(part, 10), index.Iterator(part, 15))
}

func InterestsTokens(acc *proto.Account, tokens []inverted.Token) []inverted.Token {
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, inverted.Token(interest))
	}
	return tokens
}

func DefaultParts(acc *proto.Account, parts []inverted.Partition) []inverted.Partition {
	parts = append(parts, CommonPart)
	return parts
}

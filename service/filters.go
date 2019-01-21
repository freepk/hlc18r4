package service

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	sexIndex = iota
	statusIndex
	interestsIndex
)

type PartitionEnum inverted.Partition

const (
	CommonPart = PartitionEnum(iota)
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
	indexes map[int]*inverted.InvertedIndex
}

func NewFiltersService(rep *repo.AccountsRepo) *FiltersService {
	indexes := make(map[int]*inverted.InvertedIndex, 32)
	indexes[InterestsIndex] = inverted.NewInvertedIndex()
	return &FiltersService{rep: rep, indexes: indexes}
}

func (svc *FiltersService) RebuildIndexes() {
	for _, index := range svc.indexes {
		index.Rebuild()
	}
}

func (svc *FiltersService) InterestsAny(part PartitionEnum, buf []byte) (it iterator.Iterator) {
	var name []byte
	var token inverted.Token
	index, _ := svc.indexes[InterestsIndex]
	for {
		if buf, name = parse.ScanSymbol(buf, 0x2C); len(name) == 0 {
			break
		}
		if token = interestToken(name); token == inverted.NullToken {
			break
		}
		it = iterator.NewUnionIter(it, index.TokenIterator(inverted.Partition(part), token))
	}
	return
}

func interestToken(interest []byte) inverted.Token {
	token, ok := proto.InterestDict.Id(interest)
	if !ok {
		return inverted.Token(token)
	}
	return inverted.NullToken
}

/*
func interestsTokens([]inverted.Token tokens, acc *proto.Account) []inverted.Token {
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, inverted.Token(interest))
	}
	return tokens
}

func defaultParts(parts []inverted.Partition, acc *proto.Account) []inverted.Partition {
	parts = append(parts, inverted.Part(CommonPart))
	return parts
}
*/

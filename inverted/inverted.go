package inverted

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	partsPerIndex = 256
	tokensPerPart = 2048
)

type PartitionEnum int

const (
	CommonPart = PartitionEnum(0)
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

type PartsFunc func(*proto.Account, []PartitionEnum) []PartitionEnum

type Token int

type TokensFunc func(*proto.Account, []Token) []Token

type InvertedIndex struct {
	rep        *repo.AccountsRepo
	tokens     [][][]uint32
	partsFunc  PartsFunc
	tokensFunc TokensFunc
}

func NewInvertedIndex(rep *repo.AccountsRepo, partsFunc PartsFunc, tokensFunc TokensFunc) *InvertedIndex {
	tokens := make([][][]uint32, partsPerIndex)
	for i := 0; i < partsPerIndex; i++ {
		tokens[i] = make([][]uint32, tokensPerPart)
	}
	return &InvertedIndex{rep: rep, tokens: tokens, partsFunc: partsFunc, tokensFunc: tokensFunc}
}

func (ii *InvertedIndex) Iterator(part PartitionEnum, token int) *ReverseIter {
	return NewReverseIter(ii.tokens[part][token])
}

func (ii *InvertedIndex) Rebuild() (int, int) {
	parts := make([]PartitionEnum, 0, partsPerIndex)
	tokens := make([]Token, 0, tokensPerPart)
	want := make([][]int, len(ii.tokens))
	for i := range ii.tokens {
		want[i] = make([]int, len(ii.tokens[i]))
	}
	total := 0
	ii.rep.ForEach(func(id int, acc *proto.Account) {
		parts = ii.partsFunc(acc, parts[:0])
		tokens = ii.tokensFunc(acc, tokens[:0])
		for _, part := range parts {
			for _, token := range tokens {
				total++
				want[part][token]++
			}
		}
	})
	grow := 0
	for part, tokens := range ii.tokens {
		for token, ids := range tokens {
			if want[part][token] > cap(ids) {
				grow += want[part][token]
			}
		}
	}
	buffer := make([]uint32, grow)
	for part, tokens := range ii.tokens {
		for token, ids := range tokens {
			if want[part][token] > cap(ids) {
				grow := want[part][token]
				ii.tokens[part][token], buffer = buffer[:grow], buffer[grow:]
			}
			ii.tokens[part][token] = ii.tokens[part][token][:0]
		}
	}
	ii.rep.ForEach(func(id int, acc *proto.Account) {
		parts = ii.partsFunc(acc, parts[:0])
		tokens = ii.tokensFunc(acc, tokens[:0])
		for _, part := range parts {
			for _, token := range tokens {
				reverseID := reverseMaxID - uint32(id)
				ii.tokens[part][token] = append(ii.tokens[part][token], reverseID)
			}
		}
	})
	return total, grow
}

func InterestsTokens(acc *proto.Account, tokens []Token) []Token {
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, Token(interest))
	}
	return tokens
}

func DefaultParts(acc *proto.Account, parts []PartitionEnum) []PartitionEnum {
	// Common
	parts = append(parts, CommonPart)
	// Sex
	/*
		switch acc.Sex {
		case proto.MaleSex:
			parts = append(parts, MalePart)
		case proto.FemaleSex:
			parts = append(parts, FemalePart)
		}
		// Status
		switch acc.Status {
		case proto.FreeStatus:
			parts = append(parts, FreePart, NotBusyPart, NotComplicatedPart)
		case proto.BusyStatus:
			parts = append(parts, BusyPart, NotFreePart, NotComplicatedPart)
		case proto.ComplicatedStatus:
			parts = append(parts, ComplicatedPart, NotFreePart, NotBusyPart)
		}
		// Sex & Status
		switch {
		case acc.Sex == proto.MaleSex && acc.Status == proto.FreeStatus:
			parts = append(parts, MaleFreePart, MaleNotBusyPart, MaleNotComplicatedPart)
		case acc.Sex == proto.MaleSex && acc.Status == proto.BusyStatus:
			parts = append(parts, MaleBusyPart, MaleNotFreePart, MaleNotComplicatedPart)
		case acc.Sex == proto.MaleSex && acc.Status == proto.ComplicatedStatus:
			parts = append(parts, MaleComplicatedPart, MaleNotFreePart, MaleNotBusyPart)
		case acc.Sex == proto.FemaleSex && acc.Status == proto.FreeStatus:
			parts = append(parts, FemaleFreePart, FemaleNotBusyPart, FemaleNotComplicatedPart)
		case acc.Sex == proto.FemaleSex && acc.Status == proto.BusyStatus:
			parts = append(parts, FemaleBusyPart, FemaleNotFreePart, FemaleNotComplicatedPart)
		case acc.Sex == proto.FemaleSex && acc.Status == proto.ComplicatedStatus:
			parts = append(parts, FemaleComplicatedPart, FemaleNotFreePart, FemaleNotBusyPart)
		}
		// Country part
		country := acc.Country + CountryPart
		parts = append(parts, uint8(country))
	*/
	return parts
}

package inverted

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	partsPerIndex = 256
	tokensPerPart = 2048
)

const (
	CommonPart = 0
	MalePart
	FemalePart
	FreePart = 20
	BusyPart
	ComplicatedPart
	NotFreePart = 30
	NotBusyPart
	NotComplicatedPart
	MaleFreePart = 40
	MaleBusyPart
	MaleComplicatedPart
	MaleNotFreePart = 50
	MaleNotBusyPart
	MaleNotComplicatedPart
	FemaleFreePart = 60
	FemaleBusyPart
	FemaleComplicatedPart
	FemaleNotFreePart = 70
	FemaleNotBusyPart
	FemaleNotComplicatedPart
	CountryPart = 100
)

type InvertedIndex struct {
	rep           *repo.AccountsRepo
	tokens        [][][]uint32 // [part][token][]int
	partsHandler  HandlerFunc
	tokensHandler HandlerFunc
}

type HandlerFunc func(*proto.Account, []int) []int

func NewInvertedIndex(rep *repo.AccountsRepo, partsHandler, tokensHandler HandlerFunc) *InvertedIndex {
	tokens := make([][][]uint32, partsPerIndex)
	for i := 0; i < partsPerIndex; i++ {
		tokens[i] = make([][]uint32, tokensPerPart)
	}
	return &InvertedIndex{rep: rep, tokens: tokens, partsHandler: partsHandler, tokensHandler: tokensHandler}
}

func (ii *InvertedIndex) Rebuild() int {
	parts := make([]int, 0, partsPerIndex)
	tokens := make([]int, 0, tokensPerPart)
	want := make([][]int, len(ii.tokens))
	for i := range ii.tokens {
		want[i] = make([]int, len(ii.tokens[i]))
	}
	total := 0
	ii.rep.ForEach(func(id int, acc *proto.Account) {
		parts = ii.partsHandler(acc, parts)
		tokens = ii.tokensHandler(acc, tokens)
		for _, part := range parts {
			for _, token := range tokens {
				total++
				want[part][token]++
			}
		}
	})
	return total
}

func InterestsTokens(acc *proto.Account, tokens []int) []int {
	tokens = tokens[:0]
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, int(interest))
	}
	return tokens
}

func DefaultParts(acc *proto.Account, parts []int) []int {
	parts = parts[:0]
	// Common
	parts = append(parts, CommonPart)
	// Sex
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
	country := int(acc.Country) + CountryPart
	parts = append(parts, country)
	return parts
}

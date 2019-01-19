package inverted

import (
	"fmt"

	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	partsPerIndex = 256
	tokensPerPart = 2048
)

const (
	CommonPart               = 0
	MalePart                 = 10
	FemalePart               = 11
	FreePart                 = 20
	BusyPart                 = 21
	ComplicatedPart          = 22
	NotFreePart              = 30
	NotBusyPart              = 31
	NotComplicatedPart       = 32
	MaleFreePart             = 40
	MaleBusyPart             = 41
	MaleComplicatedPart      = 42
	MaleNotFreePart          = 50
	MaleNotBusyPart          = 51
	MaleNotComplicatedPart   = 52
	FemaleFreePart           = 60
	FemaleBusyPart           = 61
	FemaleComplicatedPart    = 62
	FemaleNotFreePart        = 70
	FemaleNotBusyPart        = 71
	FemaleNotComplicatedPart = 72
	CountryPart              = 100
)

type PartsHandlerFunc func(*proto.Account, []uint8) []uint8

type TokensHandlerFunc func(*proto.Account, []uint16) []uint16

type InvertedIndex struct {
	rep           *repo.AccountsRepo
	tokens        [][][]uint32
	partsHandler  PartsHandlerFunc
	tokensHandler TokensHandlerFunc
}

func NewInvertedIndex(rep *repo.AccountsRepo, partsHandler PartsHandlerFunc, tokensHandler TokensHandlerFunc) *InvertedIndex {
	tokens := make([][][]uint32, partsPerIndex)
	for i := 0; i < partsPerIndex; i++ {
		tokens[i] = make([][]uint32, tokensPerPart)
	}
	return &InvertedIndex{rep: rep, tokens: tokens, partsHandler: partsHandler, tokensHandler: tokensHandler}
}

func (ii *InvertedIndex) Rebuild() int {
	parts := make([]uint8, 0, partsPerIndex)
	tokens := make([]uint16, 0, tokensPerPart)
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
	grow := 0
	for part, tokens := range ii.tokens {
		for token, ids := range tokens {
			if want[part][token] > cap(ids) {
				grow += want[part][token]
			}
		}
	}
	buffer := make([]uint32, grow)
	fmt.Println("grow", grow)
	for part, tokens := range ii.tokens {
		for token, ids := range tokens {
			// grow if needed
			if want[part][token] > cap(ids) {
				grow := want[part][token]
				ii.tokens[part][token], buffer = buffer[:grow], buffer[grow:]
			}
			// reset
			ii.tokens[part][token] = ii.tokens[part][token][:0]
		}
	}
	return total
}

func InterestsTokens(acc *proto.Account, tokens []uint16) []uint16 {
	tokens = tokens[:0]
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, uint16(interest))
	}
	return tokens
}

func DefaultParts(acc *proto.Account, parts []uint8) []uint8 {
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
	country := acc.Country + CountryPart
	parts = append(parts, uint8(country))
	return parts
}

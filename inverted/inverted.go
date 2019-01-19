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

type PartsFunc func(*proto.Account, []uint8) []uint8

type TokensFunc func(*proto.Account, []uint16) []uint16

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

func (ii *InvertedIndex) Rebuild() (int, int) {
	parts := make([]uint8, 0, partsPerIndex)
	tokens := make([]uint16, 0, tokensPerPart)
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
			// grow if needed
			if want[part][token] > cap(ids) {
				grow := want[part][token]
				ii.tokens[part][token], buffer = buffer[:grow], buffer[grow:]
			}
			// reset
			ii.tokens[part][token] = ii.tokens[part][token][:0]
		}
	}
	ii.rep.ForEach(func(id int, acc *proto.Account) {
		parts = ii.partsFunc(acc, parts[:0])
		tokens = ii.tokensFunc(acc, tokens[:0])
		for _, part := range parts {
			for _, token := range tokens {
				ii.tokens[part][token] = append(ii.tokens[part][token], uint32(id))
			}
		}
	})
	return total, grow
}

// sex_eq - single?
// status_eq - single?
// status_neq - single?

// email_domain
// email_lt
// email_gt
func EmailTokens(acc *proto.Account, tokens []uint16) []uint16 {
	return tokens
}

// +fname_eq
// +fname_neq
func FnameTokens(acc *proto.Account, tokens []uint16) []uint16 {
	tokens = append(tokens, uint16(acc.Fname))
	return tokens
}

// +sname_eq
// sname_starts
// +sname_null
func SnameTokens(acc *proto.Account, tokens []uint16) []uint16 {
	tokens = append(tokens, uint16(acc.Sname))
	return tokens
}

// phone_code
// phone_null
func PhoneTokens(acc *proto.Account, tokens []uint16) []uint16 {
	return tokens
}

// +country_eq
// +country_null
func CountryTokens(acc *proto.Account, tokens []uint16) []uint16 {
	tokens = append(tokens, uint16(acc.Country))
	return tokens
}

// city_eq
// city_any
// city_null
func CityTokens(acc *proto.Account, tokens []uint16) []uint16 {
	tokens = append(tokens, uint16(acc.City))
	return tokens
}

// birth_lt
// birth_gt
// birth_year

// interests_contains
// interests_any
func InterestsTokens(acc *proto.Account, tokens []uint16) []uint16 {
	for _, interest := range acc.Interests {
		if interest == 0 {
			break
		}
		tokens = append(tokens, uint16(interest))
	}
	return tokens
}

// likes_contains

// premium_now
// premium_null
func PremiumTokens(acc *proto.Account, tokens []uint16) []uint16 {
	return tokens
}

func DefaultParts(acc *proto.Account, parts []uint8) []uint8 {
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

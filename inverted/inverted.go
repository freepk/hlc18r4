package inverted

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	partsPerIndex = 256
	tokensPerPart = 2048
)

var (
	IndexedAccountsTotal = 0
)

type InvertedIndex struct {
	rep      *repo.AccountsRepo
	handler  TokenFunc
	tokens   [][][]uint32 // [part][token][]int
	counters [][][]uint32 // [part][token][]int
}

type TokenFunc func(acc *proto.Account, inParts, inTokens []int) (id uint32, outParts []int, outTokens []int)

func NewInvertedIndex(rep *repo.AccountsRepo, handler TokenFunc) *InvertedIndex {
	tokens := make([][][]uint32, partsPerIndex)
	for i := 0; i < partsPerIndex; i++ {
		tokens[i] = make([][]uint32, tokensPerPart)
	}
	return &InvertedIndex{rep: rep, handler: handler, tokens: tokens}
}

func (ii *InvertedIndex) Rebuild() {
	/*
		id := uint32(0)
		parts := make([]int, 0, partsPerIndex)
		tokens := make([]int, 0, tokensPerPart)
		counters := make([][]int, len(ii.tokens))
		for i := range counters {
			counters[i] = make([]int, len(ii.tokens[i]))
		}
		ii.rep.ForEach(func(acc *proto.Account) {
			_, parts, tokens = ii.handler(acc, parts, tokens)
			for _, part := range parts {
				for _, token := range tokens {
					IndexedAccountsTotal++
					counters[part][token]++
				}
			}
		})
		for part, counter := range counters {
			for token, count := range counter {
				if count > 0 {
					capacity := cap(ii.tokens[part][token])
					if capacity < count {
						ii.tokens[part][token] = make([]uint32, 0, count)
					} else {
						ii.tokens[part][token] = ii.tokens[part][token][:0]
					}
				}
			}
		}
		ii.rep.ForEach(func(acc *proto.Account) {
			id, parts, tokens = ii.handler(acc, parts, tokens)
			for _, part := range parts {
				for _, token := range tokens {
					ii.tokens[part][token] = append(ii.tokens[part][token], id)
				}
			}
		})
	*/
}

func InterestToken(acc *proto.Account, inParts, inTokens []int) (id uint32, outParts, outTokens []int) {
	/*
		id = acc.ID
		outParts = inParts[:0]
		outTokens = inTokens[:0]

		for _, interest := range acc.Interests {
			outTokens = append(outTokens, int(interest))
		}
		outParts = append(outParts, 0)
		// Male = 10
		// Female = 11
		switch acc.Sex {
		case proto.MaleSex:
			outParts = append(outParts, 10)
		case proto.FemaleSex:
			outParts = append(outParts, 11)
		}
		// Free = 20, NotFree = 30
		// Busy = 21, NotBusy = 31
		// Compl = 22, NotCompl = 32
		switch acc.Status {
		case proto.FreeStatus:
			outParts = append(outParts, 20, 31, 32)
		case proto.BusyStatus:
			outParts = append(outParts, 21, 30, 32)
		case proto.ComplicatedStatus:
			outParts = append(outParts, 22, 30, 31)
		}
		country := int(acc.Country) + 100
		outParts = append(outParts, country)
	*/
	return
}

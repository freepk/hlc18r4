package inverted

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	partsPerIndex = 256
	tokensPerPart = 2048
)

type Partition int

type PartsFunc func([]Partition, *proto.Account) []Partition

type Token int

const (
	NullToken    = Token(0)
	NotNullToken = Token(1)
)

type TokensFunc func([]Token, *proto.Account) []Token

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

func (ii *InvertedIndex) TokenIterator(part Partition, token int) *TokenIter {
	return NewTokenIter(ii.tokens[part][token])
}

func (ii *InvertedIndex) Rebuild() (int, int) {
	parts := make([]Partition, 0, partsPerIndex)
	tokens := make([]Token, 0, tokensPerPart)
	want := make([][]int, len(ii.tokens))
	for i := range ii.tokens {
		want[i] = make([]int, len(ii.tokens[i]))
	}
	total := 0
	ii.rep.Forward(func(id int, acc *proto.Account) {
		parts = ii.partsFunc(parts[:0], acc)
		tokens = ii.tokensFunc(tokens[:0], acc)
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
	ii.rep.Reverse(func(id int, acc *proto.Account) {
		parts = ii.partsFunc(parts[:0], acc)
		tokens = ii.tokensFunc(tokens[:0], acc)
		for _, part := range parts {
			for _, token := range tokens {
				id := 2000000 - uint32(id)
				ii.tokens[part][token] = append(ii.tokens[part][token], id)
			}
		}
	})
	return total, grow
}

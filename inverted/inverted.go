package inverted

import (
	"fmt"
)

type Document struct {
	ID     int
	Parts  []int
	Tokens [][]int
}

type Indexer interface {
	Reset()
	Next() (*Document, bool)
}

type Inverted struct {
	indexer Indexer
}

func NewInverted(indexer Indexer) *Inverted {
	return &Inverted{indexer: indexer}
}

func (inv *Inverted) prepare() (layout [][][]uint32) {
	it := inv.indexer
	it.Reset()
	for {
		doc, ok := it.Next()
		if !ok {
			break
		}
		for _, part := range doc.Parts {
			if grow := part + 1 - len(layout); grow > 0 {
				for i := 0; i < grow; i++ {
					layout = append(layout, make([][]uint32, len(doc.Tokens)))
				}
			}
			for index, tokens := range doc.Tokens {
				for _, token := range tokens {
					if grow := token + 1 - len(layout[part][index]); grow > 0 {
						layout[part][index] = append(layout[part][index], make([]uint32, grow)...)
					}
					layout[part][index][token]++
				}
			}
		}
	}
	return
}

func (inv *Inverted) Rebuild() {
	fmt.Println(inv.prepare())
}

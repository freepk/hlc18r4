package index

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

func (inv *Inverted) Rebuild() {
	inv.indexer.Reset()
	doc, ok := inv.indexer.Next()
	for ok {
		_ = doc
		doc, ok = inv.indexer.Next()
	}
}

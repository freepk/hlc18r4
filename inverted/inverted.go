package inverted

type Document struct {
	ID         int
	Partitions []int
	Indexes    [][]int
}

type Indexer interface {
	Reset()
	Next() (*Document, bool)
}

type document uint32

type vector []document

type index []vector

type partition []index

type Inverted struct {
	indexer    Indexer
	partitions []partition
}

func NewInverted(indexer Indexer) *Inverted {
	return &Inverted{indexer: indexer}
}

func (inv *Inverted) Rebuild() {
	inv.indexer.Reset()
	if doc, ok := inv.indexer.Next(); ok {
		_ = doc
	}
}

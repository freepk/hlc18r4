package index

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

type token []document

type index []token

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
		parts := len(doc.Partitions)
		if parts > len(inv.partitions) {
			inv.partitions = make([]partition, parts)
			indexes := len(doc.Indexes)
			for part := range inv.partitions {
				inv.partitions[part] = make([]index, indexes)
			}
		}
	}
}

package index

type Item struct {
	ID         int
	Partitions []int
	Tokens     [][]int
}

type Indexer interface {
	Reset()
	Next() (*Item, bool)
}

type Index struct {
	indexer Indexer
	tokens  [][][][]int
}

func NewIndex(indexer Indexer) *Index {
	return &Index{indexer: indexer}
}

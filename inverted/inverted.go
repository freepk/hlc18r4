package inverted

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
	layout  [][][][]uint32
}

func NewInverted(indexer Indexer) *Inverted {
	return &Inverted{indexer: indexer}
}

func (inv *Inverted) prepare() (layout [][][]int) {
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
					layout = append(layout, make([][]int, len(doc.Tokens)))
				}
			}
			for field, tokens := range doc.Tokens {
				for _, token := range tokens {
					if grow := token + 1 - len(layout[part][field]); grow > 0 {
						layout[part][field] = append(layout[part][field], make([]int, grow)...)
					}
					layout[part][field][token]++
				}
			}
		}
	}
	return
}

func (inv *Inverted) Rebuild() {
	layout := inv.prepare()
	if grow := len(layout) - len(inv.layout); grow > 0 {
		inv.layout = append(inv.layout, make([][][][]uint32, grow)...)
	}
	for part := range layout {
		if grow := len(layout[part]) - len(inv.layout[part]); grow > 0 {
			inv.layout[part] = append(inv.layout[part], make([][][]uint32, grow)...)
		}
		for field := range layout[part] {
			if grow := len(layout[part][field]) - len(inv.layout[part][field]); grow > 0 {
				inv.layout[part][field] = append(inv.layout[part][field], make([][]uint32, grow)...)
			}
			for token := range layout[part][field] {
				if grow := layout[part][field][token] - len(inv.layout[part][field][token]); grow > 0 {
					inv.layout[part][field][token] = append(inv.layout[part][field][token], make([]uint32, grow)...)
				}
			}
		}
	}
}

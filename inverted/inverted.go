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

func (inv *Inverted) Rebuild() {
	var layout [][][]int
	it := inv.indexer
	it.Reset()
	for {
		doc, ok := it.Next()
		if !ok {
			break
		}
		for _, part := range doc.Parts {
			if grow := part + 1 - len(layout); grow > 0 {
				layout = append(layout, make([][][]int, grow)...)
			}
			for field := range doc.Tokens {
				if grow := len(doc.Tokens) - len(layout[part]); grow > 0 {
					layout[part] = append(layout[part], make([][]int, grow)...)
				}
				for _, token := range doc.Tokens[field] {
					if grow := token + 1 - len(layout[part][field]); grow > 0 {
						layout[part][field] = append(layout[part][field], make([]int, grow)...)
					}
					layout[part][field][token]++
				}
			}
		}
	}
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
				if grow := layout[part][field][token] - cap(inv.layout[part][field][token]); grow > 0 {
					inv.layout[part][field][token] = append(inv.layout[part][field][token], make([]uint32, grow*105/100)...)
				}
				inv.layout[part][field][token] = inv.layout[part][field][token][:0]
			}
		}
	}
	it.Reset()
	for {
		doc, ok := it.Next()
		if !ok {
			break
		}
		for _, part := range doc.Parts {
			for field := range doc.Tokens {
				for _, token := range doc.Tokens[field] {
					inv.layout[part][field][token] = append(inv.layout[part][field][token], uint32(doc.ID))
				}
			}
		}
	}
}

func (inv *Inverted) Iterator(part, field, token int) *ArrayIter {
	if part+1 > len(inv.layout) ||
		field+1 > len(inv.layout[part]) ||
		token+1 > len(inv.layout[part][field]) {
		return nil
	}
	return NewArrayIter(inv.layout[part][field][token])
}

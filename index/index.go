package index

type Ref uint32

type Token []Ref

type Folder []Token

type TokenFolder struct {
	token  int
	folder int
}

type Element struct {
	ref    Ref
	parts  []int
	tokens []TokenFolder
}

type Indexer interface {
	Reset()
	Next() (*Element, bool)
}

type Index struct {
	tokens []Folder
	source Indexer
}

func (ix *Index) Rebuild() {
	src := &ix.source
	src.Reset()
	for {
		if elem, ok := src.Next(); ok {
			// ...
			continue
		}
		break
	}
}

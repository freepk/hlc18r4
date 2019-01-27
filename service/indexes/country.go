package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type countryIndexer struct {
	pos int
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newCountryIndexer(rep *repo.AccountsRepo) *countryIndexer {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 5)}
	return &countryIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *countryIndexer) Reset() {
	ix.pos = 0
}

func (ix *countryIndexer) Next() (*inverted.Document, bool) {
	if id, ok := ix.next(); ok {
		return ix.processDocument(id), true
	}
	return nil, false
}

func (ix *countryIndexer) next() (int, bool) {
	n := ix.rep.Len()
	acc := proto.Account{}
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		acc = *ix.rep.Get(id)
		if acc.Email.Len > 0 {
			ix.pos = i + 1
			return id, true
		}
	}
	return 0, false
}

func (ix *countryIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for field := range doc.Tokens {
		doc.Tokens[field] = doc.Tokens[field][:0]
	}
	return doc
}

func (ix *countryIndexer) processDocument(id int) *inverted.Document {
	acc := *ix.rep.Get(id)
	doc := ix.resetDocument()
	doc.ID = 2000000 - id
	if acc.Country > 0 {
		doc.Parts = append(doc.Parts, NotNullToken, int(acc.Country))
	} else {
		doc.Parts = append(doc.Parts, NullToken)
	}
	switch acc.Sex {
	case proto.MaleSex:
		doc.Tokens[sexField] = append(doc.Tokens[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Tokens[sexField] = append(doc.Tokens[sexField], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], ComplToken, NotSingleToken, NotInRelToken)
	}
	return doc
}

type CountryIndex struct {
	inv *inverted.Inverted
}

func NewCountryIndex(rep *repo.AccountsRepo) *CountryIndex {
	src := newCountryIndexer(rep)
	inv := inverted.NewInverted(src)
	return &CountryIndex{inv: inv}
}

func (idx *CountryIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *CountryIndex) Sex(country, sex int) *inverted.ArrayIter {
	return idx.inv.Iterator(country, sexField, sex)
}

func (idx *CountryIndex) Status(country, status int) *inverted.ArrayIter {
	return idx.inv.Iterator(country, statusField, status)
}

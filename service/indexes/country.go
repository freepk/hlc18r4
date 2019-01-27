package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type countryIter struct {
	pos int
	acc *proto.Account
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newCountryIter(rep *repo.AccountsRepo) *countryIter {
	acc := &proto.Account{}
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Fields: make([][]int, 5)}
	return &countryIter{pos: 0, acc: acc, doc: doc, rep: rep}
}

func (it *countryIter) Reset() {
	it.pos = 0
}

func (it *countryIter) Next() (*inverted.Document, bool) {
	n := it.rep.Len()
	for i := it.pos; i < n; i++ {
		id := n - i - 1
		*it.acc = *it.rep.Get(id)
		if it.acc.Email.Len > 0 {
			it.pos = i + 1
			return it.processDocument(id, it.acc), true
		}
	}
	return nil, false
}

func (it *countryIter) resetDocument() *inverted.Document {
	doc := it.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for field := range doc.Fields {
		doc.Fields[field] = doc.Fields[field][:0]
	}
	return doc
}

func (it *countryIter) processDocument(id int, acc *proto.Account) *inverted.Document {
	doc := it.resetDocument()
	doc.ID = 2000000 - id
	if acc.Country > 0 {
		doc.Parts = append(doc.Parts, NotNullToken, int(acc.Country))
	} else {
		doc.Parts = append(doc.Parts, NullToken)
	}
	switch acc.Sex {
	case proto.MaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], ComplToken, NotSingleToken, NotInRelToken)
	}
	return doc
}

type CountryIndex struct {
	inv *inverted.Inverted
}

func NewCountryIndex(rep *repo.AccountsRepo) *CountryIndex {
	src := newCountryIter(rep)
	inv := inverted.NewInverted(src)
	return &CountryIndex{inv: inv}
}

func (idx *CountryIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *CountryIndex) Sex(country, sex int) *inverted.TokenIter {
	return idx.inv.Iterator(country, sexField, sex)
}

func (idx *CountryIndex) Status(country, status int) *inverted.TokenIter {
	return idx.inv.Iterator(country, statusField, status)
}

package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type statusIter struct {
	accountsIter
}

func newStatusIter(rep *repo.AccountsRepo) *statusIter {
	return &statusIter{accountsIter: *newAccountsIter(rep, 4)}
}

func (it *statusIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *statusIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
	switch acc.Status {
	case proto.SingleStatus:
		doc.Parts = append(doc.Parts, SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Parts = append(doc.Parts, InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Parts = append(doc.Parts, ComplToken, NotSingleToken, NotInRelToken)
	}
	switch acc.Sex {
	case proto.MaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], FemaleToken)
	}
	if acc.Country > 0 {
		doc.Fields[countryField] = append(doc.Fields[countryField], NotNullToken, int(acc.Country))
	} else {
		doc.Fields[countryField] = append(doc.Fields[countryField], NullToken)
	}
	if acc.City > 0 {
		doc.Fields[cityField] = append(doc.Fields[cityField], NotNullToken, int(acc.City))
	} else {
		doc.Fields[cityField] = append(doc.Fields[cityField], NullToken)
	}
	return doc
}

type StatusIndex struct {
	inv *inverted.Inverted
}

func NewStatusIndex(rep *repo.AccountsRepo) *StatusIndex {
	src := newStatusIter(rep)
	inv := inverted.NewInverted(src)
	return &StatusIndex{inv: inv}
}

func (idx *StatusIndex) Rebuild() {
	idx.inv.Rebuild()
}

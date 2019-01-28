package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type sexIter struct {
	accountsIter
}

func newSexIter(rep *repo.AccountsRepo) *sexIter {
	return &sexIter{accountsIter: *newAccountsIter(rep, 4)}
}

func (it *sexIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *sexIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
	switch acc.Sex {
	case proto.MaleSex:
		doc.Parts = append(doc.Parts, MaleToken)
	case proto.FemaleSex:
		doc.Parts = append(doc.Parts, FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], ComplToken, NotSingleToken, NotInRelToken)
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

type SexIndex struct {
	inv *inverted.Inverted
}

func NewSexIndex(rep *repo.AccountsRepo) *SexIndex {
	src := newSexIter(rep)
	inv := inverted.NewInverted(src)
	return &SexIndex{inv: inv}
}

func (idx *SexIndex) Rebuild() {
	idx.inv.Rebuild()
}

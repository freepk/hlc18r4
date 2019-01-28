package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type cityIter struct {
	accountsIter
}

func newCityIter(rep *repo.AccountsRepo) *cityIter {
	return &cityIter{accountsIter: *newAccountsIter(rep, 5)}
}

func (it *cityIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *cityIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
	if acc.City > 0 {
		doc.Parts = append(doc.Parts, NotNullToken, int(acc.City))
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
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[interestField] = append(doc.Fields[interestField], int(acc.Interests[i]))
	}
	return doc
}

type CityIndex struct {
	inv *inverted.Inverted
}

func NewCityIndex(rep *repo.AccountsRepo) *CityIndex {
	src := newCityIter(rep)
	inv := inverted.NewInverted(src)
	return &CityIndex{inv: inv}
}

func (idx *CityIndex) Rebuild() {
	idx.inv.Rebuild()
}

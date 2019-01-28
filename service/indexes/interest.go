package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type interestIter struct {
	accountsIter
}

func newInterestIter(rep *repo.AccountsRepo) *interestIter {
	return &interestIter{accountsIter: *newAccountsIter(rep, 5)}
}

func (it *interestIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *interestIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Parts = append(doc.Parts, int(acc.Interests[i]))
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

type InterestIndex struct {
	inv *inverted.Inverted
}

func NewInterestIndex(rep *repo.AccountsRepo) *InterestIndex {
	src := newInterestIter(rep)
	inv := inverted.NewInverted(src)
	return &InterestIndex{inv: inv}
}

func (idx *InterestIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *InterestIndex) SexIter(interest, sex int) *inverted.TokenIter {
	return idx.inv.Iterator(interest, sexField, sex)
}

func (idx *InterestIndex) StatusIter(interest, status int) *inverted.TokenIter {
	return idx.inv.Iterator(interest, statusField, status)
}

func (idx *InterestIndex) CountryIter(interest, country int) *inverted.TokenIter {
	return idx.inv.Iterator(interest, countryField, country)
}

func (idx *InterestIndex) CityIter(interest, city int) *inverted.TokenIter {
	return idx.inv.Iterator(interest, cityField, city)
}

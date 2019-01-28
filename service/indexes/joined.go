package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type joinedIter struct {
	accountsIter
}

func newJoinedIter(rep *repo.AccountsRepo) *joinedIter {
	return &joinedIter{accountsIter: *newAccountsIter(rep, 5)}
}

func (it *joinedIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *joinedIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
	doc.Parts = append(doc.Parts, yearTokenTS(int(acc.BirthTS)))
	switch acc.Sex {
	case proto.MaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], SingleToken)
	case proto.InRelStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], InRelToken)
	case proto.ComplStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], ComplToken)
	}
	if acc.Country > 0 {
		doc.Fields[countryField] = append(doc.Fields[countryField], int(acc.Country))
	}
	if acc.City > 0 {
		doc.Fields[cityField] = append(doc.Fields[cityField], int(acc.City))
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[interestField] = append(doc.Fields[interestField], int(acc.Interests[i]))
	}
	return doc
}

type JoinedIndex struct {
	inv *inverted.Inverted
}

func NewJoinedIndex(rep *repo.AccountsRepo) *JoinedIndex {
	src := newJoinedIter(rep)
	inv := inverted.NewInverted(src)
	return &JoinedIndex{inv: inv}
}

func (idx *JoinedIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *JoinedIndex) Sex(year, sex int) *inverted.Token {
	return nil
}

func (idx *JoinedIndex) Status(year, status int) *inverted.Token {
	return nil
}

func (idx *JoinedIndex) Country(year, country int) *inverted.Token {
	return nil
}

func (idx *JoinedIndex) City(year, city int) *inverted.Token {
	return nil
}

func (idx *JoinedIndex) Interests(year, interest int) *inverted.Token {
	return nil
}

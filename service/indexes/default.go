package indexes

import (
	"time"

	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type defaultIndexer struct {
	pos int
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newDefaultIndexer(rep *repo.AccountsRepo) *defaultIndexer {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 6)}
	return &defaultIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *defaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *defaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for field := range doc.Tokens {
		doc.Tokens[field] = doc.Tokens[field][:0]
	}
	return doc
}

func (ix *defaultIndexer) processDocument(id int, acc *proto.Account) *inverted.Document {
	doc := ix.resetDocument()
	doc.ID = id
	doc.Parts = append(doc.Parts, defaultPartition)
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
	if acc.Country > 0 {
		doc.Tokens[countryField] = append(doc.Tokens[countryField], NotNullToken, int(acc.Country))
	} else {
		doc.Tokens[countryField] = append(doc.Tokens[countryField], NullToken)
	}
	if acc.City > 0 {
		doc.Tokens[cityField] = append(doc.Tokens[cityField], NotNullToken, int(acc.City))
	} else {
		doc.Tokens[cityField] = append(doc.Tokens[cityField], NullToken)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Tokens[interestField] = append(doc.Tokens[interestField], int(acc.Interests[i]))
	}
	if _, birth, ok := parse.ParseInt(acc.Birth[:]); ok {
		birthYear := time.Unix(int64(birth), 0).UTC().Year() - 1975
		doc.Tokens[birthYearField] = append(doc.Tokens[birthYearField], birthYear)
	}
	return doc
}

func (ix *defaultIndexer) Next() (*inverted.Document, bool) {
	n := ix.rep.Len()
	acc := &proto.Account{}
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		*acc = *ix.rep.Get(id)
		if acc.Email.Len > 0 {
			ix.pos = i + 1
			pseudo := 2000000 - id
			return ix.processDocument(pseudo, acc), true
		}
	}
	return nil, false
}

type DefaultIndex struct {
	inv *inverted.Inverted
}

func NewDefaultIndex(rep *repo.AccountsRepo) *DefaultIndex {
	src := newDefaultIndexer(rep)
	inv := inverted.NewInverted(src)
	return &DefaultIndex{inv: inv}
}

func (idx *DefaultIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *DefaultIndex) Sex(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, sexField, token)
}

func (idx *DefaultIndex) Status(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, statusField, token)
}

func (idx *DefaultIndex) Country(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, countryField, token)
}

func (idx *DefaultIndex) City(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, cityField, token)
}

func (idx *DefaultIndex) Interest(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, interestField, token)
}

func (idx *DefaultIndex) BirthYear(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, birthYearField, token)
}

package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type defaultIndexer struct {
	pos int
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newDefaultIndexer(rep *repo.AccountsRepo) *defaultIndexer {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 5)}
	return &defaultIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *defaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *defaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for index := range doc.Tokens {
		doc.Tokens[index] = doc.Tokens[index][:0]
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
	return doc
}

func (ix *defaultIndexer) Next() (*inverted.Document, bool) {
	n := ix.rep.Len()
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		acc := ix.rep.Get(id)
		if acc.Email.Len > 0 {
			ix.pos = i + 1
			pseudo := 2000000 - id
			return ix.processDocument(pseudo, acc), true
		}
	}
	return nil, false
}

type DefaultIndex struct {
	index *inverted.Inverted
}

func NewDefaultIndex(rep *repo.AccountsRepo) *DefaultIndex {
	indexer := newDefaultIndexer(rep)
	index := inverted.NewInverted(indexer)
	return &DefaultIndex{index: index}
}

func (idx *DefaultIndex) Rebuild() {
	idx.index.Rebuild()
}

func (idx *DefaultIndex) Sex(token int) *inverted.ArrayIter {
	return idx.index.Iterator(defaultPartition, sexField, token)
}

func (idx *DefaultIndex) Status(token int) *inverted.ArrayIter {
	return idx.index.Iterator(defaultPartition, statusField, token)
}

func (idx *DefaultIndex) Country(token int) *inverted.ArrayIter {
	return idx.index.Iterator(defaultPartition, countryField, token)
}

func (idx *DefaultIndex) City(token int) *inverted.ArrayIter {
	return idx.index.Iterator(defaultPartition, cityField, token)
}
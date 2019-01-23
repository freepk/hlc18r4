package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	defaultPartition = 0
)

const (
	sexIndex = iota
	statusIndex
	countryIndex
	cityIndex
	interestIndex
)

const (
	NullToken = iota
	NotNullToken
)

const (
	MaleToken = iota
	FemaleToken
)

const (
	SingleToken = iota
	InRelToken
	ComplToken
	NotSingleToken
	NotInRelToken
	NotComplToken
)

type DefaultIndexer struct {
	pos int
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func NewDefaultIndexer(rep *repo.AccountsRepo) *DefaultIndexer {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 5)}
	return &DefaultIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *DefaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *DefaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for index := range doc.Tokens {
		doc.Tokens[index] = doc.Tokens[index][:0]
	}
	return doc
}

func (ix *DefaultIndexer) processDocument(id int, acc *proto.Account) *inverted.Document {
	doc := ix.resetDocument()
	doc.ID = id
	doc.Parts = append(doc.Parts, defaultPartition)
	switch acc.Sex {
	case proto.MaleSex:
		doc.Tokens[sexIndex] = append(doc.Tokens[sexIndex], MaleToken)
	case proto.FemaleSex:
		doc.Tokens[sexIndex] = append(doc.Tokens[sexIndex], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Tokens[statusIndex] = append(doc.Tokens[statusIndex], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Tokens[statusIndex] = append(doc.Tokens[statusIndex], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Tokens[statusIndex] = append(doc.Tokens[statusIndex], ComplToken, NotSingleToken, NotInRelToken)
	}
	if acc.Country > 0 {
		doc.Tokens[countryIndex] = append(doc.Tokens[countryIndex], NotNullToken, int(acc.Country))
	} else {
		doc.Tokens[countryIndex] = append(doc.Tokens[countryIndex], NullToken)
	}
	if acc.City > 0 {
		doc.Tokens[cityIndex] = append(doc.Tokens[cityIndex], NotNullToken, int(acc.City))
	} else {
		doc.Tokens[cityIndex] = append(doc.Tokens[cityIndex], NullToken)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Tokens[interestIndex] = append(doc.Tokens[interestIndex], int(acc.Interests[i]))
	}
	return doc
}

func (ix *DefaultIndexer) Next() (*inverted.Document, bool) {
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

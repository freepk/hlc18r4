package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	commonPartition = 0
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
	doc := &inverted.Document{ID: 0, Partitions: make([]int, 1), Indexes: make([][]int, 5)}
	return &DefaultIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *DefaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *DefaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Partitions = doc.Partitions[:0]
	for i := range doc.Indexes {
		doc.Indexes[i] = doc.Indexes[i][:0]
	}
	return doc
}

func (ix *DefaultIndexer) processDocument(id int, acc *proto.Account) *inverted.Document {
	doc := ix.resetDocument()
	doc.ID = id
	doc.Partitions = append(doc.Partitions, commonPartition)
	switch acc.Sex {
	case proto.MaleSex:
		doc.Indexes[sexIndex] = append(doc.Indexes[sexIndex], MaleToken)
	case proto.FemaleSex:
		doc.Indexes[sexIndex] = append(doc.Indexes[sexIndex], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Indexes[statusIndex] = append(doc.Indexes[statusIndex], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Indexes[statusIndex] = append(doc.Indexes[statusIndex], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Indexes[statusIndex] = append(doc.Indexes[statusIndex], ComplToken, NotSingleToken, NotInRelToken)
	}
	if acc.Country > 0 {
		doc.Indexes[countryIndex] = append(doc.Indexes[countryIndex], NotNullToken, int(acc.Country))
	} else {
		doc.Indexes[countryIndex] = append(doc.Indexes[countryIndex], NullToken)
	}
	if acc.City > 0 {
		doc.Indexes[cityIndex] = append(doc.Indexes[cityIndex], NotNullToken, int(acc.City))
	} else {
		doc.Indexes[cityIndex] = append(doc.Indexes[cityIndex], NullToken)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Indexes[interestIndex] = append(doc.Indexes[interestIndex], int(acc.Interests[i]))
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

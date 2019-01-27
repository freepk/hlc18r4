package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type accountsIter struct {
	pos int
	acc *proto.Account
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newAccountsIter(rep *repo.AccountsRepo, fields int) *accountsIter {
	acc := &proto.Account{}
	doc := &inverted.Document{ID: 0, Parts: make([]int, 0), Fields: make([][]int, fields)}
	return &accountsIter{pos: 0, acc: acc, doc: doc, rep: rep}
}

func (it *accountsIter) Reset() {
	it.pos = 0
}

func (it *accountsIter) next() bool {
	n := it.rep.Len()
	for i := it.pos; i < n; i++ {
		id := n - i - 1
		*it.acc = *it.rep.Get(id)
		if it.acc.Email.Len > 0 {
			it.pos = i + 1
			return true
		}
	}
	return false
}

func (it *accountsIter) id() int {
	return it.rep.Len() - it.pos
}

func (it *accountsIter) resetDocument() {
	it.doc.ID = 0
	it.doc.Parts = it.doc.Parts[:0]
	for field := range it.doc.Fields {
		it.doc.Fields[field] = it.doc.Fields[field][:0]
	}
}

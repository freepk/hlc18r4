package service

import (
	"gitlab.com/freepk/hlc18r4/index"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsIndexer struct {
	pos int
	doc *index.Document
	rep *repo.AccountsRepo
}

func NewAccountsIndexer(rep *repo.AccountsRepo) *AccountsIndexer {
	doc := &index.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 4)}
	return &AccountsIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *AccountsIndexer) Reset() {
	ix.pos = 0
}

func (ix *AccountsIndexer) resetDocument() *index.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for i := range doc.Tokens {
		doc.Tokens[i] = doc.Tokens[i][:0]
	}
	return doc
}

func (ix *AccountsIndexer) processDocument(id int, acc *proto.Account) *index.Document {
	doc := ix.resetDocument()
	doc.ID = id
	return doc
}

func (ix *AccountsIndexer) Next() (*index.Document, bool) {
	n := ix.rep.Len()
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		acc := ix.rep.Get(id)
		if acc.Email.Len > 0 {
			ix.pos = i + 1
			return ix.processDocument(2000000 - id, acc), true
		}
	}
	return nil, false
}

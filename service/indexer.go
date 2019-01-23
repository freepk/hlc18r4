package service

import (
	"gitlab.com/freepk/hlc18r4/index"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type AccountsIndexer struct {
	last int
	doc  *index.Document
	rep  *repo.AccountsRepo
}

func NewAccountsIndexer(rep *repo.AccountsRepo, parts, tokens int) *AccountsIndexer {
	doc := &index.Document{ID: 0, Parts: make([]int, parts), Tokens: make([][]int, tokens)}
	return &AccountsIndexer{last: 0, doc: doc, rep: rep}
}

func (ix *AccountsIndexer) Reset() {
	ix.last = 0
}

func (ix *AccountsIndexer) resetDocument() {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for i := range doc.Tokens {
		doc.Tokens[i] = doc.Tokens[i][:0]
	}
}

func (ix *AccountsIndexer) processDocument(acc *proto.Account) *index.Document {
	ix.resetDocument()
	return ix.doc
}

func (ix *AccountsIndexer) Next() (*index.Document, bool) {
	n := ix.rep.Len()
	for i := ix.last; i < n; i++ {
		acc := ix.rep.Get(n - i - 1)
		if acc.Email.Len > 0 {
			return ix.doc, true
		}
	}
	return nil, false
}

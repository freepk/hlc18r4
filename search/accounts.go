package search

import (
	"github.com/freepk/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type ProcessFunc func(doc *inverted.Document, acc *proto.Account)

type AccountsProcessor struct {
	rep  *repo.AccountsRepo
	doc  *inverted.Document
	acc  *proto.Account
	proc ProcessFunc
	i    int
}

func NewAccountsProcessor(rep *repo.AccountsRepo, proc ProcessFunc) *AccountsProcessor {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 0, 16), Fields: make([][]int, 4)}
	acc := &proto.Account{}
	return &AccountsProcessor{rep: rep, doc: doc, acc: acc, proc: proc, i: 0}
}

func (prc *AccountsProcessor) Reset() {
	prc.i = 0
}

func (prc *AccountsProcessor) Next() (*inverted.Document, bool) {
	n := prc.rep.Len()
	for i := prc.i; i < n; i++ {
		id := n - i - 1
		if prc.rep.IsSet(id) {
			*prc.acc = *prc.rep.Get(id)
			prc.doc.ID = 2000000 - id
			prc.proc(prc.doc, prc.acc)
			prc.i = i + 1
			return prc.doc, true
		}
	}
	return nil, false
}

package search

import (
	"github.com/freepk/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type ProcessFunc func(ref inverted.Ref, doc *inverted.Document, acc *proto.Account)

type AccountsProcessor struct {
	rep  *repo.AccountsRepo
	acc  *proto.Account
	doc  *inverted.Document
	proc ProcessFunc
	i    int
}

func NewAccountsProcessor(rep *repo.AccountsRepo, proc ProcessFunc) *AccountsProcessor {
	acc := &proto.Account{}
	doc := &inverted.Document{}
	return &AccountsProcessor{rep: rep, acc: acc, doc: doc, proc: proc, i: 0}
}

func (prc *AccountsProcessor) Reset() {
	prc.i = 0
}

func (prc *AccountsProcessor) Next() (*inverted.Document, bool) {
	n := prc.rep.Len()
	for i := prc.i; i < n; i++ {
		id := n - i - 1
		if prc.rep.IsSet(id) {
			ref := inverted.Ref(2000000 - id)
			*prc.acc = *prc.rep.Get(id)
			prc.proc(ref, prc.doc, prc.acc)
			prc.i = i + 1
			return prc.doc, true
		}
	}
	return nil, false
}

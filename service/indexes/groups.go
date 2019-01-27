package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type groupIter struct {
	pos int
	acc *proto.Account
	doc *inverted.Document
	rep *repo.AccountsRepo
}

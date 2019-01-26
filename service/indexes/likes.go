package indexes

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type LikesFromIndex struct {
	rep       *repo.AccountsRepo
	likesFrom [][]proto.Like
}

func NewLikesFromIndex(rep *repo.AccountsRepo) *LikesFromIndex {
	likesFrom := make([][]proto.Like, rep.Len())
	return &LikesFromIndex{rep: rep, likesFrom: likesFrom}
}

func (idx *LikesFromIndex) Rebuild() {
	n := idx.rep.Len()
	acc := &proto.Account{}
	for id := 0; id < n; id++ {
		*acc = *idx.rep.Get(id)
		_ = acc
	}
}

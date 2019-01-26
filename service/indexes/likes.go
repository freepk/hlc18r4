package indexes

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type LikeIter struct {
	likes []proto.Like
}

func NewLikeIter(likes []proto.Like) *LikeIter {
	return nil
}

type LikeIndex struct {
	rep   *repo.AccountsRepo
	likes [][]proto.Like
}

func NewLikeIndex(rep *repo.AccountsRepo) *LikeIndex {
	likes := make([][]proto.Like, rep.Len())
	return &LikeIndex{rep: rep, likes: likes}
}

func (idx *LikeIndex) Rebuild() {
	//n := idx.rep.Len()
	//acc := &proto.Account{}
	//for id := 0; id < n; id++ {
	//	*acc = *idx.rep.Get(id)
	//	_ = acc
	//}
}

func (idx *LikeIndex) Iterator(id int) *LikeIter {
	return nil
}

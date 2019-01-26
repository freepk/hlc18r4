package indexes

import (
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type LikeIter struct {
	likes []proto.Like
	i     int
}

func NewLikeIter(likes []proto.Like) *LikeIter {
	return &LikeIter{likes: likes, i: 0}
}

func (it *LikeIter) Reset() {
	it.i = 0
}

func (it *LikeIter) Next() (int, bool) {
	i := it.i
	if i < len(it.likes) {
		it.i++
		return int(it.likes[i].ID), true
	}
	return 0, false
}

type LikeIndex struct {
	rep   *repo.AccountsRepo
	likes [][]proto.Like
	last  uint32
}

func NewLikeIndex(rep *repo.AccountsRepo) *LikeIndex {
	likes := make([][]proto.Like, rep.Len())
	return &LikeIndex{rep: rep, likes: likes}
}

func (idx *LikeIndex) Rebuild() {
	n := idx.rep.Len()
	want := make([]uint32, n)
	acc := &proto.Account{}
	for id := 0; id < n; id++ {
		*acc = *idx.rep.Get(id)
		for i := range acc.LikesTo {
			if acc.LikesTo[i].TS > idx.last {
				want[id]++
			}
		}
	}
}

func (idx *LikeIndex) Iterator(id int) *LikeIter {
	return NewLikeIter(idx.likes[id])
}

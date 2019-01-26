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
	var acc proto.Account
	n := idx.rep.Len()
	need := make([]int, n)
	for id := 0; id < n; id++ {
		acc = *idx.rep.Get(id)
		for _, like := range acc.LikesTo {
			need[like.ID]++
		}
	}
	grow := 0
	for id, x := range need {
		if x > cap(idx.likes[id]) {
			grow += x * 105 / 100
		}
	}
	likes := make([]proto.Like, grow)
	for i := 0; i < n; i++ {
		id := n - i - 1
		acc = *idx.rep.Get(id)
		for _, like := range acc.LikesTo {
			x := need[like.ID]
			if x > cap(idx.likes[like.ID]) {
				idx.likes[like.ID], likes = likes[:0:x], likes[x:]
			}
			if x > len(idx.likes[like.ID]) {
				pseudo := 2000000 - uint32(id)
				idx.likes[like.ID] = append(idx.likes[like.ID], proto.Like{ID: pseudo, TS: like.TS})
			}
		}
	}
}

func (idx *LikeIndex) Iterator(id int) *LikeIter {
	return NewLikeIter(idx.likes[id])
}

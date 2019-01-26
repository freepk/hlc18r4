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
	/*
		n := idx.rep.Len()
		acc := &proto.Account{}
		grow := make([]int, n)
		for id := 0; id < n; id++ {
			*acc = *idx.rep.Get(id)
			for i := range acc.LikesTo {
				if acc.LikesTo[i].TS > idx.last {
					grow[id]++
				}
			}
		}
		total := 0
		for id := 0; id < n; id++ {
			need := grow[id] + len(idx.likes[id])
			if need > cap(idx.likes[id]) {
				total += need * 105 / 100
			}
		}
		likes := make([]proto.Like, total)
		for id := 0; id < n; id++ {
			if grow[id] == 0 {
				continue
			}
			need := grow[id] + len(idx.likes[id])
			if need > cap(idx.likes[id]) {
				need = need * 105 / 100
				n := copy(likes, idx.likes[id])
				idx.likes[id], likes = likes[:n:need], likes[need:]
			}
			*acc = *idx.rep.Get(id)
			for i := range acc.LikesTo {
				ts := acc.LikesTo[i].TS
				if ts > idx.last {
				}
			}
		}
	*/
}

func (idx *LikeIndex) Iterator(id int) *LikeIter {
	return NewLikeIter(idx.likes[id])
}

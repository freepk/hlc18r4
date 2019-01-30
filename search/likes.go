package search

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
	for liker := 0; liker < n; liker++ {
		acc = *idx.rep.Get(liker)
		for _, like := range acc.LikesTo {
			need[like.ID]++
		}
	}
	grow := 0
	for likee, x := range need {
		if x > cap(idx.likes[likee]) {
			grow += x * 105 / 100
		}
		if x > len(idx.likes[likee]) {
			idx.likes[likee] = idx.likes[likee][:0]
		}
	}
	likes := make([]proto.Like, grow)
	for i := 0; i < n; i++ {
		liker := n - i - 1
		acc = *idx.rep.Get(liker)
		for _, like := range acc.LikesTo {
			likee := like.ID
			x := need[likee]
			if x > cap(idx.likes[likee]) {
				x = x * 105 / 100
				idx.likes[likee] = likes[:0:x]
				likes = likes[x:]
			}
			if x > len(idx.likes[likee]) {
				pseudo := 2000000 - uint32(liker)
				next := proto.Like{ID: pseudo, TS: like.TS}
				idx.likes[likee] = append(idx.likes[likee], next)
			}
		}
	}
}

func (idx *LikeIndex) Iterator(id int) *LikeIter {
	return NewLikeIter(idx.likes[id])
}

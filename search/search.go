package search

import (
	"gitlab.com/freepk/hlc18r4/repo"
)

type SearchService struct {
	rep   *repo.AccountsRepo
	likes *LikeIndex
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	likes := NewLikeIndex(rep)
	return &SearchService{rep: rep, likes: likes}
}

func (svc *SearchService) Rebuild() {
	svc.likes.Rebuild()
}

func (svc *SearchService) Likes(id int) *LikeIter {
	return svc.likes.Iterator(id)
}

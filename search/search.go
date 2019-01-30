package search

import (
	"github.com/freepk/iterator"
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

func (svc *SearchService) Likes(id int) iterator.Iterator {
	if svc.rep.IsSet(id) {
		return svc.likes.Iterator(id)
	}
	return nil
}

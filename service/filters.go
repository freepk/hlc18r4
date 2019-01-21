package service

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/repo"
)

type FiltersService struct {
	rep     *repo.AccountsRepo
	indexes []*inverted.InvertedIndex
}

func NewFiltersService(rep *repo.AccountsRepo) *FiltersService {
	indexes := []*inverted.InvertedIndex{
		inverted.NewInvertedIndex(rep, inverted.DefaultParts, inverted.InterestsTokens)}
	return &FiltersService{rep: rep, indexes: indexes}
}

func (svc *FiltersService) RebuildIndexes() {
	svc.indexes[0].Rebuild()
}

func (svc *FiltersService) InterestsAny(part inverted.PartitionEnum, buf []byte) iterator.Iterator {
	index := svc.indexes[0]
	return iterator.NewUnionIter(index.Iterator(part, 10), index.Iterator(part, 15))
}

func (svc *FiltersService) InterestsContains(part inverted.PartitionEnum, buf []byte) iterator.Iterator {
	index := svc.indexes[0]
	return iterator.NewInterIter(index.Iterator(part, 10), index.Iterator(part, 15))
}

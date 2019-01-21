package service

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/repo"
)

type OperatorEnum int

const (
	IntersectOper = OperatorEnum(iota)
	UnionOper
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
}

func (svc *FiltersService) ByInterests(oper OperatorEnum, part inverted.PartitionEnum, buf []byte) {
}

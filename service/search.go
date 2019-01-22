package service

import (
	"sync"

	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

type SearchService struct {
	rep           *repo.AccountsRepo
	filterQueries *sync.Pool
}

func NewSearchService(rep *repo.AccountsRepo) *SearchService {
	return &SearchService{rep: rep, filterQueries: &sync.Pool{}}
}

func (svc *SearchService) FilterQuery() *FilterQuery {
	if query, ok := svc.filterQueries.Get().(*FilterQuery); ok {
		query.Reset()
		return query
	}
	return &FilterQuery{svc: svc}
}

type FilterQuery struct {
	svc     *SearchService
	sex     int
	country int
}

func (qry *FilterQuery) Reset() {
	qry.sex = 0
	qry.country = 0
}

func (qry *FilterQuery) Close() {
	qry.svc.filterQueries.Put(qry)
}

func (qry *FilterQuery) SexEq(sex []byte) bool {
	if sex, ok := proto.SexToken(sex); ok {
		qry.sex = sex
		return true
	}
	return false
}

func (qry *FilterQuery) CountryEq(country []byte) bool {
	if country, ok := proto.CountryToken(country); ok {
		qry.country = country
		return true
	}
	return false
}

func (qry *FilterQuery) CountryNull(null bool) {
	switch null {
	case true:
		qry.country = proto.NullToken
	case false:
		qry.country = proto.NotNullToken
	}
}

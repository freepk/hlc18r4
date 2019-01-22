package service

import (
	"testing"
)

func TestSearch(t *testing.T) {
	svc := NewSearchService(nil)
	qry := svc.FilterQuery()
	t.Log(qry)
	qry.SexEq([]byte(`m`))
	t.Log(qry)
	qry.Close()
	t.Log(qry)
	qry1 := svc.FilterQuery()
	t.Log(qry1)
	qry2 := svc.FilterQuery()
	t.Log(qry2)
	qry2.SexEq([]byte(`f`))
	t.Log(qry2)
	qry1.Close()
	qry2.Close()
	qry3 := svc.FilterQuery()
	qry4 := svc.FilterQuery()
	t.Log(qry3)
	t.Log(qry4)
}

func BenchmarkSearchFilterQuery(b *testing.B) {
	svc := NewSearchService(nil)
	for i := 0; i < b.N; i++ {
		qry := svc.FilterQuery()
		qry.Close()
	}
}

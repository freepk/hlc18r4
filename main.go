package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/service"
)

func main() {
	log.Println("Restoring accounts")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts", rep.Len())

	accountsSvc := service.NewAccountsService(rep)
	searchSvc := service.NewSearchService(rep)

	handler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		switch string(path) {
		case `/accounts/new/`:
			if accountsSvc.Create(ctx.PostBody()) {
				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		case `/accounts/likes/`:
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		case `/accounts/filter/`:
			query := searchSvc.FilterQuery()
			args := ctx.QueryArgs()
			args.VisitAll(func(k, v []byte) {
				switch string(k) {
				case `sex_eq`:
					query.SexEq(v)
				case `country_eq`:
					query.CountryEq(v)
				case `country_null`:
					query.CountryNull(v)
				}
			})
			return
		}
		path, id, ok := parse.ParseInt(path[10:])
		if !ok || !accountsSvc.Exists(id) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		switch string(path) {
		case `/`:
			if accountsSvc.Update(id, ctx.PostBody()) {
				ctx.SetStatusCode(fasthttp.StatusAccepted)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		case `suggetst/`:
		case `recommend/`:
		}
	}

	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}

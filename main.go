package main

import (
	"log"

	"github.com/freepk/iterator"
	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
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
			args := ctx.QueryArgs()
			if args.Len() < 3 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			limit, err := args.GetUint(`limit`)
			if err != nil || limit > 50 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			var iter iterator.Iterator
			hasErrors := false
			fields := proto.IDField
			args.VisitAll(func(k, v []byte) {
				var next iterator.Iterator
				switch string(k) {
				case `limit`:
					return
				case `query_id`:
					return
				case `sex_eq`:
					if next = accountsSvc.BySexEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.SexField
				case `status_eq`:
					if next = accountsSvc.ByStatusEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				case `status_neq`:
					if next = accountsSvc.ByStatusNeq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				default:
					hasErrors = true
					return
				}
				if iter == nil {
					iter = next
				} else {
					iter = iterator.NewInterIter(iter, next)
				}
			})
			if hasErrors {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if _, ok := iter.Next(); ok {
				ctx.WriteString(`{"accounts":[`)
				ctx.WriteString(`]}`)
			} else {
				ctx.WriteString(`{"accounts":[]}`)
			}
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

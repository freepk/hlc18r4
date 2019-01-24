package main

import (
	"log"

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
			errs := 0
			args := ctx.QueryArgs()
			args.VisitAll(func(k, v []byte) {
				switch string(k) {
				case `limit`:
					println(`limit`, string(v))
				case `query_id`:
					println(`query_id`, string(v))
				case `sex_eq`:
					println(`sex_eq`, string(v))
				case `status_eq`:
					println(`status_eq`, string(v))
				case `country_eq`:
					println(`country_eq`, string(v))
				case `city_eq`:
					println(`city_eq`, string(v))
				default:
					errs++
				}
			})
			if errs > 0 {
				return
			}
			if args.Len() < 3 {
				return
			} else if args.Len() > 3 {
				if args.Has(`country_eq`) {
					return
				} else if args.Has(`city_eq`) {
					return
				} else if args.GetBool(`country_null`) {
					return
				} else if args.Has(`status_eq`) {
					return
				} else if args.Has(`sex_eq`) {
					return
				}
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

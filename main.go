package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/service"
)

func AccountsHandler(ctx *fasthttp.RequestCtx, svc *service.AccountsService) {
	path := ctx.Path()
	if len(path) < httpBaseLen || string(path[:httpBaseLen]) != httpBasePath {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	path = path[httpBaseLen:]
	switch string(ctx.Method()) {
	case `GET`:
		switch string(path) {
		case `filter/`:
			errs := 0
			args := ctx.QueryArgs()
			args.VisitAll(func(k, v []byte) {
				switch string(k) {
				case `limit`:
				case `query_id`:
				case `sex_eq`:
				case `country_eq`:
				case `country_null`:
				case `status_eq`:
				case `status_neq`:
				case `interests_contains`:
				case `interests_any`:
				case `likes_contains`:
				case `city_eq`:
				case `city_any`:
				case `city_null`:
				case `birth_gt`:
				case `birth_lt`:
				case `birth_year`:
				case `premium_now`:
				case `premium_null`:
				case `email_gt`:
				case `email_lt`:
				case `email_domain`:
				case `sname_starts`:
				case `sname_null`:
				case `fname_null`:
				case `fname_any`:
				case `phone_null`:
				case `phone_code`:
				default:
					errs++
				}
			})
			if errs > 0 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		case `group`:
		}
	case `POST`:
		switch string(path) {
		case `likes/`:
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		case `new/`:
			acc := &proto.Account{}
			if _, ok := acc.UnmarshalJSON(ctx.PostBody()); !ok {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if !svc.Create(acc) {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusCreated)
			return
		default:
			_, id, ok := parse.ParseInt(path)
			if !ok || !svc.Exists(id) {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
			acc := &proto.Account{}
			if _, ok := acc.UnmarshalJSON(ctx.PostBody()); !ok {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if !svc.Update(id, acc) {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func main() {
	log.Println("Restore service")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	svc := service.NewAccountsService(rep)
	svc.Reindex()
	svc.Reindex()
	return
	handler := func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, svc)
	}
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}

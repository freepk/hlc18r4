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
			/*
				sex_eq
				country_eq
				country_null
				status_eq
				status_neq
				interests_contains
				interests_any
				likes_contains
				city_eq
				city_any
				city_null
				birth_gt
				birth_lt
				birth_year
				premium_now
				premium_null
				email_gt
				email_lt
				email_domain
				sname_starts
				sname_null
				fname_null
				fname_any
				phone_null
				phone_code
			*/
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
	handler := func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, svc)
	}
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}

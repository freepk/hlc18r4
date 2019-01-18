package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/service"
)

const (
	httpBasePath   = `/accounts/`
	httpBaseLen    = len(httpBasePath)
	httpNewPath    = `/accounts/new/`
	httpLikesPath  = `/accounts/likes/`
	httpFilterPath = `/accounts/filter/`
	httpGroupPath  = `/accounts/group/`
)

func AccountsHandler(ctx *fasthttp.RequestCtx, svc *service.AccountsService) {
	var id int
	var ok bool

	path := ctx.Path()
	switch string(path) {
	case httpNewPath:
		acc := &proto.Account{}
		if _, ok = acc.UnmarshalJSON(ctx.PostBody()); !ok {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if _, id, ok = parse.ParseInt(acc.ID[:]); !ok {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if !svc.Create(id, acc) {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	case httpLikesPath:
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		return
	case httpFilterPath:
	case httpGroupPath:
	default:
		if len(path) < httpBaseLen || string(path[:httpBaseLen]) != httpBasePath {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if path, id, ok = parse.ParseInt(path[httpBaseLen:]); !ok {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if ctx.IsPost() {
			if !svc.Exists(id) {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
			acc := &proto.Account{}
			if _, ok = acc.UnmarshalJSON(ctx.PostBody()); !ok {
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

package main

import (
	"log"

	"github.com/valyala/fasthttp"
	//"gitlab.com/freepk/hlc18r4/parse"
	//"gitlab.com/freepk/hlc18r4/proto"
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
	case `POST`:
		switch string(path) {
		case `likes/`:
		case `new/`:
		default:
		}
	}
}

func main() {
	log.Println("Restore service")
	svc, err := service.RestoreAccountsService("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	handler := func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, svc)
	}
	log.Println("Start listen")
	err = fasthttp.ListenAndServe(":80", handler)
	if err != nil {
		log.Fatal(err)
	}
}

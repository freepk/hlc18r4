package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/service"
)

var (
	writesCount uint64
	accountsSvc *service.AccountsService
)

func routerHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsPost() {
		atomic.StoreUint64(&writesCount, 1)
	}
	path := ctx.Path()
	switch string(path) {
	case `/accounts/filter/`:
		filterHandler(ctx)
	case `/accounts/group/`:
		groupHandler(ctx)
	case `/accounts/new/`:
		createHandler(ctx)
	default:
		path, id, ok := parse.ParseInt(path[10:])
		if ok {
			switch string(path) {
			case `/`:
				updateHandler(id, ctx)
			case `suggetst/`:
				suggestHandler(id, ctx)
			case `recommend/`:
				recommendHandler(id, ctx)
			}
		}
	}
}

func filterHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	limit, err := args.GetUint(`limit`)
	if err != nil || limit > 50 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	args.VisitAll(func(k, v []byte) {
		switch string(k) {
		case `sex_eq`:
			if t := accountsSvc.Default().SexEq(v); t == nil {
				log.Println(string(k), string(v), t.Len())
			}
		case `status_eq`:
			if t := accountsSvc.Default().StatusEq(v); t == nil {
				log.Println(string(k), string(v), t.Len())
			}
		case `status_neq`:
			if t := accountsSvc.Default().StatusNeq(v); t != nil {
				log.Println(string(k), string(v), t.Len())
			}
		case `email_domain`:
			if t := accountsSvc.Default().EmailDomain(v); t != nil {
				log.Println(string(k), string(v), t.Len())
			}
		}
	})
}

func groupHandler(ctx *fasthttp.RequestCtx) {
}

func createHandler(ctx *fasthttp.RequestCtx) {
}

func updateHandler(id int, ctx *fasthttp.RequestCtx) {
}

func suggestHandler(id int, ctx *fasthttp.RequestCtx) {
}

func recommendHandler(id int, ctx *fasthttp.RequestCtx) {
}

func main() {
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	accountsSvc = service.NewAccountsService(rep)
	accountsSvc.RebuildIndexes()
	go func() {
		writeProcess := false
		for {
			temp := atomic.LoadUint64(&writesCount)
			if temp > 0 {
				writeProcess = true
				atomic.StoreUint64(&writesCount, 0)
			} else if writeProcess {
				writeProcess = false
				accountsSvc.RebuildIndexes()
			}
			time.Sleep(time.Second)
		}
	}()
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", routerHandler); err != nil {
		log.Fatal(err)
	}
}

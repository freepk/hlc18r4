package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/freepk/parse"
	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/accounts"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/search"
)

var (
	writesCount uint64
	accountsSvc *accounts.AccountsService
	searchSvc   *search.SearchService
)

func routerHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsPost() {
		atomic.StoreUint64(&writesCount, 1)
	}
	path := ctx.Path()
	switch string(path) {
	case `/accounts/filter/`:
		filterHandler(ctx)
	case `/accounts/likes/`:
		likesHandler(ctx)
	case `/accounts/group/`:
		groupHandler(ctx)
	case `/accounts/new/`:
		createHandler(ctx)
	default:
		if path, id, ok := parse.ParseInt(path[10:]); ok {
			switch string(path) {
			case `/`:
				updateHandler(id, ctx)
			case `suggetst/`:
				suggestHandler(id, ctx)
			case `recommend/`:
				recommendHandler(id, ctx)
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
}

func likesHandler(ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.AddLikes(ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func createHandler(ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.Create(ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func updateHandler(id int, ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.Update(id, ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	} else if err == accounts.NotFoundError {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func filterHandler(ctx *fasthttp.RequestCtx) {
}

func groupHandler(ctx *fasthttp.RequestCtx) {
}

func suggestHandler(id int, ctx *fasthttp.RequestCtx) {
}

func recommendHandler(id int, ctx *fasthttp.RequestCtx) {
}

func main() {
	log.Println("Restoring")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts service")
	accountsSvc = accounts.NewAccountsService(rep)
	log.Println("Search service")
	searchSvc = search.NewSearchService(rep)
	searchSvc.Rebuild()
	go func() {
		writeProcess := false
		for {
			temp := atomic.LoadUint64(&writesCount)
			if temp > 0 {
				writeProcess = true
				atomic.StoreUint64(&writesCount, 0)
			} else if writeProcess {
				writeProcess = false
				log.Println("Write process finished")
				searchSvc.Rebuild()
			}
			time.Sleep(time.Second)
		}
	}()
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", routerHandler); err != nil {
		log.Fatal(err)
	}
}

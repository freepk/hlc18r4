package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/freepk/hlc18r4/accounts"
	"github.com/freepk/hlc18r4/backup"
	"github.com/freepk/hlc18r4/search"
	"github.com/freepk/iterator"
	"github.com/freepk/parse"
	"github.com/valyala/fasthttp"
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

type tokenFunc func([]byte) (int, bool)

type iterFunc func(int) iterator.Iterator

type operFunc func(iterator.Iterator, iterator.Iterator) iterator.Iterator

func interOper(a, b iterator.Iterator) iterator.Iterator {
	return iterator.NewInterIter(a, b)
}

func unionOper(a, b iterator.Iterator) iterator.Iterator {
	return iterator.NewUnionIter(a, b)
}

func buildIter(iter iterator.Iterator, vals []byte, tokenFn tokenFunc, iterFn iterFunc, operFn operFunc) (iterator.Iterator, bool) {
	var res iterator.Iterator
	vals, val := parse.ScanSymbol(vals, 0x2C)
	for len(val) > 0 {
		token, ok := tokenFn(val)
		if !ok {
			return iter, false
		}
		it := iterFn(token)
		if it == nil {
			return iter, false
		}
		if res != nil {
			res = operFn(res, it)
		}
		if res == nil {
			res = it
		}
		vals, val = parse.ScanSymbol(vals, 0x2C)
	}
	if res == nil {
		return iter, true
	}
	return iterator.NewInterIter(iter, res), true
}

func intToken(b []byte) (int, bool) {
	_, token, ok := parse.ParseInt(b)
	return token, ok
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
		}
	})
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

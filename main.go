package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/service"
)

func AccountsHandler(ctx *fasthttp.RequestCtx, svc *service.AccountsService) {
	var id int
	var ok bool

	path := ctx.Path()
	switch string(path) {
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
	case `/accounts/group/`:
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
	case `/accounts/new/`:
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
	case `/accounts/likes/`:
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		return
	default:
		if len(path) < 10 || string(path[:10]) != `/accounts/` {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if path, id, ok = parse.ParseInt(path[10:]); !ok {
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
	if true {
		if fd, err := os.Create("cpu.prof"); err == nil {
			pprof.StartCPUProfile(fd)
			defer pprof.StopCPUProfile()
		}
		defer func() {
			if fd, err := os.Create("mem.prof"); err == nil {
				runtime.GC()
				pprof.WriteHeapProfile(fd)
			}
		}()
	}
	log.Println("Restore service")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts", rep.Len())
	svc := service.NewAccountsService(rep)
	handler := func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, svc)
	}
	svc.RebuildIndexes()
	go func() {
		log.Println("Start listen")
		if err := fasthttp.ListenAndServe(":80", handler); err != nil {
			log.Fatal(err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

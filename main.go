package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/database"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

func AccountsHandler(ctx *fasthttp.RequestCtx, db *database.Database) {
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
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		case `new/`:
			account := &proto.Account{}
			if _, ok := account.UnmarshalJSON(ctx.PostBody()); ok {
				if account.ID > 0 {
					ctx.SetStatusCode(fasthttp.StatusCreated)
					return
				}
			}
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		default:
			if _, id, ok := parse.ParseInt(path); ok {
				_ = id
				account := &proto.Account{}
				if _, ok := account.UnmarshalJSON(ctx.PostBody()); ok {
					ctx.SetStatusCode(fasthttp.StatusAccepted)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		}
	case `GET`:
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func main() {
	//db, err := backup.Restore("/tmp/data/")
	db, err := backup.Restore("tmp/data/")
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	db.Reindex()

	err = fasthttp.ListenAndServe(":80", func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, db)
	})
	if err != nil {
		log.Fatal(err)
	}
}

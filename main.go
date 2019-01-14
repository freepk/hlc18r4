package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
)

func AccountsHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	if len(path) < httpBaseLen || string(path[:httpBaseLen]) != httpBasePath {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	switch string(ctx.Method()) {
	case `POST`:
	case `GET`:
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func main() {
	db, err := backup.Restore("./data/")
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	db.BuildIndexes()

	err = fasthttp.ListenAndServe(":80", AccountsHandler)
	if err != nil {
		log.Fatal(err)
	}
}

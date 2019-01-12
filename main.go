package main

import (
	"log"

	"github.com/valyala/fasthttp"
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
	err := fasthttp.ListenAndServe(":80", AccountsHandler)
	if err != nil {
		log.Println(err)
	}
}

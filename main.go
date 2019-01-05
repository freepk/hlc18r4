package main

import (
	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	if string(path[:10]) != "/accounts/" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	path = path[10:]
	switch string(ctx.Method()) {
	case "GET":
	case "POST":
		switch string(path) {
		case "likes/":
			ctx.SetStatusCode(fasthttp.StatusAccepted)
		case "new/":
			ctx.SetStatusCode(fasthttp.StatusCreated)
		default:
			ok := false
			_, path, ok = parseInt(path)
			if !ok {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
			path, ok = parseSymbol(path, '/')
			if !ok {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusAccepted)
		}
		ctx.SetBody([]byte("{}"))
	}
}

func main() {
	fasthttp.ListenAndServe(":80", fastHTTPHandler)
}

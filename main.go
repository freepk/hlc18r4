package main

import "github.com/valyala/fasthttp"

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Hello, I'm Web Server!"))
}

func main() {
	fasthttp.ListenAndServe(":80", fastHTTPHandler)
}

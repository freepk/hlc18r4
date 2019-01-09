package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
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
	if true {
		if fd, err := os.Create("cpu.prof"); err == nil {
			pprof.StartCPUProfile(fd)
			defer func() {
				pprof.StopCPUProfile()
				fd.Close()
			}()
		}
		defer func() {
			if fd, err := os.Create("mem.prof"); err == nil {
				runtime.GC()
				pprof.WriteHeapProfile(fd)
				fd.Close()
			}
		}()
	}

	runtime.GC()
	outputStatm()

	err := fasthttp.ListenAndServe(":80", fastHTTPHandler)
	if err != nil {
		log.Fatal(err)
	}
}

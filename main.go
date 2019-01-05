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
	if string(path[:10]) != "/accounts/" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	path = path[10:]
	switch string(ctx.Method()) {
	case "GET":
		switch string(path) {
		case "filter/":
		case "country/":
			for i := 0; i < 70; i++ {
				ctx.Write(countryLookup.GetItemNoLock(i))
			}
		}
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
	if 0 == 1 {
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

	go func() {
		log.Println(fasthttp.ListenAndServe(":80", fastHTTPHandler))
	}()

	// Block
	var c chan int
	<-c
}

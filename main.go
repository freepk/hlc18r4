package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/db"
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

	if true {
		if fd, err := os.Create("cpu.prof"); err == nil {
			log.Println("Start CPU profile.")
			pprof.StartCPUProfile(fd)
			defer func() {
				log.Println("Stop and close CPU profile.")
				pprof.StopCPUProfile()
				fd.Close()
			}()
		}
		defer func() {
			if fd, err := os.Create("mem.prof"); err == nil {
				log.Println("Write heap profile.")
				runtime.GC()
				pprof.WriteHeapProfile(fd)
				fd.Close()
			}
		}()
	}

	log.Println("Start load DB")
	mainDb := db.NewDB()
	//mainDb.Restore("/tmp/data/data.zip")
	mainDb.Restore("./data/data.zip")
	db.Print()

	runtime.GC()
	outputStatm()

	syscall.Mlockall(syscall.MCL_FUTURE)

	go func() {
		err := fasthttp.ListenAndServe(":80", AccountsHandler)
		if err != nil {
			log.Println(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

}

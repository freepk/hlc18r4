package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func outputStatm() {
	buf, err := ioutil.ReadFile("/proc/self/statm")
	if err != nil {
		log.Fatal(err)
	}
	split := bytes.Split(buf, []byte{0x20})
	fmt.Println("Statm: size", split[0], "resident", split[1], "share", split[2])
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

	//	runtime.GC()
	//	outputStatm()

	//	n := 20
	//	m := 2 * 1024 * 1024
	//	tempItemBuffer := make([]byte, 8)
	//	lookups := make([]*Lookup, 0, n)
	//	for i := 0; i < n; i++ {
	//		lookup := NewLookup(int32(m), 0)
	//		for j := 0; j < m; j++ {
	//			binary.LittleEndian.PutUint64(tempItemBuffer, uint64(j))
	//			index := lookup.GetIndexOrSet(tempItemBuffer)
	//			item := lookup.GetItemNoLock(index)
	//			if !bytes.Equal(tempItemBuffer, item) {
	//				log.Fatal("Items not equal")
	//			}
	//		}
	//		lookups = append(lookups, lookup)
	//	}

	runtime.GC()
	outputStatm()

	go func() {
		log.Println(fasthttp.ListenAndServe(":80", fastHTTPHandler))
	}()

	// Block
	var c chan int
	<-c
}

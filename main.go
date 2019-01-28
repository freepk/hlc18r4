package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/freepk/iterator"
	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/service"
)

var (
	writesCount uint64
	accountsSvc *service.AccountsService
)

func routerHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsPost() {
		atomic.StoreUint64(&writesCount, 1)
	}
	path := ctx.Path()
	switch string(path) {
	case `/accounts/filter/`:
		filterHandler(ctx)
	case `/accounts/group/`:
		groupHandler(ctx)
	case `/accounts/new/`:
		createHandler(ctx)
	default:
		path, id, ok := parse.ParseInt(path[10:])
		if ok {
			switch string(path) {
			case `/`:
				updateHandler(id, ctx)
			case `suggetst/`:
				suggestHandler(id, ctx)
			case `recommend/`:
				recommendHandler(id, ctx)
			}
		}
	}
}

func filterHandler(ctx *fasthttp.RequestCtx) {
	var iter iterator.Iterator
	fields := proto.IDField | proto.EmailField
	args := ctx.QueryArgs()
	limit, err := args.GetUint(`limit`)
	if err != nil || limit > 50 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	args.VisitAll(func(k, v []byte) {
		var next iterator.Iterator
		switch string(k) {
		case `sex_eq`:
			if t := accountsSvc.Default().Sex(v); t == nil {
				next = t.Iterator()
			}
			fields |= proto.SexField
		case `status_eq`:
			if t := accountsSvc.Default().Status(v); t == nil {
				next = t.Iterator()
			}
			fields |= proto.StatusField
		case `status_neq`:
			if t := accountsSvc.Default().NotStatus(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.StatusField
		case `email_domain`:
			if t := accountsSvc.Default().EmailDomain(v); t != nil {
				next = t.Iterator()
			}
		case `fname_eq`:
			if t := accountsSvc.Default().Fname(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.FnameField
		case `fname_null`:
			if t := accountsSvc.Default().FnameNull(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.FnameField
		case `sname_eq`:
			if t := accountsSvc.Default().Sname(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.SnameField
		case `sname_null`:
			if t := accountsSvc.Default().SnameNull(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.SnameField
		case `phone_code`:
			if t := accountsSvc.Default().PhoneCode(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.PhoneField
		case `phone_null`:
			if t := accountsSvc.Default().PhoneNull(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.PhoneField
		case `country_eq`:
			if t := accountsSvc.Default().Country(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.CountryField
		case `country_null`:
			if t := accountsSvc.Default().CountryNull(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.CountryField
		case `city_eq`:
			if t := accountsSvc.Default().City(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.CityField
		case `city_null`:
			if t := accountsSvc.Default().CityNull(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.CityField
		case `birth_year`:
			if t := accountsSvc.Default().BirthYear(v); t != nil {
				next = t.Iterator()
			}
			fields |= proto.BirthField
		}
		if next != nil {
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewInterIter(next, iter)
			}
		}
	})
	if iter == nil {
		return
	}
	println(string(args.QueryString()))
	acc := &proto.Account{}
	comma := false
	ctx.WriteString(`{"accounts":[`)
	for limit > 0 {
		pseudo, ok := iter.Next()
		if !ok {
			break
		}
		*acc = *accountsSvc.Get(2000000 - pseudo)
		if comma {
			ctx.WriteString(`,`)
		}
		comma = true
		acc.WriteJSON(fields, ctx)
		limit--
	}
	ctx.WriteString(`]}`)
}

func groupHandler(ctx *fasthttp.RequestCtx) {
}

func createHandler(ctx *fasthttp.RequestCtx) {
}

func updateHandler(id int, ctx *fasthttp.RequestCtx) {
}

func suggestHandler(id int, ctx *fasthttp.RequestCtx) {
}

func recommendHandler(id int, ctx *fasthttp.RequestCtx) {
}

func main() {
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	accountsSvc = service.NewAccountsService(rep)
	accountsSvc.RebuildIndexes()
	go func() {
		writeProcess := false
		for {
			temp := atomic.LoadUint64(&writesCount)
			if temp > 0 {
				writeProcess = true
				atomic.StoreUint64(&writesCount, 0)
			} else if writeProcess {
				writeProcess = false
				accountsSvc.RebuildIndexes()
			}
			time.Sleep(time.Second)
		}
	}()
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", routerHandler); err != nil {
		log.Fatal(err)
	}
}

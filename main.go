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
	hasErrors := false
	birthLT := 0
	birthGT := 0
	emailLT := []byte{}
	emailGT := []byte{}
	snameStarts := []byte{}
	args.VisitAll(func(k, v []byte) {
		var next iterator.Iterator
		switch string(k) {
		case `sex_eq`:
			if it := accountsSvc.Default().Sex(v); it != nil {
				next = it
				fields |= proto.SexField
			}
		case `status_eq`:
			if it := accountsSvc.Default().Status(v); it != nil {
				next = it
				fields |= proto.StatusField
			}
		case `status_neq`:
			if it := accountsSvc.Default().NotStatus(v); it != nil {
				next = it
				fields |= proto.StatusField
			}
		case `email_domain`:
			if it := accountsSvc.Default().EmailDomain(v); it != nil {
				next = it
			}
		case `email_lt`:
			emailLT = v
		case `email_gt`:
			emailGT = v
		case `fname_eq`:
			if it := accountsSvc.Default().Fname(v); it != nil {
				next = it
				fields |= proto.FnameField
			}
		case `fname_null`:
			if it := accountsSvc.Default().FnameNull(v); it != nil {
				next = it
				fields |= proto.FnameField
			}
		case `sname_eq`:
			if it := accountsSvc.Default().Sname(v); it != nil {
				next = it
				fields |= proto.SnameField
			}
		case `sname_starts`:
			snameStarts = v
			fields |= proto.SnameField
		case `sname_null`:
			if it := accountsSvc.Default().SnameNull(v); it != nil {
				next = it
				fields |= proto.SnameField
			}
		case `phone_code`:
			if it := accountsSvc.Default().PhoneCode(v); it != nil {
				next = it
				fields |= proto.PhoneField
			}
		case `phone_null`:
			if it := accountsSvc.Default().PhoneNull(v); it != nil {
				next = it
				fields |= proto.PhoneField
			}
		case `country_eq`:
			if it := accountsSvc.Default().Country(v); it != nil {
				next = it
				fields |= proto.CountryField
			}
		case `country_null`:
			if it := accountsSvc.Default().CountryNull(v); it != nil {
				next = it
				fields |= proto.CountryField
			}
		case `city_eq`:
			if it := accountsSvc.Default().City(v); it != nil {
				next = it
				fields |= proto.CityField
			}
		case `city_null`:
			if it := accountsSvc.Default().CityNull(v); it != nil {
				next = it
				fields |= proto.CityField
			}
		case `birth_year`:
			if it := accountsSvc.Default().BirthYear(v); it != nil {
				next = it
				fields |= proto.BirthField
			}
		case `birth_lt`:
			if _, ts, ok := parse.ParseInt(v); ok {
				birthLT = ts
				fields |= proto.BirthField
			}
		case `birth_gt`:
			if _, ts, ok := parse.ParseInt(v); ok {
				birthGT = ts
				fields |= proto.BirthField
			}
		case `premium_now`:
		case `premium_null`:
		case `interests_any`:
		case `interests_contains`:
		case `fname_any`:
		case `city_any`:
		case `likes_contains`:
		case `limit`:
		case `query_id`:
		default:
			hasErrors = true
			return
		}
		if next != nil {
			if iter == nil {
				iter = next
			} else {
				iter = iterator.NewInterIter(next, iter)
			}
		}
	})
	if hasErrors {
		return
	}
	if iter == nil {
		return
	}
	acc := &proto.Account{}
	comma := false
	ctx.WriteString(`{"accounts":[`)
	for limit > 0 {
		pseudo, ok := iter.Next()
		if !ok {
			break
		}
		*acc = *accountsSvc.Get(2000000 - pseudo)
		if birthLT > 0 && birthLT < int(acc.BirthTS) {
			continue
		}
		if birthGT > 0 && birthGT > int(acc.BirthTS) {
			continue
		}
		if len(emailLT) > 0 && string(emailLT) < string(acc.Email.Buf[:acc.Email.Len]) {
			continue
		}
		if len(emailGT) > 0 && string(emailGT) > string(acc.Email.Buf[:acc.Email.Len]) {
			continue
		}
		if n := len(snameStarts); n > 0 {
			sname := acc.GetSname()
			if n > len(sname) || string(snameStarts) != string(sname[:n]) {
				continue
			}
		}
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
	if accountsSvc.Create(ctx.PostBody()) {
		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func updateHandler(id int, ctx *fasthttp.RequestCtx) {
	if accountsSvc.Update(id, ctx.PostBody()) {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
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

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

func main() {
	log.Println("Restoring accounts")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts", rep.Len())

	accountsSvc := service.NewAccountsService(rep)

	writesCount := uint64(0)
	go func() {
		writeProcess := false
		for {
			temp := atomic.LoadUint64(&writesCount)
			if temp > 0 {
				writeProcess = true
				atomic.StoreUint64(&writesCount, 0)
			} else if writeProcess {
				writeProcess = false
				log.Println("Rebuild indexes")
				accountsSvc.RebuildIndexes()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	handler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		if ctx.IsPost() {
			atomic.StoreUint64(&writesCount, 1)
		}
		switch string(path) {
		case `/accounts/new/`:
			if accountsSvc.Create(ctx.PostBody()) {
				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		case `/accounts/likes/`:
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		case `/accounts/filter/`:
			args := ctx.QueryArgs()
			if args.Len() < 3 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			limit, err := args.GetUint(`limit`)
			if err != nil || limit > 50 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			var iter iterator.Iterator
			hasErrors := false
			fields := proto.IDField | proto.EmailField
			args.VisitAll(func(k, v []byte) {
				var next iterator.Iterator
				switch string(k) {
				case `limit`:
					return
				case `query_id`:
					return
				case `sex_eq`:
					if next = accountsSvc.BySexEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.SexField
				case `status_eq`:
					if next = accountsSvc.ByStatusEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				case `status_neq`:
					if next = accountsSvc.ByStatusNeq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				case `country_eq`:
					if next = accountsSvc.ByCountryEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.CountryField
				case `country_null`:
					if next = accountsSvc.ByCountryNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.CountryField
				case `city_eq`:
					if next = accountsSvc.ByCityEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.CityField
				case `city_null`:
					if next = accountsSvc.ByCityNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.CityField
				case `city_any`:
					if next = accountsSvc.ByCityAny(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.CityField
				case `interests_any`:
					if next = accountsSvc.ByInterestsAny(v); next == nil {
						hasErrors = true
						return
					}
				case `interests_contains`:
					if next = accountsSvc.ByInterestsContains(v); next == nil {
						hasErrors = true
						return
					}
				default:
					hasErrors = true
					return
				}
				if iter == nil {
					iter = next
				} else {
					iter = iterator.NewInterIter(iter, next)
				}
			})
			if hasErrors {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			limit--
			if pseudo, ok := iter.Next(); ok {
				ctx.WriteString(`{"accounts":[`)
				acc := accountsSvc.Get(2000000 - pseudo)
				acc.WriteJSON(fields, ctx)
				for limit > 0 {
					limit--
					if pseudo, ok = iter.Next(); !ok {
						break
					}
					ctx.WriteString(`,`)
					*acc = *accountsSvc.Get(2000000 - pseudo)
					acc.WriteJSON(fields, ctx)
				}
				ctx.WriteString(`]}`)
			} else {
				ctx.WriteString(`{"accounts":[]}`)
			}
			return
		}
		path, id, ok := parse.ParseInt(path[10:])
		if !ok || !accountsSvc.Exists(id) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		switch string(path) {
		case `/`:
			if accountsSvc.Update(id, ctx.PostBody()) {
				ctx.SetStatusCode(fasthttp.StatusAccepted)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		case `suggetst/`:
		case `recommend/`:
		}
	}

	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}

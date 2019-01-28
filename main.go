package main

import (
	"log"
	"runtime"
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
	log.Println("NumCPU", runtime.NumCPU())
	log.Println("Restoring")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Creating Service", rep.Len())
	accountsSvc := service.NewAccountsService(rep)
	log.Println("Rebuild indexes")
	accountsSvc.RebuildIndexes()
	log.Println("Rebuild done")
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
				log.Println("Rebuild done")
			}
			time.Sleep(time.Second)
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
		case `/accounts/group/`:
			args := ctx.QueryArgs()
			if args.Len() < 4 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			limit, err := args.GetUint(`limit`)
			if err != nil || limit > 50 {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			order := args.Peek(`order`)
			if string(order) != `-1` && string(order) != `1` {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
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
			birthLT := 0
			birthGT := 0
			emailLT := []byte{}
			emailGT := []byte{}
			snameStarts := []byte{}
			fields := proto.IDField | proto.EmailField
			args.VisitAll(func(k, v []byte) {
				var next iterator.Iterator
				switch string(k) {
				case `limit`:
					return
				case `query_id`:
					return
				case `sex_eq`:
					if country := args.Peek(`country_eq`); len(country) > 0 {
						if next = accountsSvc.ByCountryEqSexEq(country, v); next == nil {
							hasErrors = true
							return
						}
					} else if next = accountsSvc.BySexEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.SexField
				case `status_eq`:
					if country := args.Peek(`country_eq`); len(country) > 0 {
						if next = accountsSvc.ByCountryEqStatusEq(country, v); next == nil {
							hasErrors = true
							return
						}
					} else if next = accountsSvc.ByStatusEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				case `status_neq`:
					if country := args.Peek(`country_eq`); len(country) > 0 {
						if next = accountsSvc.ByCountryEqStatusNeq(country, v); next == nil {
							hasErrors = true
							return
						}
					} else if next = accountsSvc.ByStatusNeq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.StatusField
				case `fname_eq`:
					if next = accountsSvc.ByFnameEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.FnameField
				case `fname_null`:
					if next = accountsSvc.ByFnameNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.FnameField
				case `fname_any`:
					if next = accountsSvc.ByFnameAny(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.FnameField
				case `sname_eq`:
					if next = accountsSvc.BySnameEq(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.SnameField
				case `sname_starts`:
					snameStarts = v
					fields |= proto.SnameField
					return
				case `sname_null`:
					if next = accountsSvc.BySnameNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.SnameField
				case `country_eq`:
					if args.Has(`sex_eq`) || args.Has(`status_eq`) || args.Has(`status_neq`) {
						fields |= proto.CountryField
						return
					} else if next = accountsSvc.ByCountryEq(v); next == nil {
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
				case `birth_year`:
					if next = accountsSvc.ByBirthYear(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.BirthField
				case `birth_lt`:
					if _, ts, ok := parse.ParseInt(v); !ok {
						hasErrors = true
						return
					} else {
						birthLT = ts
					}
					fields |= proto.BirthField
					return
				case `birth_gt`:
					if _, ts, ok := parse.ParseInt(v); !ok {
						hasErrors = true
						return
					} else {
						birthGT = ts
					}
					fields |= proto.BirthField
					return
				case `premium_now`:
					if next = accountsSvc.ByPremiumNow(); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.PremiumField
				case `premium_null`:
					if next = accountsSvc.ByPremiumNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.PremiumField
				case `phone_null`:
					if next = accountsSvc.ByPhoneNull(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.PhoneField
				case `phone_code`:
					if next = accountsSvc.ByPhoneCode(v); next == nil {
						hasErrors = true
						return
					}
					fields |= proto.PhoneField
				case `email_lt`:
					emailLT = v
					return
				case `email_gt`:
					emailGT = v
					return
				case `email_domain`:
					if next = accountsSvc.ByEmailDomain(v); next == nil {
						hasErrors = true
						return
					}
				case `likes_contains`:
					if next = accountsSvc.ByLikesContains(v); next == nil {
						hasErrors = true
						return
					}
				default:
					hasErrors = true
					return
				}
				if iter == nil {
					iter = next
					return
				}
				iter = iterator.NewInterIter(iter, next)
			})
			if hasErrors {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if iter == nil {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
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
				limit--
				if comma {
					ctx.WriteString(`,`)
				}
				comma = true
				acc.WriteJSON(fields, ctx)
			}
			ctx.WriteString(`]}`)
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

package main

import (
	//"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/service"
)

const (
	httpBasePath   = `/accounts/`
	httpBaseLen    = len(httpBasePath)
	httpNewPath    = `/accounts/new/`
	httpLikesPath  = `/accounts/likes/`
	httpFilterPath = `/accounts/filter/`
	httpGroupPath  = `/accounts/group/`
)

func AccountsHandler(ctx *fasthttp.RequestCtx, svc *service.AccountsService) {
	var id int
	var ok bool

	path := ctx.Path()
	switch string(path) {
	case httpFilterPath:
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
		correct := true
		fields := proto.IDField
		args.VisitAll(func(k, v []byte) {
			switch string(k) {
			case `sex_eq`:
				fields |= proto.SexField
			case `status_eq`:
				fields |= proto.StatusField
			case `status_neq`:
				fields |= proto.StatusField
			case `country_eq`:
				fields |= proto.CountryField
			case `country_null`:
				fields |= proto.CountryField
			case `city_eq`:
				fields |= proto.CityField
			case `city_null`:
				fields |= proto.CityField
			case `city_any`:
				fields |= proto.CityField
			case `interests_contains`:
			case `interests_any`:
			case `limit`:
			case `query_id`:
			default:
				correct = false
				return
			}
		})
		if !correct {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	case httpGroupPath:
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
	case httpNewPath:
		acc := &proto.Account{}
		if _, ok = acc.UnmarshalJSON(ctx.PostBody()); !ok {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if _, id, ok = parse.ParseInt(acc.ID[:]); !ok {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if !svc.Create(id, acc) {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	case httpLikesPath:
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		return
	default:
		if len(path) < httpBaseLen || string(path[:httpBaseLen]) != httpBasePath {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if path, id, ok = parse.ParseInt(path[httpBaseLen:]); !ok {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if ctx.IsPost() {
			if !svc.Exists(id) {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
			acc := &proto.Account{}
			if _, ok = acc.UnmarshalJSON(ctx.PostBody()); !ok {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if !svc.Update(id, acc) {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusAccepted)
			return
		}
	}
}

func main() {
	log.Println("Restore service")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts", rep.Len())
	svc := service.NewAccountsService(rep)
	handler := func(ctx *fasthttp.RequestCtx) {
		AccountsHandler(ctx, svc)
	}
	svc.AddInvertedIndex(inverted.DefaultParts, inverted.InterestsTokens)
	svc.AddInvertedIndex(inverted.DefaultParts, inverted.FnameTokens)
	svc.AddInvertedIndex(inverted.DefaultParts, inverted.SnameTokens)
	svc.AddInvertedIndex(inverted.DefaultParts, inverted.CountryTokens)
	svc.AddInvertedIndex(inverted.DefaultParts, inverted.CityTokens)
	svc.RebuildIndexes()
	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", handler); err != nil {
		log.Fatal(err)
	}
}

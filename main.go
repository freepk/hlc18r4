package main

import (
	"log"

	"github.com/freepk/hlc18r4/accounts"
	"github.com/freepk/hlc18r4/backup"
	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/search"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/iterator"
	"github.com/freepk/parse"
	"github.com/valyala/fasthttp"
)

var (
	writesCount uint64
	accountsSvc *accounts.AccountsService
	searchSvc   *search.SearchService
)

func routerHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	switch string(path) {
	case `/test0`:
		test0Handler(ctx)
	case `/test1`:
		test1Handler(ctx)
	case `/accounts/filter/`:
		filterHandler(ctx)
	case `/accounts/likes/`:
		likesHandler(ctx)
	case `/accounts/group/`:
		groupHandler(ctx)
	case `/accounts/new/`:
		createHandler(ctx)
	default:
		if path, id, ok := parse.ParseInt(path[10:]); ok {
			switch string(path) {
			case `/`:
				updateHandler(id, ctx)
			case `suggetst/`:
				suggestHandler(id, ctx)
			case `recommend/`:
				recommendHandler(id, ctx)
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
}

func likesHandler(ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.AddLikes(ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func createHandler(ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.Create(ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func updateHandler(id int, ctx *fasthttp.RequestCtx) {
	if err := accountsSvc.Update(id, ctx.PostBody()); err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	} else if err == accounts.NotFoundError {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func filterHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	limit, err := args.GetUint(`limit`)
	if err != nil || limit > 50 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	args.VisitAll(func(k, v []byte) {
	})
}

func test0Handler(ctx *fasthttp.RequestCtx) {
	countryKey, ok := tokens.Country([]byte(`Испезия`))
	if !ok {
		log.Fatal("Invalid country")
	}
	country := searchSvc.Countries(countryKey)
	if country == nil {
		log.Fatal("Country index is null")
	}
	interestKey, ok := tokens.Interest([]byte(`Обнимашки`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter := iterator.Iterator(country.Interest(interestKey))
	interestKey, ok = tokens.Interest([]byte(`YouTube`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter = iterator.NewUnionIter(iter, country.Interest(interestKey))
	interestKey, ok = tokens.Interest([]byte(`Солнце`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter = iterator.NewUnionIter(iter, country.Interest(interestKey))
	acc := &proto.Account{}
	limit := 22
	for {
		if limit == 0 {
			break
		}
		id, ok := iter.Next()
		if !ok {
			break
		}
		*acc = *accountsSvc.Get(2000000 - id)
		if acc.Sex != tokens.MaleSex {
			continue
		}
		acc.WriteJSON((proto.IDField | proto.EmailField | proto.SexField | proto.CountryField), ctx)
		limit--
	}

}

func test1Handler(ctx *fasthttp.RequestCtx) {
	interestKey, ok := tokens.Interest([]byte(`South Park`))
	if !ok {
		log.Fatal("Invalid interest")
	}
	iter := iterator.Iterator(searchSvc.Common().Interest(interestKey))
	acc := &proto.Account{}
	limit := 24
	for {
		if limit == 0 {
			break
		}
		id, ok := iter.Next()
		if !ok {
			break
		}
		*acc = *accountsSvc.Get(2000000 - id)
		if acc.Status != tokens.ComplStatus {
			continue
		}
		acc.WriteJSON((proto.IDField | proto.EmailField | proto.StatusField), ctx)
		limit--
	}
}

func groupHandler(ctx *fasthttp.RequestCtx) {
}

func suggestHandler(id int, ctx *fasthttp.RequestCtx) {
}

func recommendHandler(id int, ctx *fasthttp.RequestCtx) {
}

func main() {
	log.Println("Restoring")
	rep, err := backup.Restore("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accounts service")
	accountsSvc = accounts.NewAccountsService(rep)

	log.Println("Search service")
	searchSvc = search.NewSearchService(rep)
	searchSvc.Rebuild()

	log.Println("Start listen")
	if err := fasthttp.ListenAndServe(":80", routerHandler); err != nil {
		log.Fatal(err)
	}
}

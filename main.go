package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/parse"
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

	handler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
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
				return
			} else if args.Len() > 3 {
				if args.Has(`country_eq`) {
					country := args.Peek(`country_eq`)
					fmt.Println("part by country", string(country))
					return
				} else if args.Has(`city_eq`) {
					city := args.Peek(`city_eq`)
					fmt.Println("part by city", string(city), "in country")
					return
				} else if args.GetBool(`country_null`) {
					fmt.Println("part by country NULL")
					return
				} else if args.Has(`status_eq`) {
					status := args.Peek(`status_eq`)
					fmt.Println("part by status", string(status))
					return
				} else if args.Has(`sex_eq`) {
					sex := args.Peek(`sex_eq`)
					fmt.Println("part by sex", string(sex))
					return
				}
			}
			fmt.Println("default", string(args.QueryString()))
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

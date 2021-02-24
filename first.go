package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mcaci/login/db"
	"github.com/mcaci/login/route"
)

func firstTest(database *db.Database) {
	log.Println(database.Client.Set(db.Ctx, "language", "Go", 0))
	log.Println(database.Client.Get(db.Ctx, "year"))
	log.Println(database.Client.Get(db.Ctx, "language"))
	log.Println(database.Client.Del(db.Ctx, "language"))
}

func firstRouterTest(database *db.Database) *gin.Engine {
	r := gin.Default()
	route.Apply(
		route.NewDescr(
			r.POST,
			"/set/:lang",
			func(c *gin.Context) { log.Print(database.Client.Set(db.Ctx, "language", c.Param("lang"), 0)) },
		),
		route.NewDescr(
			r.POST,
			"/setgo",
			func(c *gin.Context) { log.Print(database.Client.Set(db.Ctx, "language", "go", 0)) },
		),
		route.NewDescr(
			r.POST,
			"/del",
			func(c *gin.Context) { log.Print(database.Client.Del(db.Ctx, "language")) },
		),
		route.NewDescr(
			r.GET,
			"/get",
			func(c *gin.Context) { log.Print(database.Client.Get(db.Ctx, "language")) },
		),
	)
	return r
}

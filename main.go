package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mcaci/login/db"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
)

func main() {
	database, err := db.NewDatabase(RedisAddr)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	router := initRouter(database)
	router.Run(ListenAddr)
	// log.Println(database.Client.Set(db.Ctx, "language", "Go", 0))
	// log.Println(database.Client.Get(db.Ctx, "year"))
	// log.Println(database.Client.Get(db.Ctx, "language"))
	// log.Println(database.Client.Del(db.Ctx, "language"))

	// pipe := database.Client.TxPipeline()
	// pipe.Set(db.Ctx, "language", "golang", 0)
	// pipe.Set(db.Ctx, "year", 2009, 0)
	// results, err := pipe.Exec(db.Ctx)
	// log.Println(results)
}

func initRouter(database *db.Database) *gin.Engine {
	r := gin.Default()
	r.POST("/set/:lang", func(c *gin.Context) {
		lang := c.Param("lang")
		log.Print(database.Client.Set(db.Ctx, "language", lang, 0))
	})
	r.POST("/setgo", func(c *gin.Context) {
		log.Print(database.Client.Set(db.Ctx, "language", "go", 0))
	})
	r.GET("/get", func(c *gin.Context) {
		log.Print(database.Client.Get(db.Ctx, "language"))
	})
	r.POST("/del", func(c *gin.Context) {
		log.Print(database.Client.Del(db.Ctx, "language"))
	})
	return r
}

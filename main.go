package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcaci/login/db"
	"github.com/mcaci/login/route"
)

func main() {
	cliURL := flag.String("url", "localhost", "URL of the login server. Default: localhost.")
	cliPort := flag.String("port", "8080", "Port of the login server. Default: 8080.")
	redisURL := flag.String("redis-url", "localhost", "URL of redis server. Default: localhost.")
	redisPort := flag.String("redis-port", "6379", "Port of redis server. Default: 6379.")
	flag.Parse()

	database, err := db.NewDatabase(fmt.Sprintf("%s:%s", *redisURL, *redisPort))
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	router := gin.Default()
	route.Apply(
		route.NewDescr(router.POST, "/register", route.Handle(database, route.WithRegisterHandler)),
		route.NewDescr(router.POST, "/login", route.Handle(database, route.WithLoginHandler)),
		route.NewDescr(router.GET, "/welcome", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": "welcome"}) }),
	)
	router.Run(fmt.Sprintf("%s:%s", *cliURL, *cliPort))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"unicode"

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

	name := "hello"
	name = strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return unicode.ToLower(r)
		case r >= 'a' && r <= 'z',
			unicode.IsDigit(r),
			r == '.',
			r == '-':
			return r
		default:
			return '-' // or 0 if you want to replace with 'empty' char
		}
	}, name)
	name = strings.TrimFunc(name, func(r rune) bool {
		return !unicode.IsLower(r) && !unicode.IsNumber(r)
	})

	database, err := db.NewDatabase(fmt.Sprintf("%s:%s", *redisURL, *redisPort))
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	router := gin.Default()
	route.Apply(
		route.NewDescr(router.POST, "/register", route.Handle(database, route.WithRegisterHandler)),
		route.NewDescr(router.POST, "/login", route.Handle(database, route.WithLoginHandler)),
	)
	router.Run(fmt.Sprintf("%s:%s", *cliURL, *cliPort))
}

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
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
		route.NewDescr(router.GET, "/welcome", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": "a public welcome"}) }),
		route.NewDescr(router.GET, "/private-welcome", func(c *gin.Context) {
			jwtMdw.HandlerWithNext(c.Writer, c.Request, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
				// token := authHeaderParts[1]

				// hasScope := checkScope("read:private-welcome", token)

				// if !hasScope {
				// 	c.JSON(http.StatusForbidden, gin.H{"nok": "Insufficient scope."})
				// 	return
				// }
				// message := "Hello from a private endpoint! You need to be authenticated to see this."
				c.JSON(http.StatusOK, gin.H{"ok": "as authenticated user you get a private welcome"})
			}))
		}),
	)
	log.Println(router.Run(fmt.Sprintf("%s:%s", *cliURL, *cliPort)))
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwtMdw = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "localhost" // Auth0 API Identifier
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		// Verify 'iss' claim
		iss := "https://dev-mcmg1609-m.eu.auth0.com/" // Auth0 Cluster
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-mcmg1609-m.eu.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*CustomClaims)

	hasScope := false
	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}
	return hasScope
}

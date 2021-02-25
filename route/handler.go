package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcaci/login/db"
	"github.com/mcaci/login/usr"
)

type handler func(func(interface{}) error, *db.Database) (int, interface{})

func Handle(database *db.Database, h handler) func(*gin.Context) {
	return func(c *gin.Context) { c.JSON(h(c.ShouldBindJSON, database)) }
}

func WithRegisterHandler(bind func(interface{}) error, database *db.Database) (int, interface{}) {
	var userJSON usr.User
	if err := bind(&userJSON); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	}
	err := usr.Save(database, &userJSON)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	}
	return http.StatusOK, gin.H{"user": userJSON.Username}
}

func WithLoginHandler(bind func(interface{}) error, database *db.Database) (int, interface{}) {
	var userJSON usr.User
	if err := bind(&userJSON); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	}
	user, err := usr.Get(database, userJSON.Username)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	}
	if user.Password != userJSON.Password {
		return http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"}
	}
	return http.StatusOK, gin.H{"logged as": userJSON.Username}
}

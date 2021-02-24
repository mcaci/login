package route

import (
	"github.com/gin-gonic/gin"
)

type Descr struct {
	httpF        func(string, ...gin.HandlerFunc) gin.IRoutes
	relativePath string
	handlerF     gin.HandlerFunc
}

func NewDescr(httpF func(string, ...gin.HandlerFunc) gin.IRoutes, relativePath string, handlerF gin.HandlerFunc) *Descr {
	return &Descr{
		httpF:        httpF,
		relativePath: relativePath,
		handlerF:     handlerF,
	}
}

func Apply(routes ...*Descr) {
	for _, route := range routes {
		route.httpF(route.relativePath, route.handlerF)
	}
}

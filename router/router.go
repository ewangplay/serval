package router

import (
	apiV1 "github.com/ewangplay/serval/api/v1"
	mw "github.com/ewangplay/serval/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the HTTP router
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(mw.Adapter())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", apiV1.Pong)

		v1.POST("/did/create", apiV1.CreateDid)
		v1.GET("/did/resolve/:did", apiV1.ResolveDid)
	}

	return r
}

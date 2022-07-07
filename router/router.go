package router

import (
	"io"

	apiV1 "github.com/ewangplay/serval/api/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the HTTP router
func InitRouter(w io.Writer) *gin.Engine {
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(w))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", apiV1.Pong)

		v1.POST("/did/create", apiV1.CreateDid)
		v1.GET("/did/resolve/:did", apiV1.ResolveDid)
		v1.POST("/did/revoke", apiV1.RevokeDid)
	}

	return r
}

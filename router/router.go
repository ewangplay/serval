package router

import (
	"io"
	"net/http"

	cl "github.com/ewangplay/cryptolib"
	apiV1 "github.com/ewangplay/serval/api/v1"
	ctx "github.com/ewangplay/serval/context"
	"github.com/gin-gonic/gin"
	"github.com/jerray/qsign"
	"github.com/philippgille/gokv"
)

// InitRouter initializes the HTTP router
func InitRouter(w io.Writer, store gokv.Store, csp cl.CSP, qsign *qsign.Qsign) *gin.Engine {
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(w))
	r.Use(initContext(store, csp, qsign))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", apiV1.Pong)

		v1.POST("/did/create", convert(apiV1.CreateDid))
		v1.GET("/did/resolve/:did", convert(apiV1.ResolveDid))
		v1.POST("/did/revoke", convert(apiV1.RevokeDid))
	}

	return r
}

type handlerFunc func(*ctx.Context)

func initContext(store gokv.Store, csp cl.CSP, qsign *qsign.Qsign) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := &ctx.Context{
			Context: c,
			Store:   store,
			CSP:     csp,
			Qsign:   qsign,
		}
		c.Set("context", context)

		// process handler
		c.Next()
	}
}

func convert(f handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("context")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx, ok := v.(*ctx.Context)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		f(ctx)
	}
}

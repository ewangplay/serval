package v1

import (
	"github.com/gin-gonic/gin"
)

// Pong handles the /api/v1/ping request
func Pong(c *gin.Context) {
	OkWithData(gin.H{
		"message": "pong",
	}, c)
}

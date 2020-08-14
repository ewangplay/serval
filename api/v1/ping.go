package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pong handles the /api/v1/ping request
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

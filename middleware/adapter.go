package middleware

import (
	ch "github.com/ewangplay/serval/adapter/cryptohub"
	"github.com/gin-gonic/gin"
)

// Adapter adds some adapters into gin Context
func Adapter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
		c.Next()
	}
}

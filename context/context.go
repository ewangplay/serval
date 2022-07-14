package context

import (
	cl "github.com/ewangplay/cryptolib"
	"github.com/gin-gonic/gin"
	"github.com/jerray/qsign"
	"github.com/philippgille/gokv"
)

type Context struct {
	*gin.Context
	Store gokv.Store
	CSP   cl.CSP
	Qsign *qsign.Qsign
}

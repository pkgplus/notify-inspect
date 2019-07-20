package app

import (
	"github.com/gin-gonic/gin"
)

var (
	app = gin.New()
)

func init() {
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
}

func GetApp() *gin.Engine {
	return app
}

package main

import (
	"os"

	"github.com/kataras/iris"
	"github.com/xuebing1110/notify-inspect/server/app"

	_ "github.com/xuebing1110/notify-inspect/pkg/plugin/storage/redis"
	_ "github.com/xuebing1110/notify-inspect/server/router/v2"
)

func main() {
	// http server
	irisApp := app.GetIrisApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	irisApp.Run(iris.Addr(":" + port))
}

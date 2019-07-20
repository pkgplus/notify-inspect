package main

import (
	"os"

	"github.com/xuebing1110/notify-inspect/client/app"

	_ "github.com/xuebing1110/notify-inspect/client/router/v1"
)

func main() {
	// http server

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	app.GetApp().Run(":" + port)
}

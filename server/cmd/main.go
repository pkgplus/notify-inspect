package main

import (
	"github.com/xuebing1110/notify-inspect/pkg/schedule"
	_ "github.com/xuebing1110/notify-inspect/pkg/schedule/redis"
	"github.com/xuebing1110/notify-inspect/server/app"

	_ "github.com/xuebing1110/notify-inspect/pkg/plugin/storage/redis"
	_ "github.com/xuebing1110/notify-inspect/server/router/v2"
)

func main() {
	err := schedule.Start()
	if err != nil {
		panic(err)
	}

	// http server
	app.GetApp().Run()
}

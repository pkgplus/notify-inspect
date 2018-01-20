package v2

import (
	"github.com/kataras/iris/websocket"
	"github.com/xuebing1110/notify-inspect/server/app"
	"github.com/xuebing1110/notify-inspect/server/handlers"
)

func init() {
	api := app.GetIrisApp().Party(API_ROUTE_PREFIX)

	// plugin register
	ws := websocket.New(websocket.Config{})
	ws.OnConnection(handlers.RegistePlugin)
	api.Get("/register", ws.Handler())

	// list all support plugins
	api.Get("/", handlers.ListPlugins)
	api.Get("/:pid", handlers.GetPlugin)

}

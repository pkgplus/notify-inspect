package v1

import (
	"github.com/kataras/iris/websocket"
	"github.com/xuebing1110/notify-inspect/server/handlers"
)

func plugin() {
	// plugin register
	ws := websocket.New(websocket.Config{})
	ws.OnConnection(handlers.RegistePlugin)
	api.Get("/plugins/register", ws.Handler())

	// list all support plugins
	api.Get("/plugins", handlers.ListPlugins)
	api.Get("/plugins/:pid", handlers.GetPlugin)

	// plugin api: /plugins
	plugin_api := api.Party("/plugins")
	plugin_api.Use(handlers.User)

	// user subscribe plugin
	plugin_api.Post("/:pid/sub", handlers.SavePluginSubscribe)
	plugin_api.Get("/:pid/sub", handlers.GetPluginSubscribe)
	plugin_api.Delete("/:pid/sub", handlers.DeletePluginSubscribe)

	// user plugin records
	plugin_api.Get("/:pid/records", handlers.ListPluginRecords)
	plugin_api.Post("/:pid/records", handlers.AddPluginRecord)
	plugin_api.Get("/:pid/records/:rid", handlers.GetPluginRecord)
	plugin_api.Put("/:pid/records/:rid", handlers.ModifyPluginRecord)
	plugin_api.Delete("/:pid/records/:rid", handlers.AddPluginRecord)
}

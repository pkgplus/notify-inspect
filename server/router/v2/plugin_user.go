package v2

import (
	"github.com/xuebing1110/notify-inspect/server/app"
	"github.com/xuebing1110/notify-inspect/server/handlers"
)

func init() {
	api := app.GetIrisApp().Party(API_ROUTE_PREFIX)
	api.Use(handlers.User)

	// user subscribe plugin
	api.Post("/:pid/sub", handlers.SavePluginSubscribe)
	api.Get("/:pid/sub", handlers.GetPluginSubscribe)
	api.Delete("/:pid/sub", handlers.DeletePluginSubscribe)

	// user plugin records
	api.Get("/:pid/records", handlers.ListPluginRecords)
	api.Post("/:pid/records", handlers.AddPluginRecord)
	api.Get("/:pid/records/:rid", handlers.GetPluginRecord)
	api.Put("/:pid/records/:rid", handlers.ModifyPluginRecord)
	api.Delete("/:pid/records/:rid", handlers.DeletePluginRecord)
}

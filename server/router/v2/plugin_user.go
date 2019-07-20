package v2

import (
	"github.com/xuebing1110/notify-inspect/server/handlers"
)

func init() {
	api_user := api.Group("/plugins", handlers.User)

	// user subscribe plugin
	api_user.POST("/:pid/sub", handlers.SavePluginSubscribe, handlers.CallPluginSubscribe)
	api_user.GET("/:pid/sub", handlers.GetPluginSubscribe, handlers.CallPluginSubscribeStatus)
	api_user.DELETE("/:pid/sub", handlers.DeletePluginSubscribe)

	// user plugin records
	api_user.GET("/:pid/records", handlers.ListPluginRecords)
	api_user.POST("/:pid/records", handlers.AddPluginRecord)
	api_user.GET("/:pid/records/:rid", handlers.GetPluginRecord)
	api_user.PUT("/:pid/records/:rid", handlers.ModifyPluginRecord)
	api_user.DELETE("/:pid/records/:rid", handlers.DeletePluginRecord)
}

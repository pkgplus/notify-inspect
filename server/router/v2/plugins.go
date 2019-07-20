package v2

import (
	"github.com/xuebing1110/notify-inspect/server/handlers"
)

func init() {
	// plugin register
	api.GET("/plugin_register", handlers.RegisterPlugin)

	// list all support plugins
	api.GET("/plugins", handlers.ListPlugins)
	api.GET("/plugins/:pid", handlers.GetPlugin)

}

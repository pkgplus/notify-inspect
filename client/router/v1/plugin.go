package v1

import (
	"github.com/xuebing1110/notify-inspect/client/handlers"
)

func plugin() {
	api.POST("/sub/users", handlers.Subscribe)
	api.POST("/sub/records", handlers.RecordNotice)
}

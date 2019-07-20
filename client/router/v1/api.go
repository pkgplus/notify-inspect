package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xuebing1110/notify-inspect/client/app"
)

var api *gin.RouterGroup

func init() {
	api = app.GetApp().Group("/api/v1/plugins/" + app.APP_NAME)

	plugin()
}

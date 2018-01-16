package v1

import (
	"github.com/kataras/iris/core/router"
	"github.com/xuebing1110/notify-inspect/client/app"
)

var api router.Party

func init() {
	irisApp := app.GetIrisApp()
	api = irisApp.Party("/api/v1/plugins/" + app.APP_NAME)

	plugin()
}

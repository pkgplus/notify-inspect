package v2

import "github.com/xuebing1110/notify-inspect/server/app"

const (
	API_ROUTE_PREFIX = "/api/v2/notify"
)

var (
	api = app.GetApp().Group(API_ROUTE_PREFIX)
)

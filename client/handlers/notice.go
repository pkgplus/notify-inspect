package handlers

import (
	"github.com/kataras/iris/context"
)

func RecordNotice(ctx context.Context) {
	SendNormalResponse(ctx, []string{})
}

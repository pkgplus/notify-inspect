package handlers

import (
	"net/http"

	"github.com/kataras/iris/context"
)

func User(ctx context.Context) {
	uid := ctx.GetHeader("X-Uid")
	if uid == "" {
		SendResponse(ctx, http.StatusBadRequest, "unknownUser", "the user must be specified")
		return
	}
	ctx.Values().Set(CTX_USERID, uid)
	ctx.Next()
}

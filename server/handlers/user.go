package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const CTX_USERID = "X-USER"

func User(ctx *gin.Context) {
	uid := ctx.GetHeader("X-USER")
	if uid == "" {
		SendResponse(ctx, http.StatusBadRequest, "unknownUser", "the user must be specified")
		return
	}
	ctx.Set(CTX_USERID, uid)
	ctx.Next()
}

package handlers

import (
	"github.com/gin-gonic/gin"
)

func RecordNotice(ctx *gin.Context) {
	SendNormalResponse(ctx, []string{})
}

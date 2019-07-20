package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Subscribe(ctx *gin.Context) {
	SendResponse(ctx, http.StatusOK, "OK", "")
}

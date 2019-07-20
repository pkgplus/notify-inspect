package handlers

import (
	"github.com/gin-gonic/gin"
)

func SendResponse(ctx *gin.Context, code int, msg, detail string) {
	resp := &Response{
		Code:    code,
		Message: msg,
		Detail:  detail,
	}
	ctx.JSON(resp.Code, resp)
	ctx.Abort()
}

func SendNormalResponse(ctx *gin.Context, data interface{}) {
	resp := &Response{
		Code:    200,
		Message: "OK",
		Detail:  "",
		Data:    data,
	}
	ctx.JSON(resp.Code, resp)
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data,omitempty"`
}

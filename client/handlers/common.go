package handlers

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data,omitempty"`
}

func SendResponse(ctx *gin.Context, code int, msg, detail string) {
	resp := &Response{
		Code:    code,
		Message: msg,
		Detail:  detail,
	}
	ctx.JSON(resp.Code, resp)
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

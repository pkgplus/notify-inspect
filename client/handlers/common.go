package handlers

import (
	"github.com/kataras/iris/context"
	"github.com/xuebing1110/notify/pkg/models"
)

func SendResponse(ctx context.Context, code int, msg, detail string) {
	resp := &models.Response{
		Code:    code,
		Message: msg,
		Detail:  detail,
	}
	ctx.StatusCode(resp.Code)
	ctx.JSON(resp)
}

func SendNormalResponse(ctx context.Context, data interface{}) {
	resp := &models.Response{
		Code:    200,
		Message: "OK",
		Detail:  "",
		Data:    data,
	}
	ctx.StatusCode(resp.Code)
	ctx.JSON(resp)
}

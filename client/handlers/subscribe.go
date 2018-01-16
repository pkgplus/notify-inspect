package handlers

import (
	"github.com/kataras/iris/context"
	"net/http"
)

func Subscribe(ctx context.Context) {
	SendResponse(ctx, http.StatusOK, "OK", "")
}

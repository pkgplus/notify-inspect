package app

import (
	"github.com/kataras/iris"
)

var (
	IrisApp *iris.Application
)

func init() {
	IrisApp = iris.New()
}

func GetIrisApp() *iris.Application {
	return IrisApp
}

package app

import (
	"os"

	"github.com/kataras/iris"
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/client"
)

var (
	APP_NAME string
	IrisApp  *iris.Application
)

func init() {
	// app name
	APP_NAME = os.Getenv("APP_NAME")
	if APP_NAME == "" {
		panic("the env \"APP_NAME\" is required!")
	}

	// register plugin
	err := register()
	if err != nil {
		panic("register plugin failed: " + err.Error())
	}

	IrisApp = iris.New()
}

func register() error {
	c := client.DefaultRegisterClient
	p := &plugin.Plugin{
		Id:            APP_NAME,
		Description:   "",
		ServeAddr:     "http://127.0.0.1:8081/api/v1/plugins",
		TemplateMsgId: "8U98v1g7PWLZ5p4jbWNSpY5dr-hhG5kVuMAUew4PHnY",
		Params: []plugin.PluginParam{
			plugin.PluginParam{
				Id:         "userid",
				Name:       "工号",
				Value:      "",
				Candidates: []plugin.PluginData{},
			},
		},
		RecordParams: []plugin.PluginParam{},
		Author:       "bingbaba.com",
	}

	return c.Register(p)
}

func GetIrisApp() *iris.Application {
	return IrisApp
}

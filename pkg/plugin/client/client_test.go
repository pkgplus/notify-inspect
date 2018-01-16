package client

import (
	"testing"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

func TestClientRegister(t *testing.T) {
	p := &plugin.Plugin{
		Id:            "test",
		Description:   "the test plugin",
		ServeAddr:     "http://127.0.0.1:8080",
		TemplateMsgId: "1111",
		RecordParams: []plugin.PluginParam{
			plugin.PluginParam{
				Id:         "key1",
				Name:       "字段1",
				Value:      "value1",
				Candidates: []plugin.PluginData{},
			},
		},
		Author: "zhangsan",
	}
	err := DefaultRegisterClient.Register(p)
	if err != nil {
		t.Fatal(err)
	}
}

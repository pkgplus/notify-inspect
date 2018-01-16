package plugin

import (
	"testing"
)

func TestClientRegister(t *testing.T) {
	p := &Plugin{
		Id:            "test",
		Description:   "the test plugin",
		ServeAddr:     "http://127.0.0.1:8080",
		TemplateMsgId: "1111",
		RecordParams: []PluginParam{
			PluginParam{
				Id:         "key1",
				Name:       "字段1",
				Value:      "value1",
				Candidates: []PluginData{},
			},
		},
		Author: "zhangsan",
	}
	err := DefaultRegisterClient.Register(p)
	if err != nil {
		t.Fatal(err)
	}
}

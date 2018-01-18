package plugin

import (
	"encoding/json"
	"errors"
	"time"
)

type Plugin struct {
	Id              string        `json:"id"`
	Description     string        `json:"description"`
	ServeAddr       string        `json:"serveAddr"`
	TemplateMsgId   string        `json:"templateMsgId"`
	Emphasis        string        `json:"emphasis"`
	Page            string        `json:"page"`
	Params          []PluginParam `json:"params"`
	RecordParams    []PluginParam `json:"recordParams"`
	Author          string        `json:"author"`
	RegisterTime    int64         `json:"registerTime,omitempty"`
	RegisterTimeStr string        `json:"registerTimeStr,omitempty"`
	LostTime        int64         `json:"lostTime,omitempty"`
	LostTimeStr     string        `json:"lostTimeStr,omitempty"`
}

func (p *Plugin) SetRegisterTime() {
	p.RegisterTime = time.Now().Unix()
	p.RegisterTimeStr = time.Unix(p.RegisterTime, 0).Format("2006-01-02 15:04:05")
}

func (p *Plugin) SetLost() {
	p.LostTime = time.Now().Unix()
	p.LostTimeStr = time.Unix(p.LostTime, 0).Format("2006-01-02 15:04:05")
}

func (p *Plugin) ResetLost() {
	p.LostTime = 0
}

func (p *Plugin) ToJson() []byte {
	body, _ := json.Marshal(p)
	return body
}

func NewPlugin(data []byte) (*Plugin, error) {
	p := new(Plugin)
	err := json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	if p.Id == "" {
		return nil, errors.New("the id must be specified")
	}
	if p.ServeAddr == "" {
		return nil, errors.New("the serveAddr must be specified")
	}
	if p.TemplateMsgId == "" {
		return nil, errors.New("the templateMsgId must be specified")
	}
	// if len(p.RecordParams) == 0 {
	// 	return nil, errors.New("the recordParams must be specified")
	// }

	// set registe time
	// p.SetRegisterTime()
	// p.lostTime = 0

	return p, nil
}

func (p *Plugin) GetParamValue(id string) string {
	for _, param := range p.Params {
		if param.Id == id {
			return param.Value
		}
	}
	return ""
}

func (p *Plugin) GetRecordParamValue(id string) string {
	for _, param := range p.RecordParams {
		if param.Id == id {
			return param.Value
		}
	}
	return ""
}

type PluginParam struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`

	Candidates []PluginData `json:"candidates"`
}

type PluginData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/xuebing1110/notify-inspect/pkg/notice/models"
)

var (
	client *http.Client
)

func init() {
	client = http.DefaultClient
}

type Plugin struct {
	Id              string        `json:"id"`
	Description     string        `json:"description"`
	ServerAddr      string        `json:"serverAddr"`
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

func (p *Plugin) BackendGetSubscribe(ctx context.Context, s *Subscribe) *Response {
	req_url := fmt.Sprintf("%s/sub/users/%s", p.ServerAddr, s.UserId)

	// url parameter
	url_param := url.Values{}
	for _, param := range s.Data {
		url_param.Set(param.Id, param.Value)
	}

	// http request
	req_url = req_url + "?" + url_param.Encode()
	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	if resp.StatusCode >= 400 {
		log.Printf("call %s failed: %s", req_url, resp_body)
	}

	call_resp, err := NewResponse(resp_body)
	if err != nil {
		return &Response{
			Code:    resp.StatusCode,
			Message: "SubscribePluginFailed",
			Detail:  string(resp_body),
		}
	}

	return call_resp
}

func (p *Plugin) BackendSubscribe(ctx context.Context, s *Subscribe) *Response {

	req_url := fmt.Sprintf("%s/sub/users", p.ServerAddr)
	req, err := http.NewRequest(http.MethodPost, req_url, bytes.NewReader(s.ToJson()))
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	if resp.StatusCode >= 400 {
		log.Printf("call %s failed: %s", req_url, resp_body)
	}

	call_resp, err := NewResponse(resp_body)
	if err != nil {
		return &Response{
			Code:    resp.StatusCode,
			Message: "InternalServerError",
			Detail:  string(resp_body),
		}
	}

	return call_resp
}

func (p *Plugin) BackendInspect(ctx context.Context, r *PluginRecord) (*models.Notice, error) {
	req_url := fmt.Sprintf("%s/sub/records", p.ServerAddr)

	r.Cron = nil
	body := r.ToJson()
	req, err := http.NewRequest(http.MethodPost, req_url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		log.Printf("call %s failed: %s", req_url, resp_body)
		return nil, fmt.Errorf("call %s failed, statusCode: %d", req_url, resp.StatusCode)
	}

	iresp := &Response{}
	err = json.Unmarshal(resp_body, iresp)
	if err != nil {
		return nil, err
	}

	if len(iresp.Data) == 0 {
		return nil, nil
	} else {
		return &models.Notice{
			UserID:   r.UserId,
			Template: p.TemplateMsgId,
			Emphasis: p.Emphasis,
			Page:     p.Page,
			Values:   iresp.Data,
		}, nil
	}
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
	if p.ServerAddr == "" {
		return nil, errors.New("the serverAddr must be specified")
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

type PluginParam struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`

	Candidates []PluginData `json:"candidates"`
}

type PluginData struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

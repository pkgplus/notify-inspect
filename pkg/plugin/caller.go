package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/xuebing1110/notify-inspect/pkg/log"
	"github.com/xuebing1110/notify-inspect/pkg/models"
	wxmodels "github.com/xuebing1110/notify/pkg/wechat/models"
)

var (
	client *http.Client
)

func init() {
	client = http.DefaultClient
}

func (p *Plugin) BackendGetSubscribe(s *Subscribe) *models.Response {
	req_url := fmt.Sprintf("%s/sub/users/%s", p.ServeAddr, s.UserId)

	// url parameter
	url_param := url.Values{}
	for _, param := range s.Data {
		url_param.Set(param.Id, param.Value)
	}

	// http request
	req_url = req_url + "?" + url_param.Encode()
	resp, err := http.DefaultClient.Get(req_url)
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	if resp.StatusCode >= 400 {
		log.GlobalLogger.Errorf("call %s failed: %s", req_url, resp_body)
	}

	call_resp, err := models.NewResponse(resp_body)
	if err != nil {
		return &models.Response{
			Code:    resp.StatusCode,
			Message: "InternalServerError",
			Detail:  string(resp_body),
		}
	}

	return call_resp
}

func (p *Plugin) BackendSubscribe(s *Subscribe) *models.Response {

	req_url := fmt.Sprintf("%s/sub/users", p.ServeAddr)
	req, err := http.NewRequest(http.MethodPost, req_url, bytes.NewReader(s.ToJson()))
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
			Detail:  err.Error(),
		}
	}
	if resp.StatusCode >= 400 {
		log.GlobalLogger.Errorf("call %s failed: %s", req_url, resp_body)
	}

	call_resp, err := models.NewResponse(resp_body)
	if err != nil {
		return &models.Response{
			Code:    resp.StatusCode,
			Message: "InternalServerError",
			Detail:  string(resp_body),
		}
	}

	return call_resp
}

func (p *Plugin) BackendInspect(r *PluginRecord) (*wxmodels.Notice, error) {
	req_url := fmt.Sprintf("%s/sub/records", p.ServeAddr)

	r.Cron = nil
	body := r.ToJson()
	req, err := http.NewRequest(http.MethodPost, req_url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

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
		log.GlobalLogger.Errorf("call %s failed: %s", req_url, resp_body)
		return nil, fmt.Errorf("call %s failed, statusCode: %d", req_url, resp.StatusCode)
	}

	iresp := &models.InspectResponse{}
	err = json.Unmarshal(resp_body, iresp)
	if err != nil {
		return nil, err
	}

	if len(iresp.Data) == 0 {
		return nil, nil
	} else {
		return &wxmodels.Notice{
			UserID:   r.UserId,
			Template: p.TemplateMsgId,
			Emphasis: p.Emphasis,
			Page:     p.Page,
			Values:   iresp.Data,
		}, nil
	}
}

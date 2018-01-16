package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (p *Plugin) BackendSubscribe(s *Subscribe) error {
	url := fmt.Sprintf("%s/sub/users", p.ServeAddr)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(s.ToJson()))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		log.GlobalLogger.Errorf("call %s failed: %s", url, resp_body)
		return fmt.Errorf("call %s failed, statusCode: %d", url, resp.StatusCode)
	}
	return nil
}

func (p *Plugin) BackendInspect(r *PluginRecord) (*wxmodels.Notice, error) {
	url := fmt.Sprintf("%s/sub/records", p.ServeAddr)

	r.Cron = nil
	body := r.ToJson()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
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
		log.GlobalLogger.Errorf("call %s failed: %s", url, resp_body)
		return nil, fmt.Errorf("call %s failed, statusCode: %d", url, resp.StatusCode)
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

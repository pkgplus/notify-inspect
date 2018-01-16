package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/xuebing1110/notify/pkg/models"
	wxmodels "github.com/xuebing1110/notify/pkg/wechat/models"
)

const (
	HOST_NOTICE_URL_FMT = "https://m.bingbaba.com/api/v2/notify/users/%s/notice"
)

func SendNotice(n *wxmodels.Notice) error {
	url := fmt.Sprintf(HOST_NOTICE_URL_FMT, n.UserID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(n.ToJson()))
	if err != nil {
		return err
	}

	respObj := new(models.Response)
	err = httpDo(req, respObj)
	if err != nil {
		return err
	}

	if respObj.Code < 400 {
		return nil
	} else {
		return fmt.Errorf("%s, detail: %s", respObj.Message, respObj.Detail)
	}
}

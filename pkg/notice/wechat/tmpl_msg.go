package wechat

import (
	"fmt"
	"github.com/xuebing1110/notify-inspect/pkg/notice/models"
	"strings"
)

type TemplateMsg struct {
	ToUserID        string                          `json:"touser"`
	TemplateID      string                          `json:"template_id"`
	FormID          string                          `json:"form_id"`
	Page            string                          `json:"page,omitempty"`
	Data            map[string]TemplateMsgDataValue `json:"data"`
	EmphasisKeyword string                          `json:"emphasis_keyword,omitempty"`
}

type TemplateMsgDataValue struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func NewTemplateMsgData(values []string) map[string]TemplateMsgDataValue {
	data := make(map[string]TemplateMsgDataValue)
	for i, value := range values {
		data[fmt.Sprintf("keyword%d", i+1)] = TemplateMsgDataValue{Value: value}
	}

	return data
}

func NewTemplateMsg(userid, templateid, formid string, values []string) *TemplateMsg {
	msg := &TemplateMsg{
		ToUserID:   userid,
		TemplateID: templateid,
		FormID:     formid,
		Data:       NewTemplateMsgData(values),
	}

	// data, _ := json.Marshal(msg)
	// fmt.Printf("%s\n", data)

	return msg
}

func (tmsg *TemplateMsg) SetEmphasis(i string) {
	tmsg.EmphasisKeyword = fmt.Sprintf("keyword%s.DATA", i)
}
func (tmsg *TemplateMsg) SetPage(page string) {
	tmsg.Page = page
}

func NoticeToTemplateMsg(n *models.Notice) (*TemplateMsg, error) {
	if n.Page == "" {
		n.Page = "/pages/index/index"
	} else if strings.Index(n.Page, "/") < 0 {
		n.Page = fmt.Sprintf("/pages/%s/%s", n.Page, n.Page)
	}

	if n.Emphasis == "" {
		n.Emphasis = "1"
	}

	// TODO pop a formid
	energy := ""
	openid := n.UserID
	//energy, err := storage.GlobalStore.PopEnergy(n.UserID)
	//if err != nil {
	//	return nil, err
	//}

	tmsg := NewTemplateMsg(openid, n.Template, energy, n.Values)
	tmsg.SetEmphasis(n.Emphasis)
	tmsg.SetPage(n.Page)
	return tmsg, nil
}

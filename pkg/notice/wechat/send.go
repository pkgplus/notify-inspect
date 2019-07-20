package wechat

import (
	"context"
	"github.com/esap/wechat"
	"github.com/pkg/errors"
	"github.com/xuebing1110/notify-inspect/pkg/notice"
	"github.com/xuebing1110/notify-inspect/pkg/notice/models"
	"log"
	"os"
)

const (
	FMT_URL_JSCODE2SESSION = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func init() {
	appid := os.Getenv("WX_APPID")
	secret := os.Getenv("WX_SECRET")
	if appid == "" || secret == "" {
		panic("can't get WX_APPID and WX_SECRET from env")
	}

	server := wechat.New("", appid, secret, "")
	server.MsgUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="

	err := notice.SenderRegister(&TemplateSender{Server: server})
	if err != nil {
		panic(err)
	}
}

type TemplateSender struct {
	*wechat.Server
}

func (s *TemplateSender) Name() string {
	return "WechatTemplateMessage"
}
func (s *TemplateSender) Send(ctx context.Context, n *models.Notice) error {
	msg, err := NoticeToTemplateMsg(n)
	if err != nil {
		return errors.Wrap(err, "convert to template message failed")
	}

	log.Printf("send message %+v", msg)
	wx_err := s.SendMsg(msg)
	if wx_err != nil {
		return errors.Wrap(wx_err.Error(), "send wechat message failed")
	}

	return nil
}

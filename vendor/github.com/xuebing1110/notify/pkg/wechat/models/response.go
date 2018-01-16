package models

import (
	"errors"
)

type SessionResp struct {
	OpenID  string `json:"openid"`
	SessKey string `json:"session_key"`
	Unionid string `json:"unionid"`
}

func NewSessionResp(sessMap map[string]string) (*SessionResp, error) {
	openid := sessMap["openid"]
	session_key := sessMap["session_key"]
	unionid := sessMap["unionid"]

	if openid == "" || session_key == "" {
		return nil, errors.New("openid/session_key must not null")
	}

	return &SessionResp{openid, session_key, unionid}, nil
}

func (s *SessionResp) Convert2Map() map[string]interface{} {
	return map[string]interface{}{
		"openid":      s.OpenID,
		"session_key": s.SessKey,
		"unionid":     s.Unionid,
	}
}

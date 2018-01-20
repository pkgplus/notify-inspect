package plugin

import (
	"encoding/json"
	"errors"
)

type Subscribe struct {
	UserId      string `json:"uid"`
	PluginId    string `json:"pluginId"`
	IsAvailable string `json:"isAvailable"`
	ErrMsg      string `json:"errMsg,omitempty"`

	Data []PluginData `json:"data"`
}

func (us *Subscribe) GetId() string {
	return GetSubscribeId(us.UserId, us.PluginId)
}

func GetSubscribeId(uid, pluginid string) string {
	return uid + "." + pluginid
}

func (us *Subscribe) ToJson() []byte {
	data, _ := json.Marshal(us)
	return data
}

func (us *Subscribe) Convert2Map() map[string]interface{} {
	data_bytes, _ := json.Marshal(us.Data)
	return map[string]interface{}{
		"uid":         us.UserId,
		"pluginId":    us.PluginId,
		"isAvailable": us.IsAvailable,
		"errMsg":      us.ErrMsg,
		"data":        string(data_bytes),
	}
}

func (us *Subscribe) GetParamValue(id string) string {
	for _, param := range us.Data {
		if param.Id == id {
			return param.Value
		}
	}
	return ""
}

func Map2Subscribe(values map[string]string) (*Subscribe, error) {
	data := make([]PluginData, 0)
	err := json.Unmarshal([]byte(values["data"]), &data)
	if err != nil {
		return nil, err
	}

	us := &Subscribe{
		UserId:      values["uid"],
		PluginId:    values["pluginId"],
		IsAvailable: values["isAvailable"],
		ErrMsg:      values["errMsg"],
		Data:        data,
	}
	if us.UserId == "" || us.PluginId == "" {
		return nil, errors.New("uid and pluginId must be not empty")
	}

	return us, nil
}

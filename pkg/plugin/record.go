package plugin

import (
	"encoding/json"
	"errors"

	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

type PluginRecord struct {
	Id       string                `json:"id"`
	UserId   string                `json:"uid"`
	PluginId string                `json:"pluginId"`
	Disable  string                `json:"disable"`
	Cron     *cron.CronTaskSetting `json:"cron,omitempty"`
	Data     []PluginData          `json:"data"`
	SubData  []PluginData          `json:"subData,omitempty"`
}

func (r *PluginRecord) ToJson() []byte {
	data, _ := json.Marshal(r)
	return data
}

func (pr *PluginRecord) Convert2Map() map[string]interface{} {
	data_bytes, _ := json.Marshal(pr.Data)

	pr.Cron.Init()
	pr.Cron.TaskId = pr.UserId + "." + pr.PluginId + "." + pr.Id

	cron_bytes, _ := json.Marshal(pr.Cron)
	return map[string]interface{}{
		"id":       pr.Id,
		"uid":      pr.UserId,
		"pluginId": pr.PluginId,
		"disable":  pr.Disable,
		"cron":     string(cron_bytes),
		"data":     string(data_bytes),
	}
}

func Map2PluginRecord(values map[string]string) (*PluginRecord, error) {
	data := make([]PluginData, 0)
	err := json.Unmarshal([]byte(values["data"]), &data)
	if err != nil {
		return nil, err
	}

	cron := new(cron.CronTaskSetting)
	err = json.Unmarshal([]byte(values["cron"]), cron)
	if err != nil {
		return nil, err
	}

	pur := &PluginRecord{
		Id:       values["id"],
		UserId:   values["uid"],
		PluginId: values["pluginId"],
		Disable:  values["disable"],
		Cron:     cron,
		Data:     data,
	}
	if pur.UserId == "" || pur.PluginId == "" {
		return nil, errors.New("uid and pluginId must be not empty")
	}

	return pur, nil
}

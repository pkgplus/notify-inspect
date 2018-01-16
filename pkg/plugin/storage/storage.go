package storage

import (
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

type Storage interface {
	SaveSubscribe(s *plugin.Subscribe) error
	GetSubscribe(uid, plugin string) (*plugin.Subscribe, error)
	ListSubscribes(uid string) ([]*plugin.Subscribe, error)
	// GetSubscribeCount(uid string) (int, error)
	DeleteSubscribe(uid, plugin string) error

	SavePluginRecord(r *plugin.PluginRecord) error
	GetPluginRecord(uid, pluginid, id string) (*plugin.PluginRecord, error)
	ListPluginRecords(uid, pluginid string) ([]*plugin.PluginRecord, error)
	DeleleAllPluginRecords(uid, pluginid string) (int64, error)
	GetPluginRecordCountByPlugin(uid, pluginid string) (int, error)
	DeletePluginRecord(uid, pluginid, id string) error
	ModifyPluginRecord(uid, pluginid, id string, data map[string]interface{}) (bool, error)
	// DisablePluginRecord(uid, pluginid, id string) (bool, error)
	// EnablePluginRecord(uid, pluginid, id string) (bool, error)
}

var (
	GlobalStorage Storage
)

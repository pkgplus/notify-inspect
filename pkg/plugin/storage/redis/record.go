package redis

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

const (
	USER_RECORD_PREFIX = "plugin_user_record:"
)

func (redis *RedisStorage) SavePluginRecord(r *plugin.PluginRecord) error {
	key := USER_RECORD_PREFIX + getPluginRecordKey(r.UserId, r.PluginId, r.Id)
	ret := redis.HMSet(key, r.Convert2Map())
	return ret.Err()
}

func (redis *RedisStorage) GetPluginRecord(uid, pluginid, id string) (*plugin.PluginRecord, error) {
	key := USER_RECORD_PREFIX + getPluginRecordKey(uid, pluginid, id)
	return redis.getPluginRecordByKey(key)
}

func (redis *RedisStorage) ListPluginRecords(uid, pluginid string) ([]*plugin.PluginRecord, error) {
	urs := make([]*plugin.PluginRecord, 0)

	key_prefix := USER_RECORD_PREFIX + uid + "." + pluginid + ".*"
	ret := redis.Keys(key_prefix)
	if ret.Err() != nil {
		return urs, ret.Err()
	}

	for _, key := range ret.Val() {
		ur, err := redis.getPluginRecordByKey(key)
		if err != nil {
			log.Printf("read \"%s\" from redis failed:%v", key, err)
			continue
		}
		urs = append(urs, ur)
	}

	return urs, nil
}

func (redis *RedisStorage) DeleleAllPluginRecords(uid, pluginid string) (int64, error) {
	key_prefix := USER_RECORD_PREFIX + uid + "." + pluginid + ".*"
	ret := redis.Del(key_prefix)
	return ret.Result()
}

func (redis *RedisStorage) GetPluginRecordCountByPlugin(uid, pluginid string) (int, error) {
	key_prefix := USER_RECORD_PREFIX + uid + "." + pluginid + ".*"
	ret := redis.Keys(key_prefix)
	if ret.Err() != nil {
		return 0, ret.Err()
	}

	return len(ret.Val()), nil
}

func (redis *RedisStorage) DeletePluginRecord(uid, pluginid, id string) error {
	rid := getPluginRecordKey(uid, pluginid, id)
	ret := redis.Del(USER_RECORD_PREFIX + rid)
	return ret.Err()
}

func (redis *RedisStorage) ModifyPluginRecord(uid, pluginid, id string, data map[string]interface{}) (bool, error) {
	rid := getPluginRecordKey(uid, pluginid, id)

	data_string := make(map[string]interface{})
	for k, v := range data {
		v_map, ok := v.(map[string]interface{})
		if ok {
			var err error
			data_string[k], err = json.Marshal(v_map)
			if err != nil {
				return false, fmt.Errorf("parse %s in %+v failed:%v", k, data, err)
			} else {
			}
		} else {
			data_string[k] = v
		}
	}

	ret := redis.HMSet(USER_RECORD_PREFIX+rid, data_string)
	if ret.Err() != nil {
		return false, ret.Err()
	} else if ret.Val() != "OK" {
		return false, fmt.Errorf("return %s", ret.Val())
	} else {
		return true, nil
	}
}

func (redis *RedisStorage) DisablePluginRecord(uid, pluginid, id string) (bool, error) {
	key := USER_RECORD_PREFIX + getPluginRecordKey(uid, pluginid, id)
	ret := redis.HSet(key, "disable", "1")
	return ret.Result()
}
func (redis *RedisStorage) EnablePluginRecord(uid, pluginid, id string) (bool, error) {
	key := USER_RECORD_PREFIX + getPluginRecordKey(uid, pluginid, id)
	ret := redis.HSet(key, "disable", "0")
	return ret.Result()
}

func (redis *RedisStorage) getPluginRecordByKey(key string) (*plugin.PluginRecord, error) {
	ret := redis.HGetAll(key)
	if ret.Err() != nil {
		return nil, ret.Err()
	}

	return plugin.Map2PluginRecord(ret.Val())
}

func getPluginRecordKey(uid, pluginid, id string) string {
	return uid + "." + pluginid + "." + id
}

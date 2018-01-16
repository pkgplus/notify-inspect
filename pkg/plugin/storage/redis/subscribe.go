package redis

import (
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

const (
	USER_SETTING_PREFIX = "plugin_user_setting."
)

func (redis *RedisStorage) SaveSubscribe(s *plugin.Subscribe) error {
	ret := redis.HMSet(USER_SETTING_PREFIX+s.GetId(), s.Convert2Map())
	return ret.Err()
}

func (redis *RedisStorage) ListSubscribes(uid string) ([]*plugin.Subscribe, error) {
	uss := make([]*plugin.Subscribe, 0)

	key_prefix := USER_SETTING_PREFIX + uid + ".*"
	ret := redis.Keys(key_prefix)
	if ret.Err() != nil {
		return uss, ret.Err()
	}

	for _, key := range ret.Val() {
		us, err := redis.getSubscribeByKey(key)
		if err != nil {
			redis.Errorf("read \"%s\" from redis failed:%v", key, err)
			continue
		}
		uss = append(uss, us)
	}

	return uss, nil
}

// func (redis *RedisStorage) GetSubscribeCount(uid string) (int, error) {
// 	key_prefix := USER_SETTING_PREFIX + uid + ".*"
// 	ret := redis.Keys(key_prefix)
// 	if ret.Err() != nil {
// 		return 0, ret.Err()
// 	}

// 	return len(ret.Val()), nil
// }

func (redis *RedisStorage) GetSubscribe(uid, pluginid string) (*plugin.Subscribe, error) {
	sid := plugin.GetSubscribeId(uid, pluginid)
	key := USER_SETTING_PREFIX + sid
	return redis.getSubscribeByKey(key)
}

func (redis *RedisStorage) getSubscribeByKey(key string) (*plugin.Subscribe, error) {
	ret := redis.HGetAll(key)
	if ret.Err() != nil {
		return nil, ret.Err()
	}

	return plugin.Map2Subscribe(ret.Val())
}

func (redis *RedisStorage) DeleteSubscribe(uid, pluginid string) error {
	sid := plugin.GetSubscribeId(uid, pluginid)
	ret := redis.Del(USER_SETTING_PREFIX + sid)
	return ret.Err()
}

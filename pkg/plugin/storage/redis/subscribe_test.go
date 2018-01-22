package redis

import (
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	"testing"
)

func TestGetSubscribe(t *testing.T) {
	uid := "odTIQ0dMoifGFIMrIoWA20G53-OA"
	pid := "hhsecret"
	ps, err := storage.GlobalStorage.GetSubscribe(uid, pid)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ps)
	if ps.Data == nil || len(ps.Data) != 2 {
		t.Fatal("plugin subscribe is empty!")
	}
}

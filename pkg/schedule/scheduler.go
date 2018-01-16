package schedule

import (
	"errors"
	"fmt"
	"time"

	"github.com/xuebing1110/notify-inspect/pkg/log"
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
	"github.com/xuebing1110/notify/pkg/client"
)

type Scheduler interface {
	PutTask(task *cron.CronTaskSetting, curtime time.Time) error
	RemoveTask(taskid string) error
	FetchTasks(curtime time.Time) <-chan *cron.CronTaskSetting
}

var (
	DefaultScheduler Scheduler
)

func Start() error {
	if DefaultScheduler == nil {
		return errors.New("Scheduler is null,pls import the schedule package")
	}

	go runEveryMinute()

	return nil
}

func runEveryMinute() {
	tick := time.Tick(time.Minute)
	for now := range tick {
		now_minute := now.Truncate(time.Minute)
		for task := range DefaultScheduler.FetchTasks(now_minute) {
			defer DefaultScheduler.PutTask(task, now_minute)
			runTask(task)
		}
	}
}

func runTask(task *cron.CronTaskSetting) error {
	uid, pid, rid, err := task.ParseTaskId()
	if err != nil {
		return err
	}

	p, found := plugin.DefaultRegisterServer.GetPlugin(pid)
	if !found {
		return fmt.Errorf("the plugin %s not found", pid)
	}

	r, err := storage.GlobalStorage.GetPluginRecord(uid, pid, rid)
	if err != nil {
		return err
	}

	// Disable
	if r.Disable != "False" && r.Disable != "false" && r.Disable != "0" {
		return nil
	}

	// plugin sub info
	s, err := storage.GlobalStorage.GetSubscribe(uid, pid)
	if err != nil {
		return err
	}
	r.SubData = s.Data

	// call the backend service of the plugin
	n, err := p.BackendInspect(r)
	if err != nil {
		return err
	}

	if n != nil {
		log.GlobalLogger.Infof("send a notice: %+v", n)
		err = client.SendNotice(n)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

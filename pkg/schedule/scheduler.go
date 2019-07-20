package schedule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/xuebing1110/notify-inspect/pkg/notice"
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

type Scheduler interface {
	PutTask(task *cron.CronTask, curtime time.Time) error
	RemoveTask(taskid string) error
	FetchTasks(curtime time.Time) <-chan *cron.CronTask
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
			err := runTask(task.Id)
			if err != nil {
				log.Printf("run %s failed:%v", task.Id, err.Error())
			}
			DefaultScheduler.PutTask(task, now_minute)
		}
	}
}

func runTask(taskid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uid, pid, rid, err := cron.ParseTaskId(taskid)
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
	n, err := p.BackendInspect(ctx, r)
	if err != nil {
		return err
	}

	if n != nil {
		log.Printf("send a notice: %+v", n)
		err = notice.Send(ctx, n)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

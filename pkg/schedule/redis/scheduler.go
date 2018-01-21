package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/xuebing1110/notify-inspect/pkg/log"
	myredis "github.com/xuebing1110/notify-inspect/pkg/redis"
	"github.com/xuebing1110/notify-inspect/pkg/schedule"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

const (
	TASKS_SORTSET = "crontasks"
	TASKS_DETAIL  = "crontasks.detail:"
)

type UseRedisScheduler struct {
	*redis.Client
	log.Logger
	tasks chan *cron.CronTaskSetting
}

func init() {
	schedule.DefaultScheduler = &UseRedisScheduler{
		Client: myredis.GetClient(),
		Logger: log.GlobalLogger,
	}
}

func (s *UseRedisScheduler) PutTask(task *cron.CronTask, curtime time.Time) error {
	// fmt.Printf("task: %+v  curtime:%s\n", task, curtime.String())
	next_time, err := task.Setting.NextRunTime(curtime)
	if err != nil {
		return err
	}

	hret := s.HMSet(TASKS_DETAIL+task.Id, task.Setting.Convert2Map())
	if hret.Err() != nil {
		return hret.Err()
	}

	zret := s.ZAdd(TASKS_SORTSET, redis.Z{
		float64(next_time.Unix()),
		task.Id,
	})
	return zret.Err()
}

func (s *UseRedisScheduler) RemoveTask(taskid string) error {
	ret := s.ZRem(TASKS_SORTSET, taskid)
	return ret.Err()
}

func (s *UseRedisScheduler) FetchTasks(curtime time.Time) <-chan *cron.CronTask {
	tasks := make(chan *cron.CronTask, 10)

	go s.fetch(curtime, tasks)
	return tasks
}

func (s *UseRedisScheduler) fetch(curtime time.Time, tasks chan *cron.CronTask) {
	defer close(tasks)

	ret := s.ZRevRangeByScoreWithScores(
		TASKS_SORTSET,
		redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%d", curtime.Unix()),
		})
	retZs, err := ret.Result()
	if err != nil {
		s.Errorf("fetch %s failed: %v", TASKS_SORTSET, err)
		return
	}

	for _, retZ := range retZs {
		taskid := retZ.Member.(string)
		s.Infof("found a task: %s", taskid)

		taskRet := s.HGetAll(TASKS_DETAIL + taskid)
		if err != nil {
			s.Errorf("get %s task detail failed: %v", taskid, err)
			continue
		}

		taskSetting, err := cron.NewCronTaskSettingFromMap(taskRet.Val())
		if err != nil {
			s.Errorf("parse %s failed: %v", taskid, err)
			continue
		}

		s.Infof("emit a task: %s", taskid)
		tasks <- &cron.CronTask{taskid, taskSetting}
	}
}

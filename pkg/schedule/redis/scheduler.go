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

func (s *UseRedisScheduler) PutTask(task *cron.CronTaskSetting, curtime time.Time) error {
	next_time, err := task.NextRunTime(curtime)
	if err != nil {
		return err
	}

	zret := s.ZAdd(TASKS_SORTSET, redis.Z{
		float64(next_time.Unix()),
		task.String(),
	})
	return zret.Err()
}

func (s *UseRedisScheduler) RemoveTask(taskid string) error {
	ret := s.ZRem(TASKS_SORTSET, taskid)
	return ret.Err()
}

func (s *UseRedisScheduler) FetchTasks(curtime time.Time) <-chan *cron.CronTaskSetting {
	tasks := make(chan *cron.CronTaskSetting, 10)

	go s.fetch(curtime, tasks)
	return tasks
}

func (s *UseRedisScheduler) fetch(curtime time.Time, tasks chan *cron.CronTaskSetting) {
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
		content := retZ.Member.(string)
		s.Infof("found a task: %s", content)

		task, err := schedule.NewCronTaskSetting([]byte(content))
		if err != nil {
			s.Errorf("parse %s failed: %v", content, err)
			continue
		}

		s.Infof("emit a task: %s", content)
		tasks <- task
	}
}

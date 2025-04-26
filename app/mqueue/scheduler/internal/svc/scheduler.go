package svc

import (
	"fmt"
	"github.com/hibiken/asynq"
	"graph-med/app/mqueue/scheduler/internal/config"
	"time"
)

func newScheduler(c config.Config) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")

	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Host,
			Password: c.Redis.Pass,
		},

		&asynq.SchedulerOpts{
			Location: location,
			PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
				if err != nil {
					fmt.Printf("Scheduler PostEnqueueFunc <<<<<<<===>>>>> err : %+v , taskInfo : %+v", err, info)
				}
			},
		},
	)
}

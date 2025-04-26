package logic

import (
	"context"
	"graph-med/app/mqueue/scheduler/internal/svc"
)

type MqueueScheduler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronScheduler(ctx context.Context, svcCtx *svc.ServiceContext) *MqueueScheduler {
	return &MqueueScheduler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MqueueScheduler) Register() {
	//l.settleRecordScheduler()
}

// scheduler job ------> go-zero-looklook/app/mqueue/cmd/job/internal/logic/settleRecord.go.
func (l *MqueueScheduler) settleRecordScheduler() {

	//task := asynq.NewTask(jobtype.ScheduleSettleRecord, nil)
	//
	//// every one minute exec
	//entryID, err := l.svcCtx.Scheduler.Register("*/1 * * * *", task)
	//if err != nil {
	//	logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【settleRecordScheduler】 registered  err:%+v , task:%+v", err, task)
	//}
	//fmt.Printf("【settleRecordScheduler】 registered an  entry: %q \n", entryID)
}

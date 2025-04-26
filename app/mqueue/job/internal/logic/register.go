package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"graph-med/app/mqueue/job/internal/handler"
	"graph-med/app/mqueue/job/internal/svc"
	"graph-med/app/mqueue/job/jobtype"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CronJob) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.Handle(jobtype.CaptchaSendEmail, handler.NewSendEmailHandler(l.svcCtx))

	return mux
}

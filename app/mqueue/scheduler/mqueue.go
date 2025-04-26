package main

import (
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"graph-med/app/mqueue/scheduler/internal/config"
	"graph-med/app/mqueue/scheduler/internal/logic"
	"graph-med/app/mqueue/scheduler/internal/svc"
	"os"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.DisableStat()

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	mqueueScheduler := logic.NewCronScheduler(ctx, svcContext)
	mqueueScheduler.Register()

	if err := svcContext.Scheduler.Run(); err != nil {
		logx.Errorf("!!!MqueueSchedulerErr!!!  run err:%+v", err)
		os.Exit(1)
	}
}

package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icodefans/go-extend/service"
	cron_v3 "github.com/robfig/cron/v3"
)

type cron struct {
	router []struct {
		group  string     // 业务分组
		spec   string     // 间隔时长
		hander cronHander // 处理程序
	}
}

// 任务处理程序类型定义
type cronHander func(ctx *context.Context) *service.Result

type cronJob struct {
	ctx    *context.Context
	hander cronHander
}

func (c cronJob) Run() {
	c.hander(c.ctx)
}

var cronInstance *cron

// 实例化对象，单例模式
func Cron() *cron {
	if cronInstance == nil {
		cronInstance = new(cron)
	}
	return cronInstance
}

// 增加路由配置
func (w *cron) Add(group, spec string, hander cronHander) {
	w.router = append(w.router, struct {
		group  string
		spec   string
		hander cronHander
	}{
		group,
		spec,
		hander,
	})
}

// 定时任务运行
func (w *cron) Run(group string) {
	ctx, cancel := context.WithCancel(context.Background())
	c := cron_v3.New()
	for _, value := range w.router {
		if value.group != group {
			continue
		} else if _, err := c.AddJob(value.spec, cronJob{
			hander: value.hander,
			ctx:    &ctx,
		}); err != nil {
			log.Printf("Job Spec (%s,%s)配置错误:%s", value.group, value.spec, err)
		}
	}
	c.Start()

	// 监听进程退出信号(TERM, HUP, INT, QUIT, KILL, USR1, or USR2)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2)
	sig := <-s

	// 进程上下文取消
	cancel()

	// 预留时间，执行收尾操作
	time.Sleep(3 * time.Second)
	service.Success(fmt.Sprintf("程序退出，接收到信号:%s", sig.String()), nil, "info")
}

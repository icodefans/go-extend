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
		sync   uint8      // 同步执行(0:异步执行,1:同步执行,等待上一次完成,2:同步执行,跳过重复)
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
func (w *cron) Add(group, spec string, sync uint8, hander cronHander) {
	w.router = append(w.router, struct {
		group  string
		spec   string
		sync   uint8
		hander cronHander
	}{
		group,
		spec,
		sync,
		hander,
	})
}

// 定时任务运行
func (w *cron) Run(group string) {
	// 创建基础 cron 实例（仅启用秒级，不设置全局执行链）
	c := cron_v3.New(cron_v3.WithSeconds())
	// 自定义同步链：DelayIfStillRunning（等待上一个任务完成）
	syncChain := cron_v3.NewChain(cron_v3.DelayIfStillRunning(cron_v3.DefaultLogger))
	// 自定义同步链：SkipIfStillRunning（上一个未完成则跳过）
	skipChain := cron_v3.NewChain(cron_v3.SkipIfStillRunning(cron_v3.DefaultLogger))

	ctx, cancel := context.WithCancel(context.Background())
	for _, value := range w.router {
		if value.group != group {
			continue
		}
		myJob := cronJob{hander: value.hander, ctx: &ctx}
		if value.sync != 0 {
			// next
		} else if _, err := c.AddJob(value.spec, myJob); err != nil {
			log.Printf("Job Spec (%s,%s)异步任务配置错误:%s", value.group, value.spec, err)
		}
		if value.sync != 1 {
			// next
		} else if _, err := c.AddJob(value.spec, syncChain.Then(myJob)); err != nil {
			log.Printf("Job Spec (%s,%s)同步1配置错误:%s", value.group, value.spec, err)
		}
		if value.sync != 2 {
			// next
		} else if _, err := c.AddJob(value.spec, skipChain.Then(myJob)); err != nil {
			log.Printf("Job Spec (%s,%s)同步2配置错误:%s", value.group, value.spec, err)
		}
	}
	c.Start()

	// 监听进程退出信号(TERM, HUP, INT, QUIT, KILL, USR1, or USR2)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2)
	sig := <-s

	// 停止 cron（可选，优雅退出）
	c.Stop()

	// 进程上下文取消
	cancel()

	// 预留时间，执行收尾操作
	time.Sleep(3 * time.Second)
	fmt.Printf("%s - cron server stop signal:%s\n\n", time.Now().Format(`2006/01/02 15:04:05`), sig.String())
	// service.Success(fmt.Sprintf("程序退出，接收到信号:%s", sig.String()), nil, "info")
}

package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icodefans/go-extend/service"
)

type cli struct {
	router map[string]cliHander
}

// 任务处理程序类型定义
type cliHander func(ctx *context.Context) *service.Result

var cliInstance *cli

// 实例化对象，单例模式
func Cli() *cli {
	if cliInstance == nil {
		cliInstance = new(cli)
		cliInstance.router = make(map[string]cliHander)
	}
	return cliInstance
}

// 增加路由配置
func (w *cli) Add(rule string, hander cliHander) {
	w.router[rule] = hander
}

// 路由匹配运行
func (w *cli) Match(rule string) (result *service.Result) {
	if hander, ok := w.router[rule]; !ok {
		return service.Error(101, fmt.Sprintf("没有匹配到路由规则(%s)", rule), nil, "error")
	} else {
		// 上下文取消控制
		ctx, cancel := context.WithCancel(context.Background())

		// 在这里执行异步的其他代码
		go hander(&ctx)

		// 监听进程退出信号(TERM, HUP, INT, QUIT, KILL, USR1, or USR2)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2)
		sig := <-c

		// 进程上下文取消
		cancel()

		// 预留时间，执行收尾操作
		time.Sleep(5 * time.Second)
		return service.Success("收尾成功，程序退出，执行完成", []any{fmt.Sprintf("接收到信号:%s", sig.String())}, "dev")
	}
}

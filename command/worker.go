package command

import (
	"fmt"

	"github.com/icodefans/go-extend/service"
)

type worker struct {
	router map[string]workerHander
}

// 任务处理程序类型定义
type workerHander func() *service.Result

var workerInstance *worker

// 实例化对象，单例模式
func Worker() *worker {
	if workerInstance == nil {
		workerInstance = new(worker)
		workerInstance.router = make(map[string]workerHander)
	}
	return workerInstance
}

// 增加路由配置
func (w *worker) Add(rule string, hander workerHander) {
	w.router[rule] = hander
}

// 路由匹配运行
func (w *worker) Match(rule string) (result *service.Result) {
	if hander, ok := w.router[rule]; !ok {
		return service.Error(102, fmt.Sprintf("没有匹配到路由规则(%s)", rule), nil, "error")
	} else {
		return hander()
	}
}

package sdk

import (
	"context"
	"time"
)

var exec Executor

// NewExecutor 创建执行器
func NewExecutor(ctx context.Context, options ...Option) (err error) {
	e := &executor{
		ctx: ctx,
		regList: &ListTask{
			data: make(map[string]*Task),
		},
		runList: &ListTask{
			data: make(map[string]*Task),
		},
	}
	for _, o := range options {
		o(e)
	}
	err = e.Run()
	if err != nil {
		return
	}
	return
}

type Option func(e *executor)

// WithHost 调度服务地址
func WithHost(host string) Option {
	return func(e *executor) {
		e.host = host
		//todo 连接测试
	}
}

// WithClientId 客户端ID/应用模块ID
func WithClientId(clientId string) Option {
	return func(e *executor) {
		e.clientId = clientId
	}
}

// WithTask 本地注册任务
func WithTask(handler string, task FuncTask) Option {
	return func(e *executor) {
		e.AddRegTask(handler, task)
	}
}
// WithDelay 延时获取任务实例。默认1s。
func WithDelay(delay time.Duration) Option {
	return func(e *executor) {
		e.delay = delay
	}
}
// WithCancelReg 是否取消-自动注册任务
func WithCancelReg(cancel bool) Option {
	return func(e *executor) {
		e.cancelReg = cancel
	}
}

// WithHttpTimeout 超时时间
func WithHttpTimeout(timeout time.Duration) Option {
	return func(e *executor) {
		e.httpTimeout = timeout
	}
}

// 监控回调 、 超时处理、

// GetExecutor 获取执行器实例
func GetExecutor() Executor {
	return exec
}

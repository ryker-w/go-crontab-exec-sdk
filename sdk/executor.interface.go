package sdk

import "github.com/ryker-w/go-crontab-exec-sdk/sdk/common"

// Executor 执行器
type Executor interface {
	Run() (err error)
	// AddRegTask 本地注册
	AddRegTask(handler string, task FuncTask)
	// RegTask 服务器注册
	RegTask() (err error)
	// GetTaskInstance 获取调度任务实例
	GetTaskInstance() common.RunReq
	// RunTask 运行任务
	RunTask(req common.RunReq) common.CallElement
}

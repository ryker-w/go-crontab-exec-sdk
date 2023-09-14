package sdk

import (
	"fmt"
	"github.com/lishimeng/go-log"
	"github.com/ryker-w/go-crontab-exec-sdk/sdk/common"
	"runtime/debug"
)

// FuncTask 任务执行函数
type FuncTask func(param common.RunReq) (code int, msg string)

// Task 任务
type Task struct {
	Id        string
	Name      string
	Param     common.RunReq
	fn        FuncTask
	StartTime int64
	EndTime   int64
}

// Run 运行任务
func (t *Task) Run() common.CallElement {
	defer func() common.CallElement {
		if err := recover(); err != nil {
			log.Info(err)
			debug.PrintStack() //堆栈跟踪
			return common.Callback(t.Param, common.Error, fmt.Sprintf("panic: %v", err))
		}
		return common.Callback(t.Param, common.Success, "")
	}()
	code, msg := t.fn(t.Param)
	return common.Callback(t.Param, code, msg)
}

// Info 任务信息
func (t *Task) Info() string {
	return fmt.Sprintf("任务ID[%s]任务名称[%s]", t.Id, t.Name)
}

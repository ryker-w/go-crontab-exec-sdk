package common

import (
	"github.com/lishimeng/app-starter"
	"time"
)

// CallElement 执行器执行完任务后，回调任务结果时使用
type CallElement struct {
	HandleCode      int    `json:"handleCode"` //200表示正常,500表示失败
	HandleMsg       string `json:"handleMsg"`
	JobId           string `json:"jobId"`
	InstanceId      int    `json:"instanceId"`
	ExecutorHandler string `json:"ExecutorHandler"`
}

// RunReq 触发任务请求参数
type RunReq struct {
	JobID           string                 `json:"jobId"`           // 任务ID
	ExecutorHandler string                 `json:"executorHandler"` // 任务标识
	InstanceId      int                    `json:"instanceId"`      // 任务实例ID
	Time            time.Time              `json:"time"`            // 出发请求的时间
	Data            map[string]interface{} `json:"data"`
}

type ScheduleResponse struct {
	app.Response
	Data RunReq `json:"data,omitempty"`
}

type RegReq struct {
	ClientId string   `json:"clientId"`
	Handlers []string `json:"handlers"`
}

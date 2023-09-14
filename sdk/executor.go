package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
	"github.com/ryker-w/go-crontab-exec-sdk/sdk/common"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	defaultDelay       = 1 * time.Second
	httpDefaultTimeout = 8 * time.Second
)

type executor struct {
	ctx         context.Context
	host        string        // 调度服务地址
	clientId    string        // 执行器标识
	cancelReg   bool          // 取消自动注册。默认false，即默认自动远程注册
	delay       time.Duration // 延时自动获取。单位：秒。默认1s。
	httpTimeout time.Duration // 请求超时时间
	regList     *ListTask     // 注册任务列表
	runList     *ListTask     // 正在执行任务列表
	mu          sync.RWMutex  // RW锁

}

func (e *executor) Run() (err error) {
	if len(e.host) == 0 || len(e.clientId) == 0 {
		return errors.New("host 和 clientId 不能为空")
	}
	if e.delay == 0 {
		e.delay = defaultDelay
	}
	if e.httpTimeout == 0 {
		e.httpTimeout = httpDefaultTimeout
	}
	if !e.cancelReg {
		err = e.RegTask()
		if err != nil {
			return
		}
	}
	e.run()
	return err
}

// AddRegTask 本地注册任务
func (e *executor) AddRegTask(handler string, task FuncTask) {
	if e.regList == nil {
		e.regList = &ListTask{
			data: make(map[string]*Task),
		}
	}
	var t = &Task{}
	t.fn = task
	t.Name = handler
	e.regList.Set(handler, t)
	log.Debug("执行器注册本地任务：%s", handler)
}

// RegTask 远程注册任务
func (e *executor) RegTask() (err error) {
	if len(e.clientId) == 0 {
		return errors.New("clientId 不能为空")
	}
	if e.regList == nil || e.regList.Len() == 0 {
		return
	}
	err = e.regTask()
	if err != nil {
		return err
	}
	return
}

func (e *executor) regTask() (err error) {
	req := common.RegReq{
		ClientId: e.clientId,
		Handlers: e.regList.GetKeys(),
	}
	var requestUrl = e.buildRemoteUrl(common.RouteRegTaskTpl)
	client := &http.Client{Timeout: common.HttpDefaultTimeout}
	jsonStr, _ := json.Marshal(req)
	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Debug("注册失败，client Post err")
		log.Debug(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	result, _ := io.ReadAll(resp.Body)
	var response app.Response
	err = json.Unmarshal(result, &response)
	if err != nil {
		log.Debug("注册失败，response Unmarshal err")
		log.Debug(err)
		return
	}
	if response.Code != float64(tool.RespCodeSuccess) {
		return errors.New(fmt.Sprintf("注册失败, response code err: %v", response.Code))
	}
	return
}

func (e *executor) GetTaskInstance() common.RunReq {
	return e.getTaskInstance()
}

// RunTask 运行任务
func (e *executor) RunTask(req common.RunReq) common.CallElement {
	return e.runTask(req)
}

//运行一个任务
func (e *executor) runTask(req common.RunReq) common.CallElement {
	e.mu.Lock()
	defer e.mu.Unlock()
	// 本地执行
	if !e.regList.Exists(req.ExecutorHandler) {
		return common.Callback(req, common.NotFount, fmt.Sprintf("任务[ %s ]没有注册: %s", req.JobID, req.ExecutorHandler))
	}
	task := e.regList.Get(req.ExecutorHandler)
	task.Id = req.JobID
	task.Name = req.ExecutorHandler
	task.Param = req

	if e.runList == nil {
		e.runList = &ListTask{
			data: make(map[string]*Task),
		}
	}
	e.runList.Set(task.Id, task)
	log.Debug("任务[%s]开始执行: %s", req.JobID, req.ExecutorHandler)
	result := task.Run()
	log.Debug("[%s]执行结果: code=%d, msg=%s", result.ExecutorHandler, result.HandleCode, result.HandleMsg)
	return result
}

func (e *executor) getTaskInstance() common.RunReq {
	var requestUrl = e.buildRemoteUrl(common.RouteScheduleTpl)
	client := &http.Client{Timeout: common.HttpDefaultTimeout}
	resp, err := client.Get(requestUrl)
	if err != nil {
		log.Debug("client Post err")
		log.Debug(err)
		return common.RunReq{}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	result, _ := io.ReadAll(resp.Body)
	var response common.ScheduleResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		log.Debug("response Unmarshal err")
		log.Debug(err)
		return common.RunReq{}
	}
	return response.Data
}

func (e *executor) buildRemoteUrl(routePath string) string {
	requestUrl, _ := url.JoinPath(e.host, fmt.Sprintf(routePath, e.clientId))
	return requestUrl
}

func (e *executor) response(callElement common.CallElement) error {
	var requestUrl = e.buildRemoteUrl(common.RouteScheduleTpl)
	client := &http.Client{Timeout: common.HttpDefaultTimeout}
	jsonStr, _ := json.Marshal(callElement)
	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Debug("client Post err")
		log.Debug(err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	result, _ := io.ReadAll(resp.Body)
	var response app.Response
	err = json.Unmarshal(result, &response)
	if err != nil {
		log.Debug("response Unmarshal err")
		log.Debug(err)
		return err
	}
	if response.Code != float64(tool.RespCodeSuccess) {
		return errors.New(fmt.Sprintf("response err code:%d", response.Code))
	}
	return nil
}

func (e *executor) run() {
	go func() {
		var timer = time.NewTimer(e.delay)
		defer func() {
			timer.Stop()
		}()
		for {
			select {
			case <-e.ctx.Done():
				return
			case <-timer.C:
				req := e.GetTaskInstance()
				if req.InstanceId != 0 {
					callElement := e.RunTask(req)
					err := e.response(callElement)
					if err != nil {
						log.Info(err)
					}
				}
				timer.Reset(e.delay)
			}
		}
	}()
}

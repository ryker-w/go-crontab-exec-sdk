package common

// Callback 封装任务回调
func Callback(req RunReq, code int, msg string) CallElement {
	data := CallElement{
		HandleCode:      code,
		HandleMsg:       msg,
		JobId:           req.JobID,
		ExecutorHandler: req.ExecutorHandler,
		InstanceId:      req.InstanceId,
	}
	return data
}

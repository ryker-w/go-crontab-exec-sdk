package common

const baseApi = "/monitor"

const (
	ActionRun      = baseApi + "/run"       // 启动任务
	ActionKill     = baseApi + "/kill"      // 终止任务
	ActionLog      = baseApi + "/log"       // 任务日志
	ActionBeat     = baseApi + "/heartbeat" // 心跳检测
	ActionIdleBeat = baseApi + "/busy"      // 忙碌检测
)

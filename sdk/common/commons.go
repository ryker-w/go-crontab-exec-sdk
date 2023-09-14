package common

import "time"

const RouteScheduleTpl = "/schedule/%s"
const RouteSchedulePath = "/schedule/{clientId}"

const RouteRegTaskTpl = "/register/%s"
const RouteRegTaskPath = "/register/{clientId}"

const HttpDefaultTimeout = 8 * time.Second

const (
	Success  = 200
	NotFount = 404
	Error    = 500
)

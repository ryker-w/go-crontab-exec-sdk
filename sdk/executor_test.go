package sdk

import (
	"context"
	"github.com/lishimeng/go-log"
	"github.com/ryker-w/go-crontab-exec-sdk/sdk/common"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	err := NewExecutor(context.Background(),
		WithHost("https://open.thingplecloud.com/crontab"), WithClientId("sdk-test-client1"),
		WithTask("demo1_task", demo1),
		WithTask("hello_world", helloWorld))
	if err != nil {
		t.Fatal(err)
		return
	}

	select {}
}

func demo1(param common.RunReq) (code int, msg string) {
	log.Info("demo1 run, receiving params: %+v", param)
	return common.Success, ""
}

func helloWorld(param common.RunReq) (code int, msg string) {
	log.Info("helloWorld")
	return common.Success, ""
}

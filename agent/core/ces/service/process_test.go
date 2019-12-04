package services

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//StartProcessInfoCollectTask
func TestStartProcessInfoCollectTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_StartProcessInfoCollectTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
			StartProcessInfoCollectTask(nil)
		})
	})
}

//SendProcessInfoTask
func TestSendProcessInfoTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendProcessInfoTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(BuildURL, func(destURI string) string {
				t.SkipNow()
				return ""
			})
			metrics := make(chan model.ChProcessList, 1)
			metrics <- model.ChProcessList{}
			SendProcessInfoTask(metrics)
		})
	})
}

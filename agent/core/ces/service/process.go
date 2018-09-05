package services

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/robfig/cron"
)

// StartProcessInfoCollectTask cron job for collecting top 5 cpu usage process info
func StartProcessInfoCollectTask(plist chan model.ChProcessList) {

	c := cron.New()

	c.AddFunc("0 * * * * *", func() {
		plist <- model.GetTop5CpuProcessList()
	})

	c.Start()

}

// SendProcessInfoTask task for send process info
func SendProcessInfoTask(plist chan model.ChProcessList) {
	for {
		processList := <-plist
		report.SendProcessInfo(BuildURL(cesUtils.PostProcessInfo), processList)
	}

}
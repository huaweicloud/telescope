package services

import (
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// StartProcessInfoCollectTask cron job for collecting top 5 cpu usage process info
func StartProcessInfoCollectTask(plist chan model.ChProcessList) {
	ticker := time.NewTicker(time.Duration(cesUtils.Top5ProcessCollectPeriodInSeconds) * time.Second)
	for _ = range ticker.C {
		logs.GetCesLogger().Info("GetTop5CpuProcessList starts")
		plist <- model.GetTop5CpuProcessList()
	}
}

// SendProcessInfoTask task for send process info
func SendProcessInfoTask(plist chan model.ChProcessList) {
	for {
		processList := <-plist
		report.SendProcessInfo(BuildURL(cesUtils.PostProcessInfo), processList)
	}

}

package services

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
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
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	for {
		processList := <-plist
		report.SendProcessInfo(client, BuildURL(cesUtils.PostProcessInfo), processList)
	}

}

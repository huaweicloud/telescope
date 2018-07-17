package services

import (
	"github.com/huaweicloud/telescope/agent/core/logs"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// CollectPluginTask cron job for collecting plugin data
func CollectPluginTask(data chan *model.InputMetric) {
	if !config.GetConfig().Enable || !config.GetConfig().EnablePlugin {
		return
	}

	if config.GetPluginConfig() == nil {
		return
	}

	plugins := config.GetPluginConfig().Plugins

	if len(plugins) > cesUtils.MaxPluginNum {
		plugins = plugins[:cesUtils.MaxPluginNum]
	}

	for _, eachPlugin := range plugins {
		logs.GetCesLogger().Debugf("Plugin info is %v", *eachPlugin)

		eachPluginSchedule := model.NewPluginScheduler(eachPlugin)
		if eachPluginSchedule == nil {
			return
		}
		go eachPluginSchedule.Schedule(data)
	}
}

// SendPluginTask task for post plugin data
func SendPluginTask(data chan *model.InputMetric) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	for {

		pluginData := <-data
		logs.GetCesLogger().Debugf("Plugin metric data is %v", *pluginData)
		report.SendMetricData(client, BuildURL(cesUtils.PostAggregatedMetricDataURI), pluginData, true)

	}
}

package services

import (
	"github.com/huaweicloud/telescope/agent/core/logs"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/json-iterator/go"
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// CollectCustomMonitorPluginTask cron job for collecting plugin data
func CollectCustomMonitorPluginTask(data chan model.CesMetricDataArr) {
	if !config.GetConfig().Enable || !config.GetConfig().EnablePlugin {
		return
	}

	plugins := config.GetCustomMonitorPluginConfig()
	if plugins == nil {
		return
	}
	if len(plugins) > cesUtils.MaxPluginNum {
		plugins = plugins[:cesUtils.MaxPluginNum]
	}

	for _, eachPlugin := range plugins {
		logs.GetCesLogger().Debugf("Plugin type is custom monitor, info is %v", *eachPlugin)

		eachPluginSchedule := model.NewCustomMonitorPluginScheduler(eachPlugin)
		if eachPluginSchedule == nil {
			return
		}
		go eachPluginSchedule.Schedule(data)
	}
}

// SendCustomMonitorPluginTask task for post plugin data
func SendCustomMonitorPluginTask(data chan model.CesMetricDataArr) {
	for {
		customMonitorMetricData := <-data
		logs.GetCesLogger().Debugf("Custom monitor plugin data is %v", customMonitorMetricData)

		customMonitorMetricDataInBytes, err := json.Marshal(customMonitorMetricData)
		if err != nil {
			logs.GetCesLogger().Errorf("Failed marshall custom monitor plugin data. Error: %s", err.Error())
			return
		}

		report.SendData(BuildURL(cesUtils.PostCustomMonitorMetricDataURI), customMonitorMetricDataInBytes)
	}
}

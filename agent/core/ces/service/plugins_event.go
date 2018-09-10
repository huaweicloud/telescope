package services

import (
	"github.com/huaweicloud/telescope/agent/core/logs"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
)

// CollectEventPluginTask cron job for collecting plugin data
func CollectEventPluginTask(data chan model.CesEventDataArr) {
	if !config.GetConfig().Enable || !config.GetConfig().EnablePlugin {
		return
	}

	plugins := config.GetEventPluginConfig()
	if plugins == nil {
		logs.GetCesLogger().Infof("Event plugin config is empty.")
		return
	}
	if len(plugins) > cesUtils.MaxPluginNum {
		plugins = plugins[:cesUtils.MaxPluginNum]
	}

	for _, eachPlugin := range plugins {
		logs.GetCesLogger().Debugf("Plugin type is event, info is %v", *eachPlugin)

		eachPluginSchedule := model.NewEventPluginScheduler(eachPlugin)
		if eachPluginSchedule == nil {
			return
		}
		go eachPluginSchedule.Schedule(data)
	}
}

// SendEventPluginTask task for post plugin data
func SendEventPluginTask(data chan model.CesEventDataArr) {
	for {
		eventMetricData := <-data
		logs.GetCesLogger().Debugf("Event plugin data is %v", eventMetricData)

		eventMetricDataInBytes, err := json.Marshal(eventMetricData)
		if err != nil {
			logs.GetCesLogger().Errorf("Failed marshall event data. Error: %s", err.Error())
			return
		}

		report.SendData(BuildURL(cesUtils.PostEventDataURI), eventMetricDataInBytes)
	}
}

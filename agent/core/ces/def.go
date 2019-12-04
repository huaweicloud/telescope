package ces

import (
	"os"

	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// common variables (chans and vars)
var (
	// Channels
	chPluginData, chSpeProcData chan *model.InputMetric
	chProcessInfo               chan model.ChProcessList
	chCustomMonitorData         chan model.CesMetricDataArr
	chEventData                 chan model.CesEventDataArr
)

// Initialize the original data channel
func initchPluginData() {
	chPluginData = make(chan *model.InputMetric, 100)
}

// Get the original data channel
func getchPluginData() chan *model.InputMetric {
	if chPluginData == nil {
		initchPluginData()
	}

	return chPluginData
}

func initchProcessInfo() {
	chProcessInfo = make(chan model.ChProcessList, 10)
}

func getchProcessInfo() chan model.ChProcessList {
	if chProcessInfo == nil {
		initchProcessInfo()
	}
	return chProcessInfo
}

// Initialize the original data channel
func initchSpeProcData() {
	chSpeProcData = make(chan *model.InputMetric, 100)
}

// Get the original data channel
func getchSpeProcData() chan *model.InputMetric {
	if chSpeProcData == nil {
		initchSpeProcData()
	}

	return chSpeProcData
}

func getchCustomMonitorData() chan model.CesMetricDataArr {
	if chCustomMonitorData == nil {
		chCustomMonitorData = make(chan model.CesMetricDataArr, 100)
	}
	return chCustomMonitorData
}

func getchEventData() chan model.CesEventDataArr {
	if chEventData == nil {
		chEventData = make(chan model.CesEventDataArr, 100)
	}
	return chEventData
}

func initEnvVariables() {
	err := os.Setenv(cesUtils.EnvInstanceID, utils.GetConfig().InstanceId)
	if err != nil {
		logs.GetCesLogger().Errorf("Set environment(CES_EVN_INSTANCE_ID) variable failed, error is: %v", err)
	}
}

func updateEnvVariables() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		initEnvVariables()
	}
}

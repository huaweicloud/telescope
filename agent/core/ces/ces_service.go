package ces

import (
	"github.com/buger/jsonparser"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/service"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/json-iterator/go"
	"github.com/shirou/gopsutil/process"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Service is one of the services of agent
type Service struct {
}

// Init ces Service config and channel
func (s *Service) Init() {
	config.InitConfig()
	config.InitPluginConfig()
	initchProcessInfo()
	initchPluginData()
	initEnvVariables()
	go updateEnvVariables()
}

// Start make work goroutines
func (s *Service) Start() {
	go updateConfig()
	go services.StartMetricCollectTask()

	go services.StartProcessInfoCollectTask(getchProcessInfo())
	go services.SendProcessInfoTask(getchProcessInfo())

	// plugin Agt metric-data
	go services.CollectPluginTask(getchPluginData())
	go services.SendPluginTask(getchPluginData())

	// plugin Custom Monitor
	go services.CollectCustomMonitorPluginTask(getchCustomMonitorData())
	go services.SendCustomMonitorPluginTask(getchCustomMonitorData())

	// plugin Event
	go services.CollectEventPluginTask(getchEventData())
	go services.SendEventPluginTask(getchEventData())

}

func updateConfig() {
	for {
		cesConfig := <-channel.GetCesConfigChan()
		logs.GetCesLogger().Debugf("Ces config is %s", cesConfig)

		// update ces enable
		cesEnable, err := jsonparser.GetBoolean([]byte(cesConfig), "enable")
		if err != nil {
			logs.GetCesLogger().Errorf("Failed to parse config :[%s] to get enable, error is %s", cesConfig, err.Error())
			continue
		}
		if cesEnable != config.GetConfig().Enable {
			config.GetConfig().Enable = cesEnable
		}

		// update enable_processes
		enableProcessStr, err := jsonparser.GetString([]byte(cesConfig), "enable_processes")
		if err != nil {
			logs.GetCesLogger().Errorf("Failed to parse config :[%s] to get enable_processes, error is %s", cesConfig, err.Error())
			continue
		}

		hbProcessList := config.GetConfig().EnableProcessList
		existProcessList := []config.HbProcess{}

		unmarshalErr := json.Unmarshal([]byte(enableProcessStr), &hbProcessList)
		if unmarshalErr != nil {
			logs.GetCesLogger().Errorf("Failed to unmarshal enable process list, error is %s", unmarshalErr.Error())
			continue
		}

		// check process is exist
		var notExistProcessList model.ChProcessList
		for _, eachHbProcess := range hbProcessList {
			isExist, err := process.PidExists(eachHbProcess.Pid)
			if err != nil {
				logs.GetCesLogger().Errorf("Failed to check process exist, error is %s", err.Error())
				continue
			}
			if !isExist {
				notExistProcess := model.ProcessInfo{}
				notExistProcess.Pid = eachHbProcess.Pid
				notExistProcess.Pname = eachHbProcess.Pname
				notExistProcess.Status = false
				notExistProcessList = append(notExistProcessList, &notExistProcess)
			} else {
				existProcessList = append(existProcessList, eachHbProcess)
			}
		}

		if notExistProcessList != nil && len(notExistProcessList) > 0 {
			getchProcessInfo() <- notExistProcessList
		}

		config.GetConfig().EnableProcessList = existProcessList

		// update specified process config
		specifiedProcessStr, _, _, err := jsonparser.Get([]byte(cesConfig), "specified_procs")
		if err != nil {
			logs.GetCesLogger().Errorf("Failed to parse config :[%s] to get specified_procs, error is %s", cesConfig, err.Error())
			continue
		}
		specifiedProcList := config.GetConfig().SpecifiedProcList
		unmarshalErrp := json.Unmarshal(specifiedProcessStr, &specifiedProcList)
		if unmarshalErrp != nil {
			logs.GetCesLogger().Errorf("Failed to unmarshal specified process list, error is %s", unmarshalErrp.Error())
			continue
		}

		config.GetConfig().SpecifiedProcList = specifiedProcList
	}
}

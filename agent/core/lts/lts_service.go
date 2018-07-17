package lts

import (
	"encoding/json"

	"github.com/buger/jsonparser"

	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/config"
	"github.com/huaweicloud/telescope/agent/core/lts/logdumper"
	"github.com/huaweicloud/telescope/agent/core/lts/services"
)

type LtsService struct {
}

func (s *LtsService) Init() {
	config.InitConfig()
	services.InitchData()
}

func (s *LtsService) Start() {
	go updateConfig()
	go services.StartExtractionTask()
	go services.StartDataService(services.GetchData())
}

//根据心跳返回的response更新conf_lts.json文件（不更新内存中的config）
func updateConfig() {
	//consume lts config channel
	for {
		ltsConfig := <-channel.GetLtsConfigChan()
		ltsEnable, err := jsonparser.GetBoolean([]byte(ltsConfig), "enable")
		if err != nil {
			logs.GetLtsLogger().Errorf("Failed to parse config :[%s] to get enable, error is %s", ltsConfig, err.Error())
			continue
		}

		if ltsEnable != config.GetConfig().Enable {
			config.GetConfig().Enable = ltsEnable
		}

		confStr, err := jsonparser.GetString([]byte(ltsConfig), "conf")
		if err != nil {
			logs.GetLtsLogger().Warnf("Failed to parse config :[%s] to get conf, error is %s", ltsConfig, err.Error())
			continue
		}
		remoteConfig := config.LTSConfigRemote{}
		unmarshallErr := json.Unmarshal([]byte(confStr), &remoteConfig)
		if unmarshallErr != nil {
			logs.GetLtsLogger().Warnf("Failed to unmarshal remote config, error is %s ", unmarshallErr.Error())
			continue
		}

		localGroupsBytes, _ := json.Marshal(config.GetConfig().Groups)
		remoteGroupsBytes, _ := json.Marshal(remoteConfig.Groups)

		if remoteConfig.Groups != nil && string(localGroupsBytes) != string(remoteGroupsBytes) && len(remoteGroupsBytes) > 0 {
			logs.GetLtsLogger().Debugf("Remote groups config: %s", string(remoteGroupsBytes))
			logs.GetLtsLogger().Debugf("Local groups config: %s", string(localGroupsBytes))
			newConfig := config.LTSConfig{}
			newConfig.Enable = config.GetConfig().Enable
			newConfig.Endpoint = config.GetConfig().Endpoint
			newConfig.Groups = remoteConfig.Groups
			success := config.UpdateConfig(newConfig)
			if !success {
				logs.GetLtsLogger().Error("Failed to update lts config.")
			} else {
				logs.GetLtsLogger().Debug("Success to update lts config.")
			}
		}

	}
}

func CronUpdateConfig() {
	isReload := config.ReloadConfig()
	if isReload {
		logdumper.ReloadExtractors()
	}
}

package config

import (
	"os"
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

// AssistantConfig ...
// CESConfig is the type for heartbeat response
type AssistantConfig struct {
	Enable          bool
	Endpoint        string
	EnablePlugin    bool
	ExternalService string
}

var (
	json            = jsoniter.ConfigCompatibleWithStandardLibrary
	assistantConfig *AssistantConfig
	mutex           sync.Mutex
)

// ReadConfig Read the config from conf.json
func ReadConfig() (*AssistantConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf_assistant.json")
	defer file.Close()

	if err != nil {
		logs.GetAssistantLogger().Errorf("Open assistant configuration file error: %s", err)
		return nil, utils.Errors.NoConfigFileFound
	}

	decoder := json.NewDecoder(file)
	conf := AssistantConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetAssistantLogger().Errorf("Parsing assistant configuration file error: %s", err)
		return nil, utils.Errors.ConfigFileValidationError
	}
	logs.GetAssistantLogger().Infof("Successfully loaded assistant configuration file")
	return &conf, nil
}

// InitConfig initialize configuration
func InitConfig() {
	var err error

	if assistantConfig != nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	assistantConfig, err = ReadConfig()
	if assistantConfig == nil || err != nil {
		time.Sleep(time.Second * 10)
		assistantConfig, err = ReadConfig()
		// default value
		assistantConfig = &AssistantConfig{}
		assistantConfig.Enable = true
		assistantConfig.Endpoint = config.GetConfig().Endpoint
	} else {
		assistantConfig.Enable = true
		return
	}
}

// GetConfig get wrapper
func GetConfig() *AssistantConfig {
	if assistantConfig == nil {
		InitConfig()
	}

	return assistantConfig
}

// UpdateConfig update assist config
func UpdateConfig(configBytes []byte) bool {
	c := AssistantConfig{}
	unmarshalErr := json.Unmarshal(configBytes, &c)
	if unmarshalErr != nil {
		logs.GetAssistantLogger().Errorf("Failed to unmarshal assistant config: %s, error is %s", string(configBytes), unmarshalErr.Error())
		return false
	}
	assistantConfig = &c
	pwd, err := os.Getwd()
	if err != nil {
		logs.GetAssistantLogger().Errorf("Failed to get current directory, so failed to persist assistant config to local file, error: %s", err.Error())
		return false
	}
	configPath := pwd + "/conf_assistant.json"
	err = utils.WriteStrToFile(string(configBytes), configPath)
	if err != nil {
		logs.GetAssistantLogger().Errorf("Failed to persist assistant config to local file, error: %s", err.Error())
		return false
	}
	return true
}

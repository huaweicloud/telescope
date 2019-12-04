package config

import (
	"os"
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

// CESConfig is the type for heartbeat response
type CESConfig struct {
	Enable            bool
	Endpoint          string
	EnableProcessList []HbProcess `json:"enable_processes"`
	SpecifiedProcList []string    `json:"specified_procs"`
	EnablePlugin      bool
	ExternalService   string
}

// HbProcess is the type for enable process info which used in heartbeat response
type HbProcess struct {
	Pname string `json:"name"`
	Pid   int32  `json:"pid"`
}

var (
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	cesConfig *CESConfig
	mutex     sync.Mutex
)

// ReadConfig Read the config from conf.json
func ReadConfig() (*CESConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf_ces.json")
	defer file.Close()

	if err != nil {
		logs.GetCesLogger().Errorf("Open ces configuration file error: %s", err)
		return nil, utils.Errors.NoConfigFileFound
	}

	decoder := json.NewDecoder(file)
	conf := CESConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetCesLogger().Errorf("Parsing ces configuration file error: %s", err)
		return nil, utils.Errors.ConfigFileValidationError
	}
	logs.GetCesLogger().Infof("Successfully loaded ces configuration file")
	return &conf, nil
}

// InitConfig initialize configuration
func InitConfig() {
	var err error
	if cesConfig == nil || err != nil {
		mutex.Lock()
		defer mutex.Unlock()

		cesConfig, err = ReadConfig()
		for {
			if cesConfig == nil || err != nil {
				logs.GetCesLogger().Errorf(err.Error())
				time.Sleep(time.Second * 10)
				cesConfig, err = ReadConfig()
			} else {
				cesConfig.Enable = true
				return
			}
		}
	}
}

// GetConfig get wrapper
func GetConfig() *CESConfig {
	if cesConfig == nil {
		InitConfig()
	}

	return cesConfig
}

// ReloadConfig reload config to support hot load config file
func ReloadConfig() *CESConfig {
	originalEnable := cesConfig.Enable
	originalEnableProcessList := cesConfig.EnableProcessList
	originalSpecifiedProcList := cesConfig.SpecifiedProcList
	newCesConfig, err := ReadConfig()
	for {
		if newCesConfig == nil || err != nil {
			logs.GetCesLogger().Errorf("Reload ces config error is %s", err.Error())
			time.Sleep(time.Second * 10)
			newCesConfig, err = ReadConfig()
		} else {
			newCesConfig.Enable = originalEnable
			newCesConfig.EnableProcessList = originalEnableProcessList
			newCesConfig.SpecifiedProcList = originalSpecifiedProcList
			cesConfig = newCesConfig
			return cesConfig
		}
	}
	newCesConfig.Enable = originalEnable
	cesConfig = newCesConfig
	return cesConfig
}

// UpdateConfig update ces config
func UpdateConfig(configBytes []byte) (success bool) {
	config := CESConfig{}
	unmarshalErr := json.Unmarshal(configBytes, &config)
	if unmarshalErr != nil {
		logs.GetCesLogger().Errorf("Failed to unmarshal ces config: %s, error is %s", string(configBytes), unmarshalErr.Error())
		return false
	}
	cesConfig = &config
	pwd, err := os.Getwd()
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to get current directory, so failed to persist ces config to local file, error: %s", err.Error())
		return false
	}
	configPath := pwd + "/conf_ces.json"
	err = utils.WriteStrToFile(string(configBytes), configPath)
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to persist ces config to local file, error: %s", err.Error())
		return false
	}
	return true
}

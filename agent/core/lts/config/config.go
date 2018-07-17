package config

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/errs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type TopicConfig struct {
	LogTopicId                  string
	Path                        string
	TimeExtractMode             string
	TimeExtractPattern          string
	SingleLineLog               bool
	MultiLineLogMarkTimePattern string
	MultiLineLogMarkRegex       string
	IsOsLog                     bool
}

type GroupConfig struct {
	GroupId string
	Topics  []TopicConfig
}

type LTSConfig struct {
	Enable   bool `json:"-"`
	Endpoint string
	Groups   []GroupConfig
}

//lts config from agent server
type LTSConfigRemote struct {
	Groups []GroupConfig
}

var (
	ltsConfig *LTSConfig
	mutex     sync.Mutex
)

// Read the config from conf.json
func ReadConfig() (*LTSConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf_lts.json")
	defer file.Close()

	if err != nil {
		logs.GetLtsLogger().Errorf("Open lts configuration file error: %s", err)
		return nil, utils.Errors.NoConfigFileFound
	}

	decoder := json.NewDecoder(file)
	conf := LTSConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetLtsLogger().Errorf("Parsing lts configuration file error: %s", err)
		return nil, utils.Errors.ConfigFileValidationError
	}

	logs.GetLtsLogger().Debugf("Successfully loaded lts configuration file")
	return &conf, nil
}

// Initialize configuration
func InitConfig() {
	var err error
	if ltsConfig == nil || err != nil {
		mutex.Lock()
		defer mutex.Unlock()

		ltsConfig, err = ReadConfig()
		for {
			if ltsConfig == nil || err != nil {
				logs.GetLtsLogger().Errorf(err.Error())
				time.Sleep(time.Second * 10)
				ltsConfig, err = ReadConfig()
			} else {
				return
			}
		}
	}

}

// Get wrapper
func GetConfig() *LTSConfig {
	if ltsConfig == nil {
		InitConfig()
	}

	return ltsConfig
}

//Reload config to support hot load config file; if config file update,return true, otherwise return false
func ReloadConfig() bool {
	originalEnable := ltsConfig.Enable
	newLtsConfig, err := ReadConfig()
	for {
		if newLtsConfig == nil || err != nil {
			logs.GetLtsLogger().Errorf("Reload lts config error is %s.", err.Error())
			time.Sleep(time.Second * 10)
			newLtsConfig, err = ReadConfig()
		} else {
			break
		}
	}
	originalConfigBytes, oriMarshalErr := json.Marshal(ltsConfig)
	newConfigBytes, newMarshalErr := json.Marshal(newLtsConfig)
	if oriMarshalErr == nil && newMarshalErr == nil && strings.Compare(string(originalConfigBytes), string(newConfigBytes)) != 0 {
		newLtsConfig.Enable = originalEnable
		ltsConfig = newLtsConfig
		return true
	}
	return false
}

//Update Lts config
func UpdateConfig(config LTSConfig) (success bool) {
	configBytes, marshalErr := json.Marshal(config)
	if marshalErr != nil {
		logs.GetLtsLogger().Errorf("Failed to marshal new config, error is ", marshalErr.Error())
		return false
	}
	pwd := logs.GetCurrentDirectory()
	configPath := pwd + "/conf_lts.json"
	err := utils.WriteStrToFile(string(configBytes), configPath)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to persist lts config to local file, error: %s", err.Error())
		errs.PutLtsDetail(errs.ERR_WRITE_LTS_CONFIG_FILE.Code, errs.ERR_WRITE_LTS_CONFIG_FILE.Message)
		return false
	}
	return true
}

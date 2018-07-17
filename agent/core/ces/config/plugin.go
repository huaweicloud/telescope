package config

import (
	"encoding/json"
	"os"
	"sync"

	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// PluginConfig is the type for plugins config file
type PluginConfig struct {
	Plugins []*EachPluginConfig `json:"plugins"`
}

// EachPluginConfig is the type for each plugin config
type EachPluginConfig struct {
	Path     string `json:"path"`
	Crontime int    `json:"crontime"`
}

var (
	pluginConfig *PluginConfig
	pmutex       sync.Mutex
)

// ReadPluginConfig Read the config from ../plugins/conf.json
func ReadPluginConfig() (*PluginConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/../" + ces_utils.PluginConf)
	defer file.Close()

	if err != nil {
		logs.GetCesLogger().Errorf("Open ces plugins configuration file error: %s", err)
		return nil, utils.Errors.NoConfigFileFound
	}

	decoder := json.NewDecoder(file)
	conf := PluginConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetCesLogger().Errorf("Parsing ces plugins configuration file error: %s", err)
		return nil, utils.Errors.ConfigFileValidationError
	}
	logs.GetCesLogger().Infof("Successfully loaded ces plugins configuration file")
	return &conf, nil
}

// InitPluginConfig initialize plugins configuration
func InitPluginConfig() {
	var err error
	if pluginConfig == nil || err != nil {
		pmutex.Lock()
		defer pmutex.Unlock()

		pluginConfig, err = ReadPluginConfig()

		if pluginConfig == nil || err != nil {
			logs.GetCesLogger().Errorf(err.Error())
			return
		}
	}
}

// GetPluginConfig get wrapper
func GetPluginConfig() *PluginConfig {
	if pluginConfig == nil {
		InitPluginConfig()
	}

	return pluginConfig
}

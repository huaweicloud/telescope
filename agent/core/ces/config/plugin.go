package config

import (
	"os"
	"sync"

	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// PluginConfig is the type for plugins config file
type PluginConfig struct {
	Plugins []*EachPluginConfig `json:"plugins"`
}

// EachPluginConfig is the type for each plugin config
type EachPluginConfig struct {
	Path              string `json:"path"`
	Crontime          int    `json:"crontime"`
	Type              string `json:"type"`
	MaxTimeoutProcNum int    `json:"max_timeout_proc_num"`
}

var (
	pluginConfig               *PluginConfig
	defaultPluginConfigs       []*EachPluginConfig
	customMonitorPluginConfigs []*EachPluginConfig
	eventPluginConfigs         []*EachPluginConfig
	pmutex                     sync.Mutex
)

// ReadPluginConfig Read the config from ../plugins/conf.json
func ReadPluginConfig() (*PluginConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/../" + cesUtils.PluginConf)
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

	for i, p := range conf.Plugins {
		logs.GetCesLogger().Debugf("Plugin config[%d] is %v", i, *p)

		switch p.Type {
		case cesUtils.DefaultPluginType:
			fallthrough
		case cesUtils.AgtPluginType:
			if p.Crontime < cesUtils.DefaultPluginCronTime {
				logs.GetCesLogger().Warnf("Plugin(%v) crontime is %v, less than the default crontime %v seconds. Use default crontime.", *p, p.Crontime, cesUtils.DefaultPluginCronTime)
				p.Crontime = cesUtils.DefaultPluginCronTime
			}
		case cesUtils.CustomMonitorPluginType:
			if p.Crontime < cesUtils.DefaultCustomMonitorPluginCronTime {
				logs.GetCesLogger().Warnf("Plugin(%v) crontime is %v, less than the default crontime %v seconds. Use default crontime.", *p, p.Crontime, cesUtils.DefaultCustomMonitorPluginCronTime)
				p.Crontime = cesUtils.DefaultCustomMonitorPluginCronTime
			}
		case cesUtils.EventPluginType:
			if p.Crontime < cesUtils.DefaultEventPluginCronTime {
				logs.GetCesLogger().Warnf("Plugin(%v) crontime is %v, less than the default crontime %v seconds. Use default crontime.", *p, p.Crontime, cesUtils.DefaultEventPluginCronTime)
				p.Crontime = cesUtils.DefaultEventPluginCronTime
			}
		}

		if p.MaxTimeoutProcNum == 0 ||
			p.MaxTimeoutProcNum > cesUtils.DefaultMaxTimeoutProcNum {
			logs.GetCesLogger().Warnf("Plugin(%v) max_timeout_proc_num is: %d, set it to default value(%d)", *p, p.MaxTimeoutProcNum, cesUtils.DefaultMaxTimeoutProcNum)
			p.MaxTimeoutProcNum = cesUtils.DefaultMaxTimeoutProcNum
		}

	}
	return &conf, nil
}

// InitPluginConfig initialize plugins configuration
func InitPluginConfig() {
	var err error

	if pluginConfig != nil {
		return
	}

	pmutex.Lock()
	defer pmutex.Unlock()

	pluginConfig, err = ReadPluginConfig()
	if err != nil || pluginConfig == nil {
		return
	}
}

// GetPluginConfig get wrapper
func GetPluginConfig() *PluginConfig {
	if pluginConfig == nil {
		InitPluginConfig()
	}

	return pluginConfig
}

// GetDefaultPluginConfig ...
func GetDefaultPluginConfig() []*EachPluginConfig {
	if defaultPluginConfigs == nil {
		plgConfig := GetPluginConfig()
		for _, plgConfig := range plgConfig.Plugins {
			if plgConfig.Type == cesUtils.DefaultPluginType ||
				plgConfig.Type == cesUtils.AgtPluginType {
				defaultPluginConfigs = append(defaultPluginConfigs, plgConfig)
			}
		}
	}

	return defaultPluginConfigs
}

// GetCustomMonitorPluginConfig ...
func GetCustomMonitorPluginConfig() []*EachPluginConfig {
	if customMonitorPluginConfigs == nil {
		plgConfig := GetPluginConfig()
		for _, plgConfig := range plgConfig.Plugins {
			if plgConfig.Type == cesUtils.CustomMonitorPluginType {
				customMonitorPluginConfigs = append(customMonitorPluginConfigs, plgConfig)
				logs.GetCesLogger().Debugf("Custom monitor plugin config is: %v", *plgConfig)
			}
		}
	}

	return customMonitorPluginConfigs
}

// GetEventPluginConfig ...
func GetEventPluginConfig() []*EachPluginConfig {
	if eventPluginConfigs == nil {
		plgConfig := GetPluginConfig()
		for _, plgConfig := range plgConfig.Plugins {
			if plgConfig.Type == cesUtils.EventPluginType {
				eventPluginConfigs = append(eventPluginConfigs, plgConfig)
				logs.GetCesLogger().Debugf("Event plugin config is: %v", *plgConfig)
			}
		}
	}

	return eventPluginConfigs
}

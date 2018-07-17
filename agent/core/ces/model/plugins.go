package model

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// PluginScheduler is the type for plugin scheduler
type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *config.EachPluginConfig
}

// NewPluginScheduler create a plugin scheduler by a plugin config
func NewPluginScheduler(p *config.EachPluginConfig) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	if p.Crontime < ces_utils.DefaultPluginCronTime {
		logs.GetCesLogger().Errorf("Plugin crontime is %v, less than the default crontime %v seconds. Use default crontime.", p.Crontime, ces_utils.DefaultPluginCronTime)
		p.Crontime = ces_utils.DefaultPluginCronTime
	}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Crontime) * time.Second)
	return &scheduler
}

// Schedule cron job for plugin collector
func (ps *PluginScheduler) Schedule(data chan *InputMetric) {

	for {
		select {
		case <-ps.Ticker.C:
			go func() {
				pluginData := PluginCmd(ps.Plugin)
				if pluginData != nil {
					data <- pluginData
				}

			}()

		}
	}

}

// PluginCmd output the plugin metric data by a plugin config
func PluginCmd(plugin *config.EachPluginConfig) *InputMetric {

	var result InputMetric

	if !utils.IsFileExist(plugin.Path) {
		logs.GetCesLogger().Errorf("Plugin not exist: %s", plugin.Path)
		return nil
	}

	cmd := exec.Command(plugin.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd StdoutPipe error: %v", err)
		return nil
	}
	defer stdout.Close()
	defer cmd.Wait()
	if err := cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd Start error: %v", err)
		return nil
	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin read all stdout error: %v", err)
		return nil
	}

	err = json.Unmarshal(opBytes, &result)
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin unmarshal result error: %v", err)
		return nil
	}
	return &result
}

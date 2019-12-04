package model

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// PluginScheduler is the type for plugin scheduler
type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *config.EachPluginConfig
}

// NewPluginScheduler create a plugin scheduler by a plugin config
func NewPluginScheduler(p *config.EachPluginConfig) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
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

	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.GetCesLogger().Errorf("Get current work path error: %v", err)
	}

	cmd := exec.Command(plugin.Path)
	cmd.Dir = workDir
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd StdoutPipe error: %v", err)
		return nil
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd Start error: %v", err)
		return nil
	}

	done := make(chan error, 1)
	go func() {
		defer func() {
			done <- cmd.Wait()
		}()
		opBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			logs.GetCesLogger().Errorf("Plugin read all stdout error: %v, time: %v", err, time.Now().UnixNano())
			return
		}

		if len(opBytes) == 0 {
			logs.GetCesLogger().Warn("Plugin read all stdout but get empty")
			return
		}

		err = json.Unmarshal(opBytes, &result)
		if err != nil {
			logs.GetCesLogger().Errorf("Plugin unmarshal result error: %v", err)
			return
		}

		logs.GetCesLogger().Debugf("Plugin output is: %v", result)
	}()

	timeout := plugin.MaxTimeoutProcNum * plugin.Crontime
	pid := cmd.Process.Pid
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		logs.GetCesLogger().Warnf("Plugin(%v) has not returned output for %d seconds, so kill it(PID:%d).", *plugin, timeout, pid)
		if err := cmd.Process.Kill(); err != nil {
			logs.GetCesLogger().Errorf("Failed to kill plugin(PID:%d), error is: %v", pid, err)
		} else {
			logs.GetCesLogger().Infof("Kill plugin(PID:%d) successfully", pid)
			cmd.Wait()
		}
		return nil
	case err := <-done:
		if err != nil {
			logs.GetCesLogger().Errorf("Plugin(PID:%d) process finished with error: %v", pid, err)
			return nil
		}
		logs.GetCesLogger().Infof("Plugin(PID:%d) process finished successfully", pid)
		return &result
	}
}

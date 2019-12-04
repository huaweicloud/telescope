package model

import (
	"io/ioutil"
	"os/exec"
	"time"

	"os"
	"path/filepath"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// CustomMonitorPluginScheduler ...
type CustomMonitorPluginScheduler PluginScheduler

// NewCustomMonitorPluginScheduler create a plugin scheduler by a plugin config
func NewCustomMonitorPluginScheduler(p *config.EachPluginConfig) *CustomMonitorPluginScheduler {
	scheduler := CustomMonitorPluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Crontime) * time.Second)
	return &scheduler
}

// Schedule cron job for plugin collector
func (ps *CustomMonitorPluginScheduler) Schedule(data chan CesMetricDataArr) {

	for {
		select {
		case <-ps.Ticker.C:
			go func() {
				pluginData := CustomMonitorPluginCmd(ps.Plugin)
				if pluginData != nil {
					data <- pluginData
				}
			}()
		}
	}
}

// CustomMonitorPluginCmd output the plugin metric data by a plugin config
func CustomMonitorPluginCmd(plugin *config.EachPluginConfig) CesMetricDataArr {
	var result CesMetricDataArr

	if !utils.IsFileExist(plugin.Path) {
		logs.GetCesLogger().Errorf("Custom monitor plugin not exist: %s", plugin.Path)
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
		logs.GetCesLogger().Errorf("Custom monitor plugin execute cmd StdoutPipe error: %v", err)
		return nil
	}

	if err := cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Custom monitor plugin execute cmd Start error: %v", err)
		return nil
	}
	pID := cmd.Process.Pid

	done := make(chan error, 1)
	go func() {
		defer func() {
			done <- cmd.Wait()
		}()

		opBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			logs.GetCesLogger().Errorf("Custom monitor plugin read all stdout error: %v, time: %v, PID: %d", err, time.Now().UnixNano(), pID)
			return
		}
		logs.GetCesLogger().Debugf("Custom monitor plugin original output is: %v, PID: %d", string(opBytes), pID)

		err = json.Unmarshal(opBytes, &result)
		if err != nil {
			logs.GetCesLogger().Errorf("Custom monitor plugin unmarshal error: %v, PID: %d", err, pID)
			return
		}

		logs.GetCesLogger().Debugf("Custom monitor plugin unmarshal output is: %v, PID: %d", result, pID)
	}()

	timeout := plugin.MaxTimeoutProcNum * plugin.Crontime
	pid := cmd.Process.Pid
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		logs.GetCesLogger().Warnf("Custom monitor plugin(%v) has not returned output for %d seconds, so kill it(PID:%d).", *plugin, timeout, pid)
		if err := cmd.Process.Kill(); err != nil {
			logs.GetCesLogger().Errorf("Failed to kill custom monitor plugin(PID:%d), error is: %v", pid, err)
		} else {
			logs.GetCesLogger().Infof("Kill custom monitor plugin(PID:%d) successfully", pid)
			cmd.Wait()
		}
		return nil
	case err := <-done:
		if err != nil {
			logs.GetCesLogger().Errorf("Custom monitor plugin(PID:%d) process finished with error: %v", pid, err)
			return nil
		}
		logs.GetCesLogger().Infof("Custom monitor plugin(PID:%d) process finished successfully", pid)
		return result
	}
}

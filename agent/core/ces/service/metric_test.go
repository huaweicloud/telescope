package services

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/collectors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/process"
	"github.com/yougg/assert"
	"github.com/yougg/mockfn"
)

// StartMetricCollectTask
func TestStartMetricCollectTask(t *testing.T) {
	// go InitConfig()
	t.Log("Test_StartMetricCollectTask")
	defer mockfn.RevertAll()
	mockfn.Replace(time.Sleep, func(d time.Duration) {})
	mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
		return &utils.GeneralConfig{
			ProjectId: "project_id",
		}
	})
	mockfn.Replace((*collectors.DiskCollector).Collect, func(c *collectors.DiskCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.NetCollector).Collect, func(c *collectors.NetCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.LoadCollector).Collect, func(c *collectors.LoadCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.ProcStatusCollector).Collect, func(c *collectors.ProcStatusCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	// time.NewTicker
	mockfn.Replace(time.NewTicker, func(d time.Duration) *time.Ticker {
		c := make(chan time.Time)
		go func() {
			c <- time.Time{}
			close(c)
		}()
		t := &time.Ticker{
			C: c,
		}
		return t
	})
	mockfn.Replace(config.GetConfig, func() *config.CESConfig {
		return &config.CESConfig{
			Enable: true,
			EnableProcessList: []config.HbProcess{
				{Pid: 1},
				{Pid: 2},
			},
		}
	})
	// ProcessCollector
	mockfn.Replace((*collectors.ProcessCollector).Collect, func(c *collectors.ProcessCollector, collectTime int64) *model.InputMetric {
		t.SkipNow()
		return &model.InputMetric{}
	})
	// process.PidExists
	mockfn.Replace(process.PidExists, func(pid int32) (bool, error) {
		return true, nil
	})
	mockfn.Replace(logs.GetCurrentDirectory, func() string {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}
		dir, err := filepath.Abs(filename)
		if nil != err {
			panic("get file absolute path error")
		}
		dir = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(dir))))
		t.Log(dir)
		return dir
	})
	mockfn.Replace(report.SendMetricData, func(url string, data *model.InputMetric, isAggregate bool) {
		t.Log(url)
		t.Logf("%#v", data)
	})
	StartMetricCollectTask()
}

func TestStartMetricCollectTask2(t *testing.T) {
	// go InitConfig()
	t.Log("Test_StartMetricCollectTask")
	defer mockfn.RevertAll()
	mockfn.Replace(time.Sleep, func(d time.Duration) {
		return
	})
	mockfn.Replace(config.GetConfig, func() *config.CESConfig {
		return &config.CESConfig{}
	})
	mockfn.Replace((*collectors.DiskCollector).Collect, func(c *collectors.DiskCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.NetCollector).Collect, func(c *collectors.NetCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.LoadCollector).Collect, func(c *collectors.LoadCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	mockfn.Replace((*collectors.ProcStatusCollector).Collect, func(c *collectors.ProcStatusCollector, collectTime int64) *model.InputMetric {
		return &model.InputMetric{}
	})
	// time.NewTicker
	mockfn.Replace(time.NewTicker, func(d time.Duration) *time.Ticker {
		c := make(chan time.Time)
		go func() {
			c <- time.Time{}
			c <- time.Time{}
			c <- time.Time{}
			c <- time.Time{}
			c <- time.Time{}
			c <- time.Time{}
			c <- time.Time{}
			close(c)
		}()
		t := &time.Ticker{
			C: c,
		}
		return t
	})
	mockfn.Replace(config.GetConfig, func() *config.CESConfig {
		return &config.CESConfig{
			Enable:            true,
			SpecifiedProcList: []string{"1", "2"},
		}
	})
	// ProcessCollector
	mockfn.Replace((*collectors.SpeProcCountCollector).Collect, func(c *collectors.SpeProcCountCollector, collectTime int64) *model.InputMetric {
		t.SkipNow()
		return &model.InputMetric{}
	})
	StartMetricCollectTask()
}

// BuildURL
func TestBuildURL(t *testing.T) {
	// go InitConfig()
	t.Log("Test_BuildURL")
	defer mockfn.RevertAll()
	mockfn.Replace(config.GetConfig, func() *config.CESConfig {
		return &config.CESConfig{
			Endpoint: "end",
		}
	})
	// utils.GetConfig
	mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
		return &utils.GeneralConfig{
			ProjectId: "id",
		}
	})
	url := BuildURL("url")
	a := assert.New(t)
	a.Equal("end/V1.0/idurl", url)
}

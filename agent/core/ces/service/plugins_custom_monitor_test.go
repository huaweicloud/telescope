package services

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//CollectCustomMonitorPluginTask
func TestCollectCustomMonitorPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_CollectCustomMonitorPluginTask", t, func() {
		Convey("test case 1", func() {
			Convey("test case 1", func() {
				defer mockfn.RevertAll()
				mockfn.Replace(config.GetConfig, func() *config.CESConfig {
					return &config.CESConfig{}
				})
				CollectCustomMonitorPluginTask(nil)
			})
			Convey("test case 2", func() {
				defer mockfn.RevertAll()
				mockfn.Replace(config.GetConfig, func() *config.CESConfig {
					return &config.CESConfig{
						Enable:       true,
						EnablePlugin: true,
					}
				})
				mockfn.Replace(config.GetCustomMonitorPluginConfig, func() []*config.EachPluginConfig {
					return nil
				})
				CollectCustomMonitorPluginTask(nil)
			})
			Convey("test case 3", func() {
				defer mockfn.RevertAll()
				mockfn.Replace(config.GetConfig, func() *config.CESConfig {
					return &config.CESConfig{
						Enable:       true,
						EnablePlugin: true,
					}
				})
				mockfn.Replace(config.GetCustomMonitorPluginConfig, func() []*config.EachPluginConfig {
					configs := []*config.EachPluginConfig{
						{Path: "1"},
						{Path: "2"},
						{Path: "3"},
					}
					return configs
				})
				mockfn.Replace(model.NewCustomMonitorPluginScheduler, func(p *config.EachPluginConfig) *model.CustomMonitorPluginScheduler {
					return nil
				})
				CollectCustomMonitorPluginTask(nil)
			})
			Convey("test case 4", func() {
				defer mockfn.RevertAll()
				mockfn.Replace(config.GetConfig, func() *config.CESConfig {
					return &config.CESConfig{
						Enable:       true,
						EnablePlugin: true,
					}
				})
				mockfn.Replace(config.GetCustomMonitorPluginConfig, func() []*config.EachPluginConfig {
					configs := []*config.EachPluginConfig{
						{Path: "1"},
						{Path: "2"},
						{Path: "3"},
					}
					return configs
				})
				mockfn.Replace(model.NewCustomMonitorPluginScheduler, func(p *config.EachPluginConfig) *model.CustomMonitorPluginScheduler {
					c := make(chan time.Time, 8)
					c <- time.Time{}

					t := &time.Ticker{
						C: c,
					}
					return &model.CustomMonitorPluginScheduler{
						Ticker: t,
					}
				})
				mockfn.Replace((*model.CustomMonitorPluginScheduler).Schedule, func(m *model.CustomMonitorPluginScheduler, data chan model.CesMetricDataArr) {
					return
					t.SkipNow()
				})
				metrics := make(chan model.CesMetricDataArr, 1)
				metrics <- model.CesMetricDataArr{}
				CollectCustomMonitorPluginTask(metrics)
			})
		})
	})
}

//SendCustomMonitorPluginTask
func TestSendCustomMonitorPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendCustomMonitorPluginTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(BuildURL, func(destURI string) string {
				t.SkipNow()
				return ""
			})
			mockfn.Replace((*model.PluginScheduler).Schedule, func(m *model.PluginScheduler, data chan *model.InputMetric) {
				return
				t.SkipNow()
			})
			metrics := make(chan model.CesMetricDataArr, 1)
			metrics <- model.CesMetricDataArr{}
			SendCustomMonitorPluginTask(metrics)
		})
	})
}

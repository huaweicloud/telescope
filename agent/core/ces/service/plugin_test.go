package services

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//CollectPluginTask
func TestCollectPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_CollectPluginTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{}
			})
			CollectPluginTask(nil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetDefaultPluginConfig, func() []*config.EachPluginConfig {
				return nil
			})
			CollectPluginTask(nil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetDefaultPluginConfig, func() []*config.EachPluginConfig {
				configs := []*config.EachPluginConfig{
					{Path: "1"},
					{Path: "2"},
					{Path: "3"},
				}
				return configs
			})
			mockfn.Replace(model.NewPluginScheduler, func(p *config.EachPluginConfig) *model.PluginScheduler {
				return nil
			})
			CollectPluginTask(nil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetDefaultPluginConfig, func() []*config.EachPluginConfig {
				configs := []*config.EachPluginConfig{
					{Path: "1"},
					{Path: "2"},
					{Path: "3"},
				}
				return configs
			})
			mockfn.Replace(model.NewPluginScheduler, func(p *config.EachPluginConfig) *model.PluginScheduler {
				c := make(chan time.Time, 8)
				c <- time.Time{}

				t := &time.Ticker{
					C: c,
				}
				return &model.PluginScheduler{
					Ticker: t,
				}
			})
			mockfn.Replace((*model.PluginScheduler).Schedule, func(m *model.PluginScheduler, data chan *model.InputMetric) {
				return
				t.SkipNow()
			})
			metrics := make(chan *model.InputMetric, 1)
			metrics <- &model.InputMetric{}
			//CollectPluginTask(metrics)
		})
	})
}

//SendPluginTask
func TestSendPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendPluginTask", t, func() {
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
			metrics := make(chan *model.InputMetric, 1)
			metrics <- &model.InputMetric{}
			SendPluginTask(metrics)
		})
	})
}

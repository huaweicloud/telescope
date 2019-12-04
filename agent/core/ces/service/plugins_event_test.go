package services

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//CollectEventPluginTask
func TestCollectEventPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_CollectEventPluginTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{}
			})
			CollectEventPluginTask(nil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetEventPluginConfig, func() []*config.EachPluginConfig {
				return nil
			})
			mockfn.Replace(model.NewEventPluginScheduler, func(p *config.EachPluginConfig) *model.EventPluginScheduler {
				return nil
			})
			CollectEventPluginTask(nil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetEventPluginConfig, func() []*config.EachPluginConfig {
				configs := []*config.EachPluginConfig{
					{Path: "1"},
					{Path: "2"},
					{Path: "3"},
				}
				return configs
			})
			mockfn.Replace(model.NewEventPluginScheduler, func(p *config.EachPluginConfig) *model.EventPluginScheduler {
				return nil
			})
			CollectEventPluginTask(nil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{
					Enable:       true,
					EnablePlugin: true,
				}
			})
			mockfn.Replace(config.GetEventPluginConfig, func() []*config.EachPluginConfig {
				configs := []*config.EachPluginConfig{
					{Path: "1"},
					{Path: "2"},
					{Path: "3"},
				}
				return configs
			})
			mockfn.Replace(model.NewEventPluginScheduler, func(p *config.EachPluginConfig) *model.EventPluginScheduler {
				c := make(chan time.Time, 8)
				c <- time.Time{}
				c <- time.Time{}

				t := &time.Ticker{
					C: c,
				}
				return &model.EventPluginScheduler{
					Ticker: t,
				}
			})
			/*	mockfn.Replace((*model.EventPluginScheduler).Schedule, func(m *model.EventPluginScheduler, data chan model.CesEventDataArr) {
				return
				t.SkipNow()
			})*/
			/*		metrics := make(chan model.CesEventDataArr, 1)
					metrics <- model.CesEventDataArr{}
					CollectEventPluginTask(metrics)*/
		})
	})
}

//SendEventPluginTask
func TestSendEventPluginTask(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendEventPluginTask", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(BuildURL, func(destURI string) string {
				t.SkipNow()
				return ""
			})
			metrics := make(chan model.CesEventDataArr, 1)
			metrics <- model.CesEventDataArr{}
			SendEventPluginTask(metrics)
		})
	})
}

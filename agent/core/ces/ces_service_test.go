package ces

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/buger/jsonparser"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/json-iterator/go"
	"github.com/shirou/gopsutil/process"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//Init
func TestInit(t *testing.T) {
	//go InitConfig()
	Convey("Test_Init", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.InitConfig, func() {
				return
			})
			mockfn.Replace(config.InitPluginConfig, func() {
				return
			})
			mockfn.Replace(initchRawData, func() {
				return
			})
			mockfn.Replace(initchAgRawData, func() {
				return
			})
			mockfn.Replace(initchAgResult, func() {
				return
			})
			mockfn.Replace(initchProcessInfo, func() {
				return
			})
			mockfn.Replace(initchPluginData, func() {
				return
			})
			mockfn.Replace(initEnvVariables, func() {
				return
			})
			service := Service{}
			service.Init()
		})
	})
}

//Start
func TestStart(t *testing.T) {
	//go InitConfig()
	Convey("Test_Start", t, func() {
		Convey("test case 1", func() {
			service := Service{}
			service.Start()
		})
	})
}

//updateConfig
func TestUpdateConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_updateConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(jsonparser.GetBoolean, func(data []byte, keys ...string) (val bool, err error) {
				return true, nil
			})
			mockfn.Replace(config.GetConfig, func() *config.CESConfig {
				return &config.CESConfig{Enable: false, EnableProcessList: []config.HbProcess{{Pid: 1}}}
			})
			mockfn.Replace(jsonparser.GetString, func(data []byte, keys ...string) (val string, err error) {
				//t.SkipNow()
				return "", nil
			})
			mockfn.Replace(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal, func(data []byte, v interface{}) error {
				//t.SkipNow()
				return nil
			})
			mockfn.Replace(process.PidExists, func(pid int32) (bool, error) {
				//t.SkipNow()
				return false, nil
			})
			configChan := channel.GetCesConfigChan()
			configChan <- "123"
			go updateConfig()
			time.Sleep(time.Second * 2)
		})
	})
}

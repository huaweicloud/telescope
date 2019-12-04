package ces

import (
	. "github.com/smartystreets/goconvey/convey"

	"errors"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/yougg/mockfn"
	"os"
	"testing"
)

//initchRawData
func TestInitchRawData(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchRawData", t, func() {
		Convey("test case 1", func() {
			initchRawData()
		})
	})
}

///getchRawData
func TestGetchRawData(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchRawData", t, func() {
		Convey("test case 1", func() {
			chRawData = nil
			data := getchRawData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chRawData = make(chan *model.InputMetric, 100)
			data := getchRawData()
			So(data, ShouldNotBeNil)
		})
	})
}

//initchPluginData
func TestInitchPluginData(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchPluginData", t, func() {
		Convey("test case 1", func() {
			initchPluginData()
		})
	})
}

//getchPluginData
func TestGetchPluginData(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchPluginData", t, func() {
		Convey("test case 1", func() {
			chPluginData = nil
			data := getchPluginData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chPluginData = make(chan *model.InputMetric, 100)
			data := getchPluginData()
			So(data, ShouldNotBeNil)
		})
	})
}

//initchAgRawData
func TestInitchAgRawData(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchAgRawData", t, func() {
		Convey("test case 1", func() {
			initchAgRawData()
		})
	})
}

//getchAgRawData
func TestGetchAgRawData(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchAgRawData", t, func() {
		Convey("test case 1", func() {
			chAgRawData = nil
			data := getchAgRawData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chAgRawData = make(chan model.InputMetricSlice, 100)
			data := getchAgRawData()
			So(data, ShouldNotBeNil)
		})
	})
}

//initchAgResult
func TestInitchAgResult(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchAgResult", t, func() {
		Convey("test case 1", func() {
			initchAgResult()
		})
	})
}

//getchAgResult
func TestGetchAgResult(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchAgResult", t, func() {
		Convey("test case 1", func() {
			chAgResult = nil
			data := getchAgResult()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chAgResult = make(chan *model.InputMetric, 100)
			data := getchAgResult()
			So(data, ShouldNotBeNil)
		})
	})
}

//initchProcessInfo
func TestInitchProcessInfo(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchProcessInfo", t, func() {
		Convey("test case 1", func() {
			initchProcessInfo()
		})
	})
}

//getchProcessInfo
func TestGetchProcessInfo(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchProcessInfo", t, func() {
		Convey("test case 1", func() {
			chProcessInfo = nil
			data := getchProcessInfo()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chProcessInfo = make(chan model.ChProcessList, 10)
			data := getchProcessInfo()
			So(data, ShouldNotBeNil)
		})
	})
}

//initchSpeProcData
func TestInitchSpeProcData(t *testing.T) {
	//go InitConfig()
	Convey("Test_initchSpeProcData", t, func() {
		Convey("test case 1", func() {
			initchSpeProcData()
		})
	})
}

//getchSpeProcData
func TestGetchSpeProcData(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchSpeProcData", t, func() {
		Convey("test case 1", func() {
			chSpeProcData = nil
			data := getchSpeProcData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chSpeProcData = make(chan *model.InputMetric, 100)
			data := getchSpeProcData()
			So(data, ShouldNotBeNil)
		})
	})
}

//getchCustomMonitorData
func TestGetchCustomMonitorData(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetchCustomMonitorData", t, func() {
		Convey("test case 1", func() {
			chCustomMonitorData = nil
			data := getchCustomMonitorData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chCustomMonitorData = make(chan model.CesMetricDataArr, 100)
			data := getchCustomMonitorData()
			So(data, ShouldNotBeNil)
		})
	})
}

//getchEventData
func TestGetchEventData(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetchEventData", t, func() {
		Convey("test case 1", func() {
			chEventData = nil
			data := getchEventData()
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			chEventData = make(chan model.CesEventDataArr, 100)
			data := getchEventData()
			So(data, ShouldNotBeNil)
		})
	})
}

//initEnvVariables
func TestInitEnvVariables(t *testing.T) {
	//go InitConfig()
	Convey("Test_initEnvVariables", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Setenv, func(key, value string) error {
				return errors.New("")
			})
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			initEnvVariables()
		})
	})
}

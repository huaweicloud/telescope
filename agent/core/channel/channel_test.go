package channel

import (
	"testing"

	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//GetCesConfigChan
func TestGetCesConfigChan(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetCesConfigChan", t, func() {
		Convey("test case 1", func() {
			configChan := GetCesConfigChan()
			So(configChan, ShouldNotBeNil)
		})
	})
}

//GetHeartBeatChan
func TestGetHeartBeatChan(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetHeartBeatChan", t, func() {
		Convey("test case 1", func() {
			beatChan := GetHeartBeatChan()
			So(beatChan, ShouldNotBeNil)
		})
	})
}

//NewHBEntity
func Test_NewHBEntity(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewHBEntity", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			entity := NewHBEntity(0, 0, "")
			So(entity, ShouldNotBeNil)
		})
	})
}

//MapStatus
func Test_MapStatus(t *testing.T) {
	//go InitConfig()
	Convey("Test_MapStatus", t, func() {
		Convey("test case 1", func() {
			status := Running.MapStatus()
			So(status, ShouldEqual, "running")
		})
		Convey("test case 2", func() {
			status := Shutdown.MapStatus()
			So(status, ShouldEqual, "stopped")
		})
		Convey("test case 3", func() {
			status := Upgrading.MapStatus()
			So(status, ShouldEqual, "upgrading")
		})
		Convey("test case 4", func() {
			const a StatusEnum = 5
			status := a.MapStatus()
			So(status, ShouldEqual, "unknown")
		})
	})
}

//GetServicesChData
func TestGetServicesChData(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetServicesChData", t, func() {
		Convey("test case 1", func() {
			data := GetServicesChData()
			So(data, ShouldNotBeNil)
		})
	})
}

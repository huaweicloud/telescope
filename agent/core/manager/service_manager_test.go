package manager

import (
	"testing"

	"github.com/huaweicloud/telescope/agent/core/assistant"
	"github.com/huaweicloud/telescope/agent/core/ces"
	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//NewServicemanager
func TestNewServicemanager(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewServicemanager", t, func() {
		Convey("test case 1", func() {
			newServicemanager := NewServicemanager()
			So(newServicemanager, ShouldNotBeNil)
		})
	})
}

//Init
func TestInit(t *testing.T) {
	//go InitConfig()
	Convey("Test_Init", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.InitConfig, func() {
				return
			})
			service := servicemanager{}
			service.Init()
		})
	})
}

//RegisterService
func TestRegisterService(t *testing.T) {
	//go InitConfig()
	Convey("Test_RegisterService", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*ces.Service).Init, func(*ces.Service) {
				t.SkipNow()
				return
			})
			mockfn.Replace((*assistant.Assistant).Init, func(*assistant.Assistant) {
				t.SkipNow()
				return
			})
			/*	service := servicemanager{}
				service.RegisterService()*/
		})
	})
}

//InitService
func TestInitService(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitService", t, func() {
		Convey("test case 1", func() {
			service := servicemanager{}
			service.InitService()
		})
	})
}

//StartService
func TestStartService(t *testing.T) {
	//go InitConfig()
	Convey("Test_StartService", t, func() {
		Convey("test case 1", func() {
			service := servicemanager{}
			service.StartService()
		})
	})
}

//HeartBeat
func TestHeartBeat(t *testing.T) {
	//go InitConfig()
	Convey("Test_HeartBeat", t, func() {
		Convey("test case 1", func() {
			service := servicemanager{}
			service.HeartBeat()
		})
	})
}

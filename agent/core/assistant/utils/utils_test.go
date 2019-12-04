package utils

import (
	"errors"
	"fmt"
	"github.com/huaweicloud/telescope/agent/core/assistant/config"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//InitConfig
func TestInitConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
		})
	})
}

//getCSTLocation
func TestBuildURL(t *testing.T) {

	Convey("Test_InitConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(config.GetConfig, func() *config.AssistantConfig {
				assistantConfig := &config.AssistantConfig{
					Endpoint: "Endpoint",
				}
				return assistantConfig
			})
			//utils.GetConfig
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				assistantConfig := &utils.GeneralConfig{
					ProjectId: "ProjectId",
				}
				return assistantConfig
			})
			url := BuildURL("123")
			So(url, ShouldEqual, "Endpoint/V1.0/ProjectId123")
		})
	})
}

//GetMarshalledRequestBody
func TestGetMarshalledRequestBody(t *testing.T) {

	Convey("Test_InitConfig", t, func() {
		Convey("test case 1", func() {
			bytes, e := GetMarshalledRequestBody("123", "")
			So(bytes, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("123")
			})
			bytes, e := GetMarshalledRequestBody(json.Marshal, "12")
			So(bytes, ShouldBeNil)
			So(e, ShouldNotBeNil)
			fmt.Println(bytes, e)
		})
	})
}

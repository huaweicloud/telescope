package config

import (
	"errors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"os"
	"testing"
	"time"
)

//InitConfig
func TestInitConfig(t *testing.T) {

	//ReadConfig
	defer mockfn.RevertAll()
	mockfn.Replace(ReadConfig, func() (*AssistantConfig, error) {
		return nil, errors.New("123")
	})
	//time.Sleep
	mockfn.Replace(time.Sleep, func(d time.Duration) {
		return
	})
	mockfn.Replace(config.GetConfig, func() *config.CESConfig {
		cesConfig := &config.CESConfig{}
		return cesConfig
	})
	//logs.GetAssistantLogger().Errorf
	mockfn.Replace(logs.GetAssistantLogger().Errorf, func(format string, params ...interface{}) error {
		return nil
	})
	assistantConfig = nil
	InitConfig()

}

//ReadConfig
func TestReadConfig(t *testing.T) {
	Convey("Test_ReadConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			readConfig, e := ReadConfig()
			So(readConfig, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, errors.New("123")
			})
			readConfig, e := ReadConfig()
			So(readConfig, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*jsoniter.Decoder).Decode, func(j *jsoniter.Decoder, obj interface{}) error {
				return nil
			})
			readConfig, e := ReadConfig()
			So(readConfig, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
	})
}

//GetConfig
func TestGetConfig(t *testing.T) {
	Convey("Test_GetConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(InitConfig, func() {

			})
			assistantConfig = nil
			readConfig := GetConfig()
			So(readConfig, ShouldBeNil)
		})
	})
}

//UpdateConfig
func TestUpdateConfig(t *testing.T) {
	Convey("Test_UpdateConfig", t, func() {
		Convey("test case 1", func() {
			assistantConfig = nil
			success := UpdateConfig(nil)
			So(success, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			//utils.WriteStrToFile
			defer mockfn.RevertAll()
			mockfn.Replace(utils.WriteStrToFile, func(str string, filepath string) error {
				return nil
			})
			config := AssistantConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeTrue)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Getwd, func() (dir string, err error) {
				return "", errors.New("123")
			})
			config := AssistantConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeFalse)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.WriteStrToFile, func(str string, filepath string) error {
				return errors.New("123")
			})
			config := AssistantConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeFalse)
		})
	})
}

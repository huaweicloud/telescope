package config

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//ReadConfig
func TestReadConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_ReadConfig", t, func() {
		Convey("test case 1", func() {
			//os.Open
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, errors.New("123")
			})
			config, e := ReadConfig()
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*jsoniter.Decoder).Decode, func(j *jsoniter.Decoder, v interface{}) error {
				return errors.New("123")
			})
			config, e := ReadConfig()
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
	})
}

//InitConfig
func TestInitConfig(t *testing.T) {

	Convey("Test_InitConfig", t, func() {
		Convey("test case 2", func() {
			//os.Open
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*CESConfig, error) {
				return nil, errors.New("123")
			})
			//time.Sleep
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				t.SkipNow()
			})
			InitConfig()
		})
		Convey("test case 1", func() {
			InitConfig()
		})
	})
}

//GetConfig
func TestGetConfig(t *testing.T) {
	Convey("Test_collectNewFile", t, func() {
		Convey("test case 2", func() {
			//InitConfig
			defer mockfn.RevertAll()
			mockfn.Replace(InitConfig, func() {

			})
			config := GetConfig()
			So(config, ShouldBeNil)
		})
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(InitConfig, func() {

			})
			cesConfig = nil
			config := GetConfig()
			So(config, ShouldBeNil)
		})
	})
}

//ReloadConfig
func TestReloadConfig(t *testing.T) {

	Convey("Test_ReloadConfig", t, func() {

		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*CESConfig, error) {
				config := &CESConfig{}
				return config, nil
			})
			mockfn.Replace(json.Marshal, func(v interface{}) ([]byte, error) {
				return []byte("a"), nil
			})
			mockfn.Replace(strings.Compare, func(a, b string) int {
				return 1
			})
			cesConfig = &CESConfig{}
			isReload := ReloadConfig()
			So(isReload, ShouldNotBeNil)
		})
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*CESConfig, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				t.SkipNow()
			})
			mockfn.Replace(json.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("123")
			})
			cesConfig = &CESConfig{}
			isReload := ReloadConfig()
			So(isReload, ShouldNotBeNil)
		})
	})
}

//UpdateConfig
func TestUpdateConfig(t *testing.T) {

	Convey("Test_UpdateConfig", t, func() {
		Convey("test case 1", func() {
			config := CESConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeTrue)
		})
		Convey("test case 2", func() {
			//ReadConfig
			defer mockfn.RevertAll()
			mockfn.Replace(json.Unmarshal, func(data []byte, v interface{}) error {
				return errors.New("123")
			})
			config := CESConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeTrue)
		})
		Convey("test case 3", func() {
			//ReadConfig
			defer mockfn.RevertAll()
			mockfn.Replace(utils.WriteStrToFile, func(str string, filepath string) error {
				return errors.New("123")
			})
			config := CESConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeFalse)
		})
		Convey("test case 4", func() {
			//ReadConfig
			defer mockfn.RevertAll()
			mockfn.Replace(os.Getwd, func() (dir string, err error) {
				return "", errors.New("123")
			})
			config := CESConfig{}
			bytes, _ := json.Marshal(config)
			success := UpdateConfig(bytes)
			So(success, ShouldBeFalse)
		})
		pwd, _ := os.Getwd()
		defer os.Remove(pwd + "/conf_ces.json")
	})
}

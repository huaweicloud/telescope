package config

import (
	"errors"
	"github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"os"
	"testing"
	"time"
)

//InitConfig
func TestInitConfig1(t *testing.T) {
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

//ReadPluginConfig
func TestReadPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_ReadPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, errors.New("123")
			})
			config, e := ReadPluginConfig()
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*jsoniter.Decoder).Decode, func(j *jsoniter.Decoder, v interface{}) error {
				//ReadPluginConfig
				mockfn.Replace(ReadPluginConfig, func() (*PluginConfig, error) {
					conf := &PluginConfig{
						Plugins: []*EachPluginConfig{
							{Path: ""},
						},
					}
					return conf, nil
				})
				return errors.New("123")
			})
			config, e := ReadPluginConfig()
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*jsoniter.Decoder).Decode, func(j *jsoniter.Decoder, v interface{}) error {
				v = &PluginConfig{
					Plugins: []*EachPluginConfig{
						{Path: ""},
					},
				}
				return nil
			})
			config, e := ReadPluginConfig()
			So(config, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
	})
}

//InitPluginConfig
func TestInitPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ReadPluginConfig, func() (*PluginConfig, error) {
				return nil, errors.New("123")
			})
			InitPluginConfig()
		})
	})
}

//GetPluginConfig
func TestGetPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(InitPluginConfig, func() {
				return
			})
			config := GetPluginConfig()
			So(config, ShouldBeNil)
		})
	})
}

//GetDefaultPluginConfig
func TestGetDefaultPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetDefaultPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetPluginConfig, func() *PluginConfig {
				config := &PluginConfig{
					Plugins: []*EachPluginConfig{
						{Path: "12"},
					},
				}
				return config
			})
			config := GetDefaultPluginConfig()
			So(config, ShouldNotBeNil)
		})
	})
}

//GetCustomMonitorPluginConfig
func TestGetCustomMonitorPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetCustomMonitorPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetPluginConfig, func() *PluginConfig {
				config := &PluginConfig{
					Plugins: []*EachPluginConfig{
						{Path: "12", Type: "Custom Monitor"},
					},
				}
				return config
			})
			config := GetCustomMonitorPluginConfig()
			So(config, ShouldNotBeNil)
		})
	})
}

//GetEventPluginConfig
func TestGetEventPluginConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetEventPluginConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetPluginConfig, func() *PluginConfig {
				config := &PluginConfig{
					Plugins: []*EachPluginConfig{
						{Path: "12", Type: "Event"},
					},
				}
				return config
			})
			config := GetEventPluginConfig()
			So(config, ShouldNotBeNil)
		})
	})
}

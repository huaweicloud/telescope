package model

import (
	. "github.com/smartystreets/goconvey/convey"

	"errors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	"github.com/yougg/mockfn"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"
)

//NewPluginScheduler
func TestNewPluginScheduler(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewPluginScheduler", t, func() {
		Convey("test case 1", func() {
			pluginConfig := &config.EachPluginConfig{Crontime: 12}
			scheduler := NewPluginScheduler(pluginConfig)
			So(scheduler, ShouldNotBeNil)
		})
	})
}

//Schedule
func TestSchedule(t *testing.T) {
	//go InitConfig()
	Convey("Test_Schedule", t, func() {
		Convey("test case 1", func() {
			/*defer mockfn.RevertAll()
			mockfn.Replace(PluginCmd, func(plugin *config.EachPluginConfig) *InputMetric {
				os.Exit(0)
				return nil
			})
			scheduler := &PluginScheduler{
				Ticker: time.NewTicker(1 * time.Second),
				Plugin: &config.EachPluginConfig{},
			}
			metrics := make(chan *InputMetric, 1)
			scheduler.Schedule(metrics)*/
		})
	})
}

//PluginCmd
func TestPluginCmd(t *testing.T) {
	//go InitConfig()
	Convey("Test_PluginCmd", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return false
			})

			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			//cmd.StdoutPipe
			mockfn.Replace((*exec.Cmd).StdoutPipe, func(*exec.Cmd) (io.ReadCloser, error) {
				return nil, errors.New("")
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, errors.New("")
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			//ioutil.ReadAll
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 6", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			//ioutil.ReadAll
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
		Convey("test case 7", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "", errors.New("")
			})
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return true
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			//ioutil.ReadAll
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			//json.Unmarshal
			mockfn.Replace(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal, func(data []byte, v interface{}) error {
				return nil
			})
			pluginConfig := &config.EachPluginConfig{}
			cmd := PluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
	})
}

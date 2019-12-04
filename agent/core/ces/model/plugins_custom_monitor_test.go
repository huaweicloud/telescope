package model

import (
	. "github.com/smartystreets/goconvey/convey"

	"errors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/yougg/mockfn"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

//NewCustomMonitorPluginScheduler
func TestNewCustomMonitorPluginScheduler(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewCustomMonitorPluginScheduler", t, func() {
		Convey("test case 1", func() {
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			scheduler := NewCustomMonitorPluginScheduler(pluginConfig)
			So(scheduler, ShouldNotBeNil)
		})
	})
}

//Schedule
func TestScheduleMonitor(t *testing.T) {
	//go InitConfig()
	Convey("Test_ScheduleMonitor", t, func() {
		Convey("test case 1", func() {

		})
	})
}

//CustomMonitorPluginCmd
func TestCustomMonitorPluginCmd(t *testing.T) {
	//go InitConfig()
	Convey("Test_CustomMonitorPluginCmd", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.IsFileExist, func(path string) bool {
				return false
			})
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			cmd := CustomMonitorPluginCmd(pluginConfig)
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
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			cmd := CustomMonitorPluginCmd(pluginConfig)
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
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			cmd := CustomMonitorPluginCmd(pluginConfig)
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
			mockfn.Replace((*exec.Cmd).StdoutPipe, func(*exec.Cmd) (io.ReadCloser, error) {
				return nil, nil
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, errors.New("")
			})
			//exec.Command
			mockfn.Replace(exec.Command, func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{
					Process: &os.Process{Pid: 1},
				}
			})
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			cmd := CustomMonitorPluginCmd(pluginConfig)
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
			mockfn.Replace((*exec.Cmd).StdoutPipe, func(*exec.Cmd) (io.ReadCloser, error) {
				return nil, nil
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			//exec.Command
			mockfn.Replace(exec.Command, func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{
					Process: &os.Process{Pid: 1},
				}
			})
			pluginConfig := &config.EachPluginConfig{Crontime: 1}
			cmd := CustomMonitorPluginCmd(pluginConfig)
			So(cmd, ShouldBeNil)
		})
	})
}

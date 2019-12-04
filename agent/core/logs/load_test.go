package logs

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//loadAppConfig
func TestLoadAppConfig(t *testing.T) {

	Convey("Test_loadAppConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getCommonLog, func() (config string) {
				return
			})
			mockfn.Replace(seelog.LoggerFromConfigAsBytes, func(data []byte) (seelog.LoggerInterface, error) {
				return nil, errors.New("")
			})
			loadAppConfig()
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getCommonLog, func() (config string) {
				return
			})
			mockfn.Replace(getCesLog, func() (config string) {
				return
			})
			mockfn.Replace(getAssistantLog, func() (config string) {
				return
			})
			mockfn.Replace(seelog.LoggerFromConfigAsBytes, func(data []byte) (seelog.LoggerInterface, error) {
				return nil, nil
			})
			loadAppConfig()
		})
	})
}

//disable
func TestDisableLog(t *testing.T) {
	Convey("Test_DisableLog", t, func() {
		Convey("test case 1", func() {
			disable()
		})
	})
}

//GetAssistantLogger
func TestGetAssistantLogger(t *testing.T) {
	Convey("Test_GetAssistantLogger", t, func() {
		Convey("test case 1", func() {
			getAssistantLogger := GetAssistantLogger()
			So(getAssistantLogger, ShouldNotBeNil)
		})
	})
}

//GetCesLogger
func TestGetCesLogger(t *testing.T) {
	Convey("Test_GetCesLogger", t, func() {
		Convey("test case 1", func() {
			getCesLogger := GetCesLogger()
			So(getCesLogger, ShouldNotBeNil)
		})
	})
}

//GetLogger
func TestGetLogger(t *testing.T) {
	Convey("Test_GetLogger", t, func() {
		Convey("test case 1", func() {
			getLogger := GetLogger()
			So(getLogger, ShouldNotBeNil)
		})
	})
}

//GetCurrentDirectory
func TestGetCurrentDirectory(t *testing.T) {
	Convey("Test_GetCurrentDirectory", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "123", nil
			})
			directory := GetCurrentDirectory()
			So(directory, ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Abs, func(path string) (string, error) {
				return "123", errors.New("123")
			})
			directory := GetCurrentDirectory()
			So(directory, ShouldEqual, "123")
		})
	})
}

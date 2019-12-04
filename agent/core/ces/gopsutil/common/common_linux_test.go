package common

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//getLsblkResult
func TestGetLsblkResult(t *testing.T) {
	//go InitConfig()
	Convey("Test_getLsblkResult", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			/*	mockfn.Replace((*exec.Cmd).StdoutPipe, func(*exec.Cmd) (io.ReadCloser, error) {
				return nil, errors.New("123")
			})*/
			result := getLsblkResult()
			So(result, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*exec.Cmd).StdoutPipe, func(*exec.Cmd) (io.ReadCloser, error) {
				return nil, errors.New("123")
			})
			result := getLsblkResult()
			So(result, ShouldBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			mockfn.Replace((*bufio.Reader).ReadString, func(b *bufio.Reader, delim byte) (string, error) {
				return "", errors.New("")
			})
			result := getLsblkResult()
			So(result, ShouldBeNil)
		})
	})
}

//GetDeviceTypeMap
func TestGetDeviceTypeMap(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetDeviceTypeMap", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getLsblkResult, func() []string {
				strings := []string{
					"1 2 3 4",
					"disk disk disk disk disk disk disk",
					"part part part part part part part",
					"lvm lvm lvm lvm lvm lvm lvm lvm",
					"5 5 5 5 5 5 5 5 5", "6", "7"}
				return strings
			})
			typeMap := GetDeviceTypeMap()
			So(typeMap, ShouldNotBeNil)
		})
	})
}

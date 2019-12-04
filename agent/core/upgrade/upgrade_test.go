package upgrade

import (
	"errors"
	"fmt"
	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//isDownloaded
func TestIsDownloaded(t *testing.T) {
	//go InitConfig()
	Convey("Test_isDownloaded", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				return nil, errors.New("")
			})
			downloaded := isDownloaded("")
			So(downloaded, ShouldBeFalse)
		})
		Convey("test case2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				dir, _ := os.Getwd()
				path := dir + "/def.go"
				info, _ := os.Lstat(path)
				return info, nil
			})
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, errors.New("")
			})
			downloaded := isDownloaded("")
			So(downloaded, ShouldBeFalse)
		})
		Convey("test case3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				dir, _ := os.Getwd()
				path := dir + "/def.go"
				info, _ := os.Lstat(path)
				return info, nil
			})
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, nil
			})
			downloaded := isDownloaded("")
			So(downloaded, ShouldBeFalse)
		})
	})
}

//Download
func TestDownload(t *testing.T) {
	//go InitConfig()
	Convey("Test_Download", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return true
			})
			download := Download("", "", "")
			So(download, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return false
			})
			mockfn.Replace(utils.CreateDir, func(dir string) error {
				return errors.New("123")
			})
			download := Download("", "", "")
			So(download.Error(), ShouldEqual, "123")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return false
			})
			mockfn.Replace(utils.CreateDir, func(dir string) error {
				return nil
			})
			mockfn.Replace(utils.HTTPGet, func(url string) ([]byte, error) {
				return nil, errors.New("123")
			})
			download := Download("", "", "")
			So(download.Error(), ShouldEqual, "123")
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return false
			})
			mockfn.Replace(utils.CreateDir, func(dir string) error {
				return nil
			})
			mockfn.Replace(utils.HTTPGet, func(url string) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(ioutil.WriteFile, func(filename string, data []byte, perm os.FileMode) error {
				return errors.New("123")
			})
			download := Download("", "", "")
			So(download.Error(), ShouldEqual, "123")
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return false
			})
			mockfn.Replace(utils.CreateDir, func(dir string) error {
				return nil
			})
			mockfn.Replace(utils.HTTPGet, func(url string) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(ioutil.WriteFile, func(filename string, data []byte, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(fmt.Sprintf, func(format string, a ...interface{}) string {
				return "str"
			})
			download := Download("", "", "str1")
			So(download, ShouldNotBeNil)
		})
		Convey("test case 6", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isDownloaded, func(version string) bool {
				return false
			})
			mockfn.Replace(utils.CreateDir, func(dir string) error {
				return nil
			})
			mockfn.Replace(utils.HTTPGet, func(url string) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(ioutil.WriteFile, func(filename string, data []byte, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(fmt.Sprintf, func(format string, a ...interface{}) string {
				return "str"
			})
			mockfn.Replace(filepath.Join, func(elem ...string) string {
				return ":::"
			})
			download := Download("", "", "str")
			So(download, ShouldBeNil)
		})
	})
}

//extractNameFromUrl
func TestExtractNameFromUrl(t *testing.T) {
	//go InitConfig()
	Convey("Test_extractNameFromUrl", t, func() {
		Convey("test case 1", func() {
			url := extractNameFromUrl("a/b", "c")
			So(url, ShouldEqual, "b")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(strings.Split, func(s, sep string) []string {
				return []string{}
			})
			url := extractNameFromUrl("", "c")
			So(url, ShouldEqual, "c")
		})
	})
}

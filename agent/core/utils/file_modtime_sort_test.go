package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"os"
	"sort"
	"testing"
)

//Less
func TestLess(t *testing.T) {
	//go InitConfig()
	Convey("Test_Less", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				dir, _ := os.Getwd()
				path := dir + "/file_modtime_sort.go"
				info, _ := os.Lstat(path)
				return info, nil
			})
			infoStrs := fileInfoStr{}
			strs := []string{""}
			infoStrs = strs
			less := infoStrs.Less(0, 0)
			So(less, ShouldBeTrue)
		})
	})
}

//Len
func TestLenErrInvalidFields(t *testing.T) {
	//go InitConfig()
	Convey("Test_Len", t, func() {
		Convey("test case 1", func() {
			infoStrs := fileInfoStr{}
			len := infoStrs.Len()
			So(len, ShouldEqual, 0)
		})
	})
}

//Swap
func TestSwap(t *testing.T) {
	//go InitConfig()
	Convey("Test_Swap", t, func() {
		Convey("test case 1", func() {
			infoStrs := fileInfoStr{}
			strs := []string{""}
			infoStrs = strs
			infoStrs.Swap(0, 0)
		})
	})
}

//FileListSortTimeAsc
func TestFileListSortTimeAsc(t *testing.T) {

	Convey("Test_FileListSortTimeAsc", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(sort.Sort, func(data sort.Interface) {
				return
			})
			strs := []string{"1", "2"}
			files := FileListSortTimeAsc(strs)
			So(files, ShouldNotBeEmpty)
		})
		Convey("test case 2", func() {
			strs := []string{}
			files := FileListSortTimeAsc(strs)
			So(files, ShouldBeEmpty)
		})
	})
}

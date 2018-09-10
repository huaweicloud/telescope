package model

import (
	. "github.com/smartystreets/goconvey/convey"

	"bytes"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/process"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//GetTop5CpuProcessList
func TestGetTop5CpuProcessList(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetTop5CpuProcessList", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(process.Processes, func() ([]*process.Process, error) {
				processes := []*process.Process{
					{Pid: 1},
					{Pid: 3},
					{Pid: 2},
					{Pid: 7},
					{Pid: 5},
				}
				return processes, nil
			})
			//Percent
			mockfn.Replace((*process.Process).Percent, func(p *process.Process, interval time.Duration) (float64, error) {
				return 0, nil
			})
			mockfn.Replace((*process.Process).CreateTime, func(p *process.Process) (int64, error) {
				return 0, nil
			})
			//Cmdline
			mockfn.Replace((*process.Process).Cmdline, func(p *process.Process) (string, error) {
				var buffer bytes.Buffer
				s := "A"
				for i := 0; i < 4096+1; i++ {
					buffer.WriteString(s)
				}
				s = buffer.String()
				return s, nil
			})
			list := GetTop5CpuProcessList()
			So(list, ShouldNotBeNil)
		})
	})
}

//Len
func TestLen(t *testing.T) {
	//go InitConfig()
	Convey("Test_Len", t, func() {
		Convey("test case 1", func() {
			lists := CPUProcessList{}
			i := lists.Len()
			So(i, ShouldEqual, 0)
		})
	})
}

//Swap
func TestSwap(t *testing.T) {
	//go InitConfig()
	Convey("Test_Swap", t, func() {
		Convey("test case 1", func() {
			lists := CPUProcessList{
				{Pid: 1},
				{Pid: 2},
			}
			lists.Swap(0, 1)
			So(lists[0].Pid, ShouldEqual, 2)
		})
	})
}

//Less
func TestLess(t *testing.T) {
	//go InitConfig()
	Convey("Test_Less", t, func() {
		Convey("test case 1", func() {
			lists := CPUProcessList{
				{Pid: 1},
				{Pid: 2},
			}
			less := lists.Less(0, 1)
			So(less, ShouldBeFalse)
		})
	})
}

//BuildProcessInfoByList
func TestBuildProcessInfoByList(t *testing.T) {
	//go InitConfig()
	Convey("Test_BuildProcessInfoByList", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			lists := ChProcessList{
				{Pid: 1},
				{Pid: 2},
			}
			list := BuildProcessInfoByList(lists)
			So(list, ShouldNotBeNil)
		})
	})
}

//GenerateHashID
func TestGenerateHashIDProcess(t *testing.T) {
	//go InitConfig()
	Convey("Test_GenerateHashID", t, func() {
		Convey("test case 1", func() {
			id := GenerateHashID("1", 2)
			So(id, ShouldNotBeBlank)
		})
	})
}

//GenerateHashIDByPname
func TestGenerateHashIDByPname(t *testing.T) {
	//go InitConfig()
	Convey("Test_GenerateHashIDByPname", t, func() {
		Convey("test case 1", func() {
			pname := GenerateHashIDByPname("1")
			So(pname, ShouldNotBeBlank)
		})
	})
}

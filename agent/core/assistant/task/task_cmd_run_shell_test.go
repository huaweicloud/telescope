package task

import (
	"errors"
	"github.com/huaweicloud/telescope/agent/core/logs"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"os"
	"os/exec"
	"testing"
)

//InstantiateTaskRunShell
func TestInstantiateTaskRunShell(te *testing.T) {

	Convey("Test_InstantiateTaskRunShell", te, func() {
		Convey("test case 1", func() {
			pulled := TaskPulled{}
			shell := InstantiateTaskRunShell(&pulled)
			So(shell, ShouldNotBeNil)
		})
	})
}

//Run TaskRunShell
func TestRunTaskRunShell(te *testing.T) {

	Convey("Test_RunTaskRunShell", te, func() {
		Convey("test case 1", func() {
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD:        nil,
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := trs.Run(ints, strings)
			So(run, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isPathExisted, func(path string) (bool, error) {
				return true, errors.New("123")
			})
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD:        nil,
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := trs.Run(ints, strings)
			So(run, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isPathExisted, func(path string) (bool, error) {
				return true, nil
			})
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD:        nil,
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := trs.Run(ints, strings)
			So(run, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isPathExisted, func(path string) (bool, error) {
				return true, nil
			})
			mockfn.Replace((*exec.Cmd).Start, func(*exec.Cmd) error {
				return nil
			})
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD:        nil,
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := trs.Run(ints, strings)
			So(run, ShouldBeNil)
		})
	})
}

//Pause TaskRunShell
func TestPauseTaskRunShell(te *testing.T) {

	Convey("Test_TaskRunShell", te, func() {
		Convey("test case 1", func() {
			trs := TaskRunShell{
				CMD: nil,
			}
			pause := trs.Pause()
			So(pause, ShouldBeNil)
		})
	})
}

//Cancel TaskRunShell
func TestCancelTaskRunShell(te *testing.T) {

	Convey("Test_CancelTaskRunShell", te, func() {
		Convey("test case 1", func() {
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD:        &exec.Cmd{},
			}
			cancel := trs.Cancel()
			So(cancel, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*os.Process).Kill, func(*os.Process) error {
				return errors.New("")
			})
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD: &exec.Cmd{
					Process: &os.Process{},
				},
			}
			cancel := trs.Cancel()
			So(cancel, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*os.Process).Kill, func(*os.Process) error {
				return nil
			})
			pulled := TaskPulled{
				CmdMeta: CmdMeta{
					Command:  "",
					WorkPath: "",
				},
			}
			trs := TaskRunShell{
				TaskPulled: &pulled,
				CMD: &exec.Cmd{
					Process: &os.Process{},
				},
			}
			trs.CMD.Process = &os.Process{}
			cancel := trs.Cancel()
			So(cancel, ShouldBeNil)
		})
	})
}

//CarryOn TaskRunShell
func TestCarryOnTaskRunShell(te *testing.T) {

	Convey("Test_CarryOnTaskRunShell", te, func() {
		Convey("test case 1", func() {
			trs := TaskRunShell{
				CMD: &exec.Cmd{},
			}
			on := trs.CarryOn()
			So(on, ShouldBeNil)
		})
	})
}

//isPathExisted
func TestIsPathExisted(te *testing.T) {
	Convey("Test_isPathExisted", te, func() {
		Convey("test case 1", func() {
			b, e := isPathExisted("")
			So(b, ShouldBeFalse)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.IsNotExist, func(err error) bool {
				return false
			})
			b, e := isPathExisted("")
			So(b, ShouldBeTrue)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			pwd := logs.GetCurrentDirectory()
			b, e := isPathExisted(pwd)
			So(b, ShouldBeTrue)
			So(e, ShouldBeNil)
		})
	})
}

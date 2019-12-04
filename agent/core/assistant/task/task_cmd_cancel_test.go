package task

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//InstantiateTaskCancel
func TestInstantiateTaskCancel(te *testing.T) {
	Convey("Test_InstantiateTaskCancel", te, func() {
		Convey("test case 1", func() {
			cancel := InstantiateTaskCancel(nil)
			So(cancel, ShouldNotBeNil)
		})
	})
}

//Run(exitCodeChan chan int, outputChan chan *string)
func TestRunTaskCancel(te *testing.T) {
	Convey("Test_Run", te, func() {
		Convey("test case 1", func() {
			cancel := &TaskCancel{
				TaskPulled: &TaskPulled{
					InvocationID: "",
				},
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := cancel.Run(ints, strings)
			So(run, ShouldBeNil)
		})
		Convey("test case 2", func() {
			tasks := make(map[string]*Task)
			tasks["123"] = &Task{
				TaskPulled: &TaskPulled{},
			}
			sprint := fmt.Sprint(tasks)
			cancel := &TaskCancel{
				TaskPulled: &TaskPulled{
					InvocationID: "",
					CmdMeta:      CmdMeta{InvocationID: sprint},
				},
			}
			ints := make(chan int, 1)
			strings := make(chan *string, 1)
			run := cancel.Run(ints, strings)
			So(run, ShouldBeNil)
		})
	})
}

//Pause
func TestPauseTaskCancel(te *testing.T) {
	Convey("Test_PauseTaskCancel", te, func() {
		Convey("test case 1", func() {
			cancel := &TaskCancel{
				TaskPulled: &TaskPulled{
					InvocationID: "",
				},
			}
			pause := cancel.Pause()
			So(pause, ShouldBeNil)
		})
	})
}

//Cancel
func TestCancelTaskCancel(te *testing.T) {

	Convey("Test_CancelTaskCancel", te, func() {
		Convey("test case 1", func() {
			cancel := &TaskCancel{
				TaskPulled: &TaskPulled{
					InvocationID: "",
				},
			}
			e := cancel.Cancel()
			So(e, ShouldBeNil)
		})
	})
}

//CarryOn
func TestCarryOnTaskCancel(te *testing.T) {

	Convey("Test_CarryOnTaskCancel", te, func() {
		Convey("test case 1", func() {
			cancel := &TaskCancel{
				TaskPulled: &TaskPulled{
					InvocationID: "",
				},
			}
			on := cancel.CarryOn()
			So(on, ShouldBeNil)
		})
	})
}

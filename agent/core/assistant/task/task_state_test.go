package task

import (
	"container/list"
	"errors"
	"github.com/smallnest/fsm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
)

//OnExit
func TestOnExit(te *testing.T) {
	Convey("Test_OnExit", te, func() {
		Convey("test case 1", func() {
			processor := TaskEventProcessor{}
			processor.OnExit("", nil)
		})
	})
}

//Action
func TestAction(te *testing.T) {

	Convey("Test_Action", te, func() {
		Convey("test case 1", func() {
			args := []interface{}{""}
			processor := TaskEventProcessor{}
			processor.Action("", "", "", args)
		})
		Convey("test case 2", func() {
			//ExecTaskAction
			defer mockfn.RevertAll()
			mockfn.Replace(ExecTaskAction, func(action string, task *TaskExecEntity) error {
				return nil
			})
			//updateStates
			mockfn.Replace((*TaskExecEntity).updateStates, func(t *TaskExecEntity, toState string) {
				return
			})
			pulled := TaskPulled{}
			entity := TaskExecEntity{
				TaskPulled: &pulled,
				State:      STATE_CANCELED,
			}
			args := []interface{}{""}
			args[0] = &entity
			processor := TaskEventProcessor{}
			processor.Action("", "", "", args)
		})
		Convey("test case 3", func() {
			//ExecTaskAction
			defer mockfn.RevertAll()
			mockfn.Replace(ExecTaskAction, func(action string, task *TaskExecEntity) error {
				return errors.New("")
			})
			//updateStates
			mockfn.Replace((*TaskExecEntity).updateStates, func(t *TaskExecEntity, toState string) {
				return
			})
			pulled := TaskPulled{}
			entity := TaskExecEntity{
				TaskPulled: &pulled,
				State:      STATE_CANCELED,
			}
			args := []interface{}{""}
			args[0] = &entity
			processor := TaskEventProcessor{}
			processor.Action("", "", "", args)
		})
	})
}

//initTaskFSM
func TestInitTaskFSM(te *testing.T) {
	//initTaskFSM()
	Convey("Test_initTaskFSM", te, func() {
		Convey("test case 1", func() {
			//fsm.NewStateMachine
			defer mockfn.RevertAll()
			mockfn.Replace(fsm.NewStateMachine, func(delegate fsm.Delegate, transitions ...fsm.Transition) *fsm.StateMachine {
				return nil
			})
			initTaskFSM()
		})
	})
}

//ExecTaskAction
func TestExecTaskAction(te *testing.T) {
	Convey("Test_ExecTaskAction", te, func() {
		Convey("test case 1", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).Pause, func(*TaskExecEntity) error {
				return errors.New("123")
			})
			ExecTaskAction(ACTION_PAUSE, entity)
		})
		Convey("test case 2", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).CarryOn, func(*TaskExecEntity) error {
				return errors.New("123")
			})
			ExecTaskAction(ACTION_CARRY_ON, entity)
		})
		Convey("test case 3", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).Cancel, func(*TaskExecEntity) error {
				return errors.New("123")
			})
			ExecTaskAction(ACTION_TIMEOUT, entity)
		})
		Convey("test case 4", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).Cancel, func(*TaskExecEntity) error {
				return errors.New("123")
			})
			ExecTaskAction(ACTION_CANCEL, entity)
		})
		Convey("test case 5", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).Cancel, func(*TaskExecEntity) error {
				return errors.New("123")
			})
			//ACTION_RETURN_NONZERO
			ExecTaskAction(ACTION_RETURN_NONZERO, entity)
		})
		Convey("test case 6", func() {
			//fsm.NewStateMachine
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			str := "123"
			strings <- &str
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- 1
			entity := &TaskExecEntity{
				TaskInterface: &TaskCancel{},
				TaskPulled:    &TaskPulled{},
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
				States:        list.New(),
				Events:        list.New(),
			}
			ExecTaskAction("default", entity)
		})
	})
}

//RunAndListen
func TestRunAndListen(te *testing.T) {

	Convey("Test_RunAndListen", te, func() {
		Convey("test case 1", func() {
			pulled := TaskPulled{}
			cancel := TaskCancel{}
			ints := make(chan int, MAX_RETURN_CODE_CHAN_SIZE)
			ints <- EXIT_CODE_ZERO
			strings := make(chan *string, MAX_OUTPUT_CHAN_SIZE)
			temp := "123"
			strings <- &temp
			entity := TaskExecEntity{
				TaskInterface: &cancel,
				TaskPulled:    &pulled,
				EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
				ExitCodeChan:  ints,
				ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
				OutputChan:    strings,
				Output:        nil,
				State:         STATE_CANCELED,
			}
			//updateEvents
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).updateEvents, func(t *TaskExecEntity, event string) {
				te.SkipNow()
			})
			entity.RunAndListen()
		})
	})
}

//getOutputString
func TestGetOutputString(te *testing.T) {
	Convey("Test_getOutputString", te, func() {
		Convey("test case 1", func() {
			s := getOutputString(nil)
			So(s, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			str := "123"
			s := getOutputString(&str)
			So(s, ShouldEqual, "123")
		})
	})
}

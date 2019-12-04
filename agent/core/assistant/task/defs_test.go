package task

import (
	"container/list"
	"github.com/robfig/cron"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
)

//updateEvents
func TestUpdateEvents(te *testing.T) {
	Convey("Test_updateEvents", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{
				Output: nil,
				State:  STATE_CREATED,
				States: list.New(),
				Events: list.New(),
			}
			taskExecEntity.updateEvents("")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*list.List).Remove, func(l *list.List, e *list.Element) interface{} {
				return 1
			})
			taskExecEntity := TaskExecEntity{
				Output: nil,
				State:  STATE_CREATED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			for i := 0; i < 20; i++ {
				l.PushBack(i)
			}
			taskExecEntity.States = l
			taskExecEntity.Events = l
			taskExecEntity.updateEvents("")
		})
	})
}

//updateStates
func TestUpdateStates(te *testing.T) {
	Convey("Test_updateStates", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{
				States: list.New(),
				Events: list.New(),
			}
			taskExecEntity.updateStates("")
		})
		Convey("test case 2", func() {
			taskExecEntity := TaskExecEntity{
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			for i := 0; i < 20; i++ {
				l.PushBack(i)
			}
			taskExecEntity.States = l
			taskExecEntity.updateStates("")
		})
	})
}

//rollbackStates
func TestRollbackStates(te *testing.T) {

	Convey("Test_rollbackStates", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{
				State:  STATE_CREATED,
				States: list.New(),
				Events: list.New(),
			}
			taskExecEntity.rollbackStates("")
		})
		Convey("test case 2", func() {
			taskExecEntity := TaskExecEntity{
				State:  STATE_CREATED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushBack(123)
			l.PushBack(123)
			taskExecEntity.States = l
			taskExecEntity.rollbackStates("")
		})
		Convey("test case 3", func() {
			taskExecEntity := TaskExecEntity{
				State:  STATE_CREATED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushBack("")
			l.PushBack("")
			taskExecEntity.States = l
			taskExecEntity.rollbackStates("")
		})
	})
}

//isFinalState
func TestIsFinalState(te *testing.T) {
	Convey("Test_isFinalState", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{
				State: STATE_CANCELED,
			}
			state := taskExecEntity.isFinalState(true)
			So(state, ShouldBeTrue)
		})
		Convey("test case 2", func() {
			taskExecEntity := TaskExecEntity{
				State: STATE_CANCELED,
			}
			state := taskExecEntity.isFinalState(false)
			So(state, ShouldBeTrue)
		})
		Convey("test case 3", func() {
			taskExecEntity := TaskExecEntity{
				State: STATE_CANCELED,
			}
			taskExecEntity.State = "aa"
			state := taskExecEntity.isFinalState(false)
			So(state, ShouldBeFalse)
		})
	})
}

//pushTaskExec
func TestPushTaskExec(te *testing.T) {
	Convey("Test_pushTaskExec", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{}
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			task.pushTaskExec(&taskExecEntity)
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{}
			l := list.New()
			for i := 0; i < 20; i++ {
				l.PushBack(i)
			}
			task.TaskExecEntityList = l
			task.pushTaskExec(&taskExecEntity)
		})
		Convey("test case 3", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{}
			defer mockfn.RevertAll()
			mockfn.Replace((*list.List).Back, func(*list.List) *list.Element {
				element := &list.Element{
					Value: &taskExecEntity,
				}
				return element
			})
			//TaskInterface
			mockfn.Replace(TaskInterface.Cancel, func(TaskInterface) error {
				return nil
			})
			l := list.New()
			for i := 0; i < 20; i++ {
				l.PushBack(i)
			}
			task.TaskExecEntityList = l
			task.pushTaskExec(&taskExecEntity)
		})
	})
}

//sendCancelEvent
func TestSendCancelEvent(te *testing.T) {

	Convey("Test_sendCancelEvent", te, func() {
		Convey("test case 1", func() {
			taskExecEntity := TaskExecEntity{
				EventChan: make(chan string, MAX_EVENT_CHAN_SIZE),
				States:    list.New(),
				Events:    list.New(),
			}
			l := list.New()
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: l,
			}
			l.PushFront(&taskExecEntity)
			task.sendCancelEvent()
		})
		Convey("test case 2", func() {
			taskExecEntity := TaskExecEntity{
				EventChan: make(chan string, MAX_EVENT_CHAN_SIZE),
				States:    list.New(),
				Events:    list.New(),
			}
			l := list.New()
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: l,
			}

			l.PushFront(taskExecEntity)
			task.TaskExecEntityList = l
			task.sendCancelEvent()
		})
		Convey("test case 3", func() {
			taskExecEntity := TaskExecEntity{
				EventChan: make(chan string, MAX_EVENT_CHAN_SIZE),
				States:    list.New(),
				Events:    list.New(),
			}
			l := list.New()
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: l,
			}

			l.PushFront(&taskExecEntity)
			task.Cron = new(cron.Cron)
			task.TaskExecEntityList = l
			task.sendCancelEvent()
		})
		Convey("test case 4", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			task.Cron = new(cron.Cron)
			task.sendCancelEvent()
		})
	})
}

//isCronTask
func TestIsCronTask(te *testing.T) {

	Convey("Test_isCronTask", te, func() {
		Convey("test case 1", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			cronTask := task.isCronTask()
			So(cronTask, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			task.TaskPulled.Cron = "123"
			cronTask := task.isCronTask()
			So(cronTask, ShouldBeTrue)
		})
	})
}

//isCancelTask
func TestIsCancelTask(te *testing.T) {

	Convey("Test_isCancelTask", te, func() {
		Convey("test case 1", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			cancelTask := task.isCancelTask()
			So(cancelTask, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			task.TaskPulled.Command = TASK_COMMAND_CANCEL
			cancelTask := task.isCancelTask()
			So(cancelTask, ShouldBeTrue)
		})
	})
}

//getTaskState
func TestGetTaskState(t *testing.T) {

	Convey("Test_getTaskState", t, func() {
		Convey("test case 1", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(&taskExecEntity)
			task.TaskExecEntityList = l
			state := task.getTaskState()
			So(state, ShouldEqual, "SUCCEEDED")
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(taskExecEntity)
			task.TaskExecEntityList = l
			state := task.getTaskState()
			So(state, ShouldBeBlank)
		})
		Convey("test case 3", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(&taskExecEntity)
			task.TaskExecEntityList = l
			task.TaskPulled.Cron = "123"
			state := task.getTaskState()
			So(state, ShouldEqual, "RUNNING")
		})
		Convey("test case 4", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			state := task.getTaskState()
			So(state, ShouldEqual, "RUNNING")
		})
	})
}

//getTaskOutput
func TestGetTaskOutput(t *testing.T) {

	Convey("Test_getTaskOutput", t, func() {
		Convey("test case 1", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			output := task.getTaskOutput()
			So(output, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(&taskExecEntity)
			task.TaskExecEntityList = l
			output := task.getTaskOutput()
			So(output, ShouldBeBlank)
		})
		Convey("test case 3", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(taskExecEntity)
			task.TaskExecEntityList = l
			output := task.getTaskOutput()
			So(output, ShouldBeBlank)
		})
		Convey("test case 4", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			temp := "123" // *string cannot be initialized
			taskExecEntity.Output = &temp
			l.PushFront(&taskExecEntity)
			task.TaskExecEntityList = l
			task.TaskPulled.Cron = "123"
			output := task.getTaskOutput()
			So(output, ShouldEqual, "123")
		})
		Convey("test case 5", func() {
			task := Task{
				TaskPulled:         &TaskPulled{},
				TaskExecEntityList: list.New(),
			}

			l := list.New()
			l.PushFront(nil)
			task.TaskExecEntityList = l
			output := task.getTaskOutput()
			So(output, ShouldBeBlank)
		})
	})
}

//getErrNo
func TestGetErrNo(te *testing.T) {

	Convey("Test_getErrNo", te, func() {
		Convey("test case 1", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(taskExecEntity)
			task.TaskExecEntityList = l

			no := task.getErrNo()
			So(no, ShouldNotBeBlank)
		})
		Convey("test case 2", func() {
			task := Task{
				TaskPulled: &TaskPulled{},
			}
			taskExecEntity := TaskExecEntity{
				State:  STATE_SUCCEEDED,
				States: list.New(),
				Events: list.New(),
			}
			l := list.New()
			l.PushFront(&taskExecEntity)
			task.TaskExecEntityList = l

			no := task.getErrNo()
			So(no, ShouldBeBlank)
		})
	})
}

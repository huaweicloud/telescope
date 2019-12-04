package task

import (
	"container/list"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/robfig/cron"
)

// TaskPulled ...
type TaskPulled struct {
	TaskID       string  `json:"task_id"`
	InvocationID string  `json:"invocation_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Command      string  `json:"command"`
	CmdMeta      CmdMeta `json:"cmd_meta"`
	Cron         string  `json:"cron"`
	Timeout      int     `json:"timeout"`
	Type         string  `json:"type"`
}

// CmdMeta ...
type CmdMeta struct {
	WorkPath string `json:"work_path"`
	Command  string `json:"command"`
	// TaskID/InvocationID appears only when Command is CANCEL/STOP/PAUSE
	TaskID       string `json:"task_id"`
	InvocationID string `json:"invocation_id"`
}
type PullTaskRequestBody struct {
	InvokeScope string `json:"invoke_scope"`
	Resources   string `json:"resources"`
}

// TaskExecEntity ...
type TaskExecEntity struct {
	TaskInterface

	TaskPulled   *TaskPulled
	EventChan    chan string
	ExitCodeChan chan int
	ExitChan     chan bool
	OutputChan   chan *string
	Output       *string
	State        string
	States       *list.List
	Events       *list.List
}

func (t *TaskExecEntity) updateEvents(event string) {
	if t.States.Len() == MAX_EVENTS_SIZE {
		t.Events.Remove(t.Events.Back())
		t.Events.PushFront(event)
	} else {
		t.Events.PushFront(event)
	}
}

func (t *TaskExecEntity) updateStates(toState string) {
	if t.States.Len() == MAX_STATES_SIZE {
		t.States.Remove(t.States.Back())
		t.States.PushFront(toState)
	} else {
		t.States.PushFront(toState)
	}
	t.State = toState
}

func (t *TaskExecEntity) rollbackStates(toState string) {
	if t.States.Len() > 0 {
		t.States.Remove(t.States.Front())
		if state, ok := t.States.Front().Value.(string); ok {
			t.State = state
		} else {
			logs.GetAssistantLogger().Errorf("Convert t.States.Front().Value to string failed.")
		}
	}
}

func (t *TaskExecEntity) isFinalState(cronFlag bool) bool {
	if cronFlag && t.State == STATE_CANCELED {
		return true
	}

	if !cronFlag && (t.State == STATE_CANCELED ||
		t.State == STATE_FAILED ||
		t.State == STATE_SUCCEEDED ||
		t.State == STATE_TIMEOUT) {
		return true
	}

	return false
}

// Task ...
type Task struct {
	TaskPulled         *TaskPulled
	Cron               *cron.Cron
	TaskExecEntityList *list.List
}

func (t *Task) pushTaskExec(taskExecEntity *TaskExecEntity) {
	if t.TaskExecEntityList == nil {
		t.TaskExecEntityList = list.New()
	}
	if t.TaskExecEntityList.Len() >= CRON_TASKEXEC_MAINTAIN_COUNT {
		logs.GetAssistantLogger().Debug("CRON_TASKEXEC_MAINTAIN_COUNT over")
		e := t.TaskExecEntityList.Back()
		backTask, ok := e.Value.(*TaskExecEntity)
		if ok {
			// cancel directly, it's state don't need to trans and record
			backTask.Cancel()
		} else {
			logs.GetAssistantLogger().Errorf("Failed to converse e.Value: %v to TaskExecEntity", e.Value)
		}
		t.TaskExecEntityList.Remove(e)
	}

	t.TaskExecEntityList.PushFront(taskExecEntity)
}

func (t *Task) sendCancelEvent() {
	logs.GetAssistantLogger().Debugf("Task get canceled now, TaskExecEntityList len is %d, %v", t.TaskExecEntityList.Len, t.TaskExecEntityList)
	// stop cron
	if t.Cron != nil {
		logs.GetAssistantLogger().Debug("Is cron task been canceling")
		t.Cron.Stop()
	}
	for e := t.TaskExecEntityList.Front(); e != nil; e = e.Next() {
		taskExecEntity, ok := e.Value.(*TaskExecEntity)
		if ok {
			if t.Cron == nil || taskExecEntity.State == STATE_RUNNING {
				taskExecEntity.EventChan <- EVENT_CANCEL
			} else {
				taskExecEntity.State = STATE_CANCELED
			}
		} else {
			logs.GetAssistantLogger().Errorf("Failed to converse e.Value: %v to TaskExecEntity", e.Value)
		}
	}
}

// IsCronTask return bool indicates the task is cron or not
func (t *Task) isCronTask() bool {
	if t.TaskPulled.Cron == "" {
		return false
	} else {
		return true
	}
}

// IsCancelTask return bool indicates the task is CANCEL type or not
func (t *Task) isCancelTask() bool {
	if t.TaskPulled.Command == TASK_COMMAND_CANCEL {
		return true
	} else {
		return false
	}
}

// getTaskState
func (t *Task) getTaskState() string {
	if t.TaskExecEntityList.Front() == nil {
		logs.GetAssistantLogger().Errorf("t.TaskExecEntityList.Front() is nil")
		return STATE_RUNNING
	}

	if taskExecEntity, ok := t.TaskExecEntityList.Front().Value.(*TaskExecEntity); ok {
		// task state calculate from taskExec, cron task only has STATE_RUNNING and STATE_CANCELED
		if t.isCronTask() && taskExecEntity.State != STATE_CANCELED {
			return STATE_RUNNING
		}

		return taskExecEntity.State
	} else {
		logs.GetAssistantLogger().Errorf("Convert to TaskExecEntity failed, return empty state")
		return ""
	}
}

// getTaskOutput
func (t *Task) getTaskOutput() string {
	taskExecEntityListFront := t.TaskExecEntityList.Front()
	if taskExecEntityListFront == nil {
		logs.GetAssistantLogger().Warnf("Task(ID-%s) list front is nil.", t.TaskPulled.TaskID)
		return ""
	}

	taskExecEntityListFrontValue := t.TaskExecEntityList.Front().Value
	if taskExecEntityListFrontValue == nil {
		logs.GetAssistantLogger().Warnf("Task(ID-%s) list front value is nil.", t.TaskPulled.TaskID)
		return ""
	}

	taskExecEntity, ok := taskExecEntityListFrontValue.(*TaskExecEntity)
	if !ok {
		logs.GetAssistantLogger().Errorf("Convert to TaskExecEntity failed, return empty output")
		return ""
	}

	if taskExecEntity.Output == nil {
		logs.GetAssistantLogger().Warnf("taskExecEntity.Output is nil, return empty output")
		return ""
	} else {
		logs.GetAssistantLogger().Debugf("taskExecEntity.Output is: %s", *(taskExecEntity.Output))
		return *(taskExecEntity.Output)
	}
}

// getErrNo
func (t *Task) getErrNo() string {
	if _, ok := t.TaskExecEntityList.Front().Value.(*TaskExecEntity); ok {
		return ""
	} else {
		logs.GetAssistantLogger().Errorf("Convert to TaskExecEntity failed, return %s", ERROR_INTERNAL)
		return ERROR_INTERNAL
	}
}

// TaskInterface ...
type TaskInterface interface {
	// Blocking
	Run(chan int, chan *string) error
	// Unblocking
	Pause() error
	Cancel() error
	CarryOn() error
}

// PullTaskRespBody ...
type PullTaskRespBody struct {
	Tasks []TaskPulled `json:"tasks"`
}

// ReplyTaskRespBody ...
type ReplyTaskRespBody struct {
	SucceededList []string `json:"succeeded_list"`
	FailedList    []string `json:"failed_list"`
}

// ReplyTask ...
type ReplyTask struct {
	TaskID       string `json:"task_id"`
	InvocationID string `json:"invocation_id"`
	Status       string `json:"status"`
	// ErrNum indicates the internal error in our system, not the return code of certain execution
	ErrNum string `json:"err_no"`
	Output string `json:"output"`
}

// ReplyTaskRequestBody ...
type ReplyTaskRequestBody struct {
	InstanceID string      `json:"instance_id"`
	Tasks      []ReplyTask `json:"tasks"`
}

// Output ...
type Output struct {
	ExitCode int    `json:"exit_code"`
	StdOut   string `json:"std_out"`
	StdErr   string `json:"std_err"`
}

const (
	// TIME_LOCATION_CST represents UTC +08:00
	TIME_LOCATION_CST = "Asia/Shanghai"

	// Reply Task Error NO.
	ERROR_INTERNAL = "asssit.001 Internal Error"

	TASK_COMMAND_INSTALL   = "INSTALL"
	TASK_COMMAND_RUN_SHELL = "RUN_SHELL"
	TASK_COMMAND_CANCEL    = "CANCEL"

	MAX_EVENT_CHAN_SIZE       = 3
	MAX_RETURN_CODE_CHAN_SIZE = 1
	MAX_EXIT_CHAN_SIZE        = 1
	MAX_OUTPUT_CHAN_SIZE      = 1
	MAX_STATES_SIZE           = 20
	MAX_EVENTS_SIZE           = 20

	STATE_CREATED   = "CREATED"
	STATE_RUNNING   = "RUNNING"
	STATE_SLEEPING  = "SLEEPING"
	STATE_TIMEOUT   = "TIMEOUT"
	STATE_FAILED    = "FAILED"
	STATE_SUCCEEDED = "SUCCEEDED"
	STATE_CANCELED  = "CANCELED"

	EVENT_RUN            = "EVENT_RUN"
	EVENT_PAUSE          = "EVENT_PAUSE"
	EVENT_CARRY_ON       = "EVENT_CARRY_ON"
	EVENT_CANCEL         = "EVENT_CANCEL"
	EVENT_TIMEOUT        = "EVENT_TIMEOUT"
	EVENT_RETURN_NONZERO = "EVENT_RETURN_NONZERO"
	EVENT_RETURN_ZERO    = "EVENT_RETURN_ZERO"

	ACTION_RUN            = "ACTION_RUN"
	ACTION_PAUSE          = "ACTION_PAUSE"
	ACTION_CARRY_ON       = "ACTION_CARRY_ON"
	ACTION_CANCEL         = "ACTION_CANCEL"
	ACTION_TIMEOUT        = "ACTION_TIMEOUT"
	ACTION_RETURN_NONZERO = "ACTION_RETURN_NONZERO"
	ACTION_RETURN_ZERO    = "ACTION_RETURN_ZERO"

	CRON_TASKEXEC_MAINTAIN_COUNT = 10
)

var (
	// TaskMap ...
	TaskMap = make(map[string]*Task)
	TaskFSM = initTaskFSM()
)

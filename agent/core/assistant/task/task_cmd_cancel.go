package task

import (
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// TaskCancel the struct type for Run,Pause,Cancel,CarryOn
type TaskCancel struct {
	TaskPulled *TaskPulled
}

// InstantiateTaskCancel ...
func InstantiateTaskCancel(taskPulled *TaskPulled) *TaskCancel {
	return &TaskCancel{
		TaskPulled: taskPulled,
	}
}

// Run command execution
func (tc *TaskCancel) Run(exitCodeChan chan int, outputChan chan *string) error {
	logs.GetAssistantLogger().Debugf("TaskExecEntity.TaskInterface is cancel. %v", tc.TaskPulled)
	aimedTask := TaskMap[tc.TaskPulled.CmdMeta.InvocationID]
	// aimedTask not exists
	if aimedTask == nil {
		logs.GetAssistantLogger().Warnf("AimedTask to cancel does not exist. Invocation id: %s", tc.TaskPulled.CmdMeta.InvocationID)
	} else {
		// send EVENT_CANCEL to aimed Task
		aimedTask.sendCancelEvent()
	}

	outputChan <- nil
	exitCodeChan <- 0
	return nil
}

// Pause command pause, it may continue to be executed, like CarryOn()
func (tc *TaskCancel) Pause() error {
	return nil
}

// Cancel command Cancel, the process can be killed
func (tc *TaskCancel) Cancel() error {
	return nil
}

// CarryOn command CarryOn, the paused task can continue
func (tc *TaskCancel) CarryOn() error {
	return nil
}

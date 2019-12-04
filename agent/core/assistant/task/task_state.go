package task

import (
	"errors"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/smallnest/gofsm"
	"time"
)

// TaskEventProcessor ...
type TaskEventProcessor struct{}

// OnExit ...
func (p *TaskEventProcessor) OnExit(fromState string, args []interface{}) {
}

// Action ...
func (p *TaskEventProcessor) Action(action string, fromState string, toState string, args []interface{}) {
	logs.GetAssistantLogger().Infof("[Action] action is %s, fromState is %s, toState is %s", action, fromState, toState)
	var (
		task *TaskExecEntity
		ok   bool
	)

	if len(args) >= 1 {
		task, ok = args[0].(*TaskExecEntity)
		if !(ok && task.State != "") {
			logs.GetAssistantLogger().Infof("[Action] convert parameter to task failed.")
			return
		} else {
			logs.GetAssistantLogger().Infof("[Action] convert parameter to task successfully.")
		}
	}

	err := ExecTaskAction(action, task)
	if err == nil {
		task.updateStates(toState)
	}
}

// OnEnter ...
func (p *TaskEventProcessor) OnEnter(toState string, args []interface{}) {
}

func initTaskFSM() *fsm.StateMachine {
	delegate := &fsm.DefaultDelegate{P: &TaskEventProcessor{}}

	transitions := []fsm.Transition{
		{From: STATE_CREATED, Event: EVENT_RUN, To: STATE_RUNNING, Action: ACTION_RUN},
		{From: STATE_RUNNING, Event: EVENT_PAUSE, To: STATE_SLEEPING, Action: ACTION_PAUSE},
		{From: STATE_SLEEPING, Event: EVENT_CARRY_ON, To: STATE_RUNNING, Action: ACTION_CARRY_ON},
		{From: STATE_RUNNING, Event: EVENT_TIMEOUT, To: STATE_TIMEOUT, Action: ACTION_TIMEOUT},
		{From: STATE_RUNNING, Event: EVENT_RETURN_NONZERO, To: STATE_FAILED, Action: ACTION_RETURN_NONZERO},
		{From: STATE_RUNNING, Event: EVENT_RETURN_ZERO, To: STATE_SUCCEEDED, Action: ACTION_RETURN_ZERO},
		{From: STATE_RUNNING, Event: EVENT_CANCEL, To: STATE_CANCELED, Action: ACTION_CANCEL},
		{From: STATE_CREATED, Event: EVENT_CANCEL, To: STATE_CANCELED, Action: ACTION_CANCEL},
	}

	return fsm.NewStateMachine(delegate, transitions...)
}

// ExecTaskAction ...
func ExecTaskAction(action string, task *TaskExecEntity) error {
	logs.GetAssistantLogger().Infof("[ExecTaskAction] action is: %s, task is: %v", action, *task)
	var err error = nil

	switch action {
	case ACTION_RUN:
		err = task.Run(task.ExitCodeChan, task.OutputChan)
	case ACTION_PAUSE:
		err = task.Pause()
	case ACTION_CARRY_ON:
		err = task.CarryOn()
	case ACTION_TIMEOUT:
		err = task.Cancel()
	case ACTION_CANCEL:
		err = task.Cancel()
	case ACTION_RETURN_NONZERO:
		fallthrough
	case ACTION_RETURN_ZERO:
		break
	default:
		err = errors.New("Unknown action")
		logs.GetAssistantLogger().Errorf("Unknown action")
	}

	if err != nil {
		logs.GetAssistantLogger().Errorf("[InvocationID-%s] Task executes %s failed and error is: %s", task.TaskPulled.InvocationID, action, err.Error())
	} else {
		logs.GetAssistantLogger().Debugf("[InvocationID-%s] Task executes %s successfully", task.TaskPulled.InvocationID, action)
	}
	return err
}

// RunAndListen ...
func (t *TaskExecEntity) RunAndListen() {
	logs.GetAssistantLogger().Infof("[RunAndListen InvocationID(%s)] Entering(task id: %s)", t.TaskPulled.InvocationID, t.TaskPulled.TaskID)
	t.EventChan <- EVENT_RUN
	var timeout <-chan time.Time
	timeout = time.After(time.Second * time.Duration(t.TaskPulled.Timeout))

L:
	for {
		select {
		// run/carry on/pause/cancel
		case event := <-t.EventChan:
			logs.GetAssistantLogger().Debugf("[RunAndListen InvocationID(%s)] Event(%s) trigger", t.TaskPulled.InvocationID, event)
			logs.GetAssistantLogger().Debugf("[RunAndListen InvocationID(%s)] State is: %s, event is: %s, TaskExecEntity is: %v", t.TaskPulled.InvocationID, t.State, event, t)

			t.updateEvents(event)
			TaskFSM.Trigger(t.State, event, t)

			if t.State == STATE_TIMEOUT || t.State == STATE_FAILED || t.State == STATE_SUCCEEDED || t.State == STATE_CANCELED {
				t.ExitChan <- true
			}

		// timeout
		case <-timeout:
			t.EventChan <- EVENT_TIMEOUT
		case exitCode := <-t.ExitCodeChan:
			logs.GetAssistantLogger().Debugf("[RunAndListen InvocationID(%s)] Exit code(%d) trigger", t.TaskPulled.InvocationID, exitCode)
			output := <-t.OutputChan
			t.Output = output
			logs.GetAssistantLogger().Debugf("[RunAndListen InvocationID(%s)] Recieve *output from channel, and output is: %s", t.TaskPulled.InvocationID, getOutputString(output))
			switch exitCode {
			case 0:
				t.EventChan <- EVENT_RETURN_ZERO
			default:
				t.EventChan <- EVENT_RETURN_NONZERO
			}
		// exit the go routine
		case <-t.ExitChan:
			logs.GetAssistantLogger().Infof("[RunAndListen InvocationID(%s)] Task(%s) received exit signal and break out", t.TaskPulled.InvocationID, t.TaskPulled.TaskID)
			break L
		}
	}

	logs.GetAssistantLogger().Infof("[RunAndListen InvocationID(%s)] Leaving(task id: %s)", t.TaskPulled.InvocationID, t.TaskPulled.TaskID)
}

func getOutputString(output *string) string {
	if output != nil {
		return *output
	}
	return ""
}

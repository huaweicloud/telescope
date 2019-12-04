package task

import (
	"bytes"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"os"
	"os/exec"
	"syscall"
)

const (
	// MAX_STDOUT_IN_BYTES ...
	MAX_STDOUT_IN_BYTES                 = 8192
	MAX_STDERR_IN_BYTES                 = 8192
	EXIT_CODE_ZERO                      = 0
	EXIT_CODE_NON_EXISTENT_WORKING_PATH = -1
	NON_EXISTENT_WORKING_PATH_INFO      = "Working directory does not exist."
)

// TaskRunShell
// the struct type for Run,Pause,Cancel,CarryOn
type TaskRunShell struct {
	TaskPulled *TaskPulled
	CMD        *exec.Cmd
}

// InstantiateTaskRunShell ...
func InstantiateTaskRunShell(taskPulled *TaskPulled) *TaskRunShell {
	logs.GetAssistantLogger().Debugf("[InvocationID-%s] Initialize TaskRunShell", taskPulled.InvocationID)
	return &TaskRunShell{
		TaskPulled: taskPulled,
		CMD:        &exec.Cmd{},
	}
}

// Run ...
//command execution
func (trs *TaskRunShell) Run(exitCodeChan chan int, outputChan chan *string) error {
	cmdMeta := trs.TaskPulled.CmdMeta
	cmd := cmdMeta.Command
	workPath := cmdMeta.WorkPath

	// check working directory failed
	isPathExisted, err := isPathExisted(workPath)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[InvocationID-%s] Check path(%s) failed and error is: %s",
			trs.TaskPulled.InvocationID, workPath, err.Error())
		return err
	}

	// working directory does not exist
	if !isPathExisted {
		logs.GetAssistantLogger().Errorf("[InvocationID-%s] Path(%s) does not exist", trs.TaskPulled.InvocationID, workPath)
		output := Output{
			ExitCode: EXIT_CODE_NON_EXISTENT_WORKING_PATH,
			StdOut:   "",
			StdErr:   NON_EXISTENT_WORKING_PATH_INFO,
		}
		outputBytes, _ := json.Marshal(output)
		outputStr := string(outputBytes)
		outputChan <- &outputStr
		exitCodeChan <- EXIT_CODE_NON_EXISTENT_WORKING_PATH
		return nil
	}

	// execute shell
	trs.CMD = exec.Command("/bin/bash", "-c", cmd)
	trs.CMD.Dir = workPath
	var (
		stdoutBuf, stderrBuf     bytes.Buffer
		output                   Output
		stdoutBytes, stderrBytes []byte
	)
	trs.CMD.Stdout = &stdoutBuf
	trs.CMD.Stderr = &stderrBuf
	err = trs.CMD.Start()
	// Run failed may due to OOM or sth. else
	if err != nil {
		logs.GetAssistantLogger().Errorf("[InvocationID-%s] Exec run of \"%s\" failed and error is: ",
			trs.TaskPulled.InvocationID, cmd, err.Error())
		return err
	}

	go func() {
		if err := trs.CMD.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0
				// This works on both Unix and Windows. Although package
				// syscall is generally platform dependent, WaitStatus is
				// defined for both Unix and Windows and in both cases has
				// an ExitStatus() method with the same signature.
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					logs.GetAssistantLogger().Errorf("[InvocationID-%s] Run \"%s\" finished and exit code is: %d",
						trs.TaskPulled.InvocationID, cmd, status.ExitStatus())
					output.ExitCode = status.ExitStatus()
				}
			} else {
				logs.GetAssistantLogger().Infof("[InvocationID-%s] Run \"%s\" finished and exit code is: %d",
					trs.TaskPulled.InvocationID, cmd, EXIT_CODE_ZERO)
				output.ExitCode = EXIT_CODE_ZERO
			}
		}
		logs.GetAssistantLogger().Debugf("[InvocationID-%s] Run \"%s\" finished", trs.TaskPulled.InvocationID, cmd)

		// cut
		stdoutBytes = stdoutBuf.Bytes()
		if len(stdoutBytes) > MAX_STDOUT_IN_BYTES {
			stdoutBytes = stdoutBytes[0:MAX_STDOUT_IN_BYTES]
		}
		stderrBytes = stderrBuf.Bytes()
		if len(stderrBytes) > MAX_STDERR_IN_BYTES {
			stderrBytes = stderrBytes[0:MAX_STDERR_IN_BYTES]
		}
		output.StdOut, output.StdErr = string(stdoutBytes), string(stderrBytes)

		// struct to string
		outputBytes, err := json.Marshal(output)
		if err != nil {
			logs.GetAssistantLogger().Errorf("[InvocationID-%s] Marshal strut to string failed and error is: ",
				trs.TaskPulled.InvocationID, err.Error())
		}
		outputStr := string(outputBytes)
		exitCodeChan <- output.ExitCode
		outputChan <- &outputStr
		logs.GetAssistantLogger().Debugf("[InvocationID-%s] Run finished, output is: %s", trs.TaskPulled.InvocationID, outputStr)
	}()

	return nil
}

// Pause ...
//command pause, it may continue to be executed, like CarryOn()
func (trs *TaskRunShell) Pause() error {
	return nil
}

// Cancel ...
//command Cancel, the process can be killed
func (trs *TaskRunShell) Cancel() error {
	if trs.CMD.Process == nil {
		return nil
	}
	err := trs.CMD.Process.Kill()
	if err != nil {
		logs.GetAssistantLogger().Errorf("[InvocationID-%s] Kill process failed and error is: %s", trs.TaskPulled.InvocationID, err.Error())
		return err
	}
	logs.GetAssistantLogger().Infof("[InvocationID-%s] Kill process(%d) successfully", trs.TaskPulled.InvocationID, trs.CMD.Process.Pid)

	return nil
}

// CarryOn ...
//command CarryOn, the paused task can continue
func (trs *TaskRunShell) CarryOn() error {
	return nil
}

//judge whether a path exists
func isPathExisted(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

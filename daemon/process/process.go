package process

import (
	"github.com/huaweicloud/telescope/agent/core/manager"
	"os"
	"os/exec"
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
)

var (
	AgentProcDoneChan    = make(chan *os.Process, 10)
	sendSignalRetryCount = 3
)

// GetAgentVersion get current agent version, command: agent -version
func GetAgentVersion(binPath string) (string, error) {
	cmd := exec.Command(binPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// SendAgentSignal signal agent, command: agent stop/upgrade
func SendAgentSignal(binPath string, signal os.Signal) error {
	var cmd *exec.Cmd
	switch signal {
	case upgrade.SIG_UPGRADE:
		cmd = exec.Command(binPath, "upgrade")
	case manager.SIG_STOP:
		cmd = exec.Command(binPath, "stop")
	default:
		cmd = exec.Command(binPath, "-version")
	}

	_, err := cmd.Output()
	return err
}

// StartProcess start new child process
func StartProcess(binPath string) (*os.Process, error) {
	cmd := exec.Command(binPath)
	e := os.Environ()
	cmd.Env = e
	cmd.Args = os.Args
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logs.GetLogger().Errorf("Start process failed, err:%s", err.Error())
		return nil, err
	}

	proc := cmd.Process
	//wait process, otherwise the process becomes zombie after kill
	go func() {
		err := cmd.Wait()
		if err != nil {
			logs.GetLogger().Errorf("Agent process wait out of error: %s.", err.Error())
		}
		AgentProcDoneChan <- cmd.Process
		logs.GetLogger().Infof("AgentProcDoneChan length: %d.", len(AgentProcDoneChan))
	}()
	return proc, nil
}

// KillProcess kill child process
func KillProcess(proc *os.Process) error {
	if proc == nil {
		return nil
	}

	osName := runtime.GOOS
	if osName == "windows" {
		return proc.Kill()
	} else {
		err := proc.Signal(upgrade.SIG_UPGRADE)
		if err != nil {
			logs.GetLogger().Errorf("Stop(SIG_UPGRADE) linux agent process failed, err:%s", err.Error())
			return err
		}
		// wait old process finished
		proc.Wait()
		return nil
	}
}

// StopProcess send and kill child process
func SigAndKillProcess(binPath string, signal os.Signal, proc *os.Process) error {
	if proc == nil {
		return nil
	}
	err := proc.Kill()

	if err != nil {
		return err
	}
	tryCount := 0
	for tryCount < sendSignalRetryCount {
		err = SendAgentSignal(binPath, signal)
		if err == nil {
			break
		}
		tryCount = tryCount + 1
	}

	return nil
}

// StopProcess ...
func StopProcess(binPath string, proc *os.Process) error {
	return SigAndKillProcess(binPath, manager.SIG_STOP, proc)
}

// UpgradeProcess ...
func UpgradeProcess(binPath string, proc *os.Process) error {
	return SigAndKillProcess(binPath, upgrade.SIG_UPGRADE, proc)
}

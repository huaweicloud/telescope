package circuitbreaker

import (
	"os"
	"path/filepath"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/manager"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/nightlyone/lockfile"
	"github.com/shirou/gopsutil/process"
)

// SingleInstanceCheck check whether the agent process is running before start
// error indicates there is another agent process is running or sth. wrong to check
func SingleInstanceCheck() (err error) {
	lock, err := lockfile.New(filepath.Join(utils.GetWorkingPath(), AgentPIDFileName))
	if err != nil {
		logs.GetCesLogger().Errorf("Cannot init pid file and error is:%v", err)
		return
	}
	logs.GetCesLogger().Debugf("New lockfile(%s) successfully or existed", string(lock))

	err = lock.TryLock()
	if err != nil {
		logs.GetCesLogger().Warnf("Cannot lock, reason is: %v", err)
		p, errInner := lock.GetOwner()
		if errInner == nil {
			logs.GetCesLogger().Debugf("GetOwner successfully and pid in file is: %d", p.Pid)
			if !isCurrentAgentPID(p.Pid) {
				return nil
			}
		}
		return
	}

	//defer lock.Unlock()
	return
}

// isCurrentAgentPID check the pid in pid file is the current agent process or not
// any error in it just return true which may be overkill
func isCurrentAgentPID(pid int) bool {
	pidInt32 := int32(pid)

	exist, err := process.PidExists(pidInt32)
	if err != nil {
		logs.GetCesLogger().Errorf("Exec PidExists(%d) failed and error is:%v", pidInt32, err)
		return true
	}
	if !exist {
		logs.GetCesLogger().Infof("PID(%d) in %s does not exist in system, it may be the legacy file from last process", pidInt32, AgentPIDFileName)
		return false
	}

	// check agent process name
	p, err := process.NewProcess(pidInt32)
	if err != nil {
		logs.GetCesLogger().Errorf("Exec process.NewProcess(%d) failed and error is:%v", pidInt32, err)
		return true
	}
	name, err := p.Name()
	if err != nil {
		logs.GetCesLogger().Errorf("Exec p.Name() failed and error is:%v", err)
		return true
	}
	logs.GetCesLogger().Debugf("Agent process name is:%s", name)
	if name == utils.AgentNameLinux || name == utils.AgentNameWin {
		logs.GetCesLogger().Infof("PID(%d) in %s does exist and the name(%s) does MATCH the agent name, which indicates the agent process is running", pidInt32, AgentPIDFileName, name)
		killAgentProcess(p)
		return true
	}
	logs.GetCesLogger().Infof("PID(%d) in %s does exist and the name(%s) does NOT MATCH the agent name, which indicates the agent process is running", pidInt32, AgentPIDFileName, name)

	// check agent's parent process name
	pp, err := p.Parent()
	if err != nil {
		logs.GetCesLogger().Errorf("Exec p.Parent() failed and error is:%v", err)
		return true
	}
	logs.GetCesLogger().Debugf("Parent process ID is: %d", pp.Pid)
	name, err = pp.Name()
	if err != nil {
		logs.GetCesLogger().Errorf("Exec p.Name()[parent] failed and error is:%v", err)
		return true
	}
	logs.GetCesLogger().Debugf("Agent's parent process name is:%s", name)
	if name == utils.DaemonNameLinux || name == utils.DaemonNameWin {
		logs.GetCesLogger().Infof("Agent's parent(daemon - %d) process is running", pp.Pid)
		return true
	}

	return false
}

func killAgentProcess(p *process.Process) {
	agentPid := int(p.Pid)
	sameRunningProc, err := os.FindProcess(agentPid)
	if err != nil {
		logs.GetCesLogger().Errorf("Cannot get process by pid when try to kill agent. %s", err.Error())
		return
	}

	if sameRunningProc == nil {
		logs.GetCesLogger().Warnf("Get process nil by pid: %d, cannot kill it", agentPid)
		return
	}

	logs.GetCesLogger().Errorf("Other agent is running, is a panic scene, kill it. Pid: %d", sameRunningProc.Pid)
	err = sameRunningProc.Signal(manager.SigUserStop)

	if err == nil {
		logs.GetCesLogger().Infof("Send SigUserStop signal success, pid:%d", agentPid)
		return
	}

	logs.GetCesLogger().Errorf("Cannot kill running agent. %s, %v", err.Error(), manager.SigUserStop)
	// check agent's parent process name
	pp, err := p.Parent()
	if err != nil {
		logs.GetCesLogger().Errorf("Exec p.Parent() failed and error is:%v", err)
		return
	}
	logs.GetCesLogger().Debugf("Parent process ID is: %d", pp.Pid)
	name, err := pp.Name()
	if err != nil {
		logs.GetCesLogger().Errorf("Exec p.Name()[parent] failed and error is:%v", err)
		return
	}
	logs.GetCesLogger().Debugf("Agent's parent process name is:%s", name)
	if name == utils.DaemonNameLinux || name == utils.DaemonNameWin {
		logs.GetCesLogger().Infof("Agent's parent(daemon - %d) process is running", pp.Pid)
		err = sameRunningProc.Kill()
		if err != nil {
			logs.GetCesLogger().Warnf("Try to kill agent process error:%v", err)
		}
		return
	}

	return
}

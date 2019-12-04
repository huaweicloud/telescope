package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	agent "github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/huaweicloud/telescope/daemon/process"
	"github.com/huaweicloud/telescope/daemon/upgrade"
	"github.com/kardianos/service"
)

var (
	agentHome     string
	agentTmpHome  string
	agentVersion  string
	agentName     string
	daemonName    string
	agentProcess  *os.Process
	upgradeSignal = make(chan *agent.Info, 1)
	isUpgrading   = false
	isStopping    = false
	upgradeMutex  = &sync.Mutex{}
)

type program struct{}

func (p *program) Start(s service.Service) error {
	if !isAgentFileExisted() {
		return errors.New("error")
	}
	go run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	stopAgent()
	return nil
}

func init() {
	agentHome = logs.GetCurrentDirectory()
	agentTmpHome = agentHome + "/.tmp"

	err := utils.CreateDir(agentTmpHome)
	if err != nil {
		logs.GetLogger().Errorf("Create agent tmp dir failed, err:%v", err)
		logs.GetLogger().Flush()
		os.Exit(-1)
	}

	osName := runtime.GOOS
	switch osName {
	case "windows":
		agentName = utils.AgentNameWin
		daemonName = utils.DaemonNameWin
	case "linux":
		agentName = utils.AgentNameLinux
		daemonName = utils.DaemonNameLinux
	default:
		logs.GetLogger().Errorf("Unsupported os type: %s.", osName)
		logs.GetLogger().Flush()
		os.Exit(-1)
	}

	agentVersion, _ = process.GetAgentVersion(agentHome + "/" + agentName)
	logs.GetLogger().Infof("Current agent version: %s", agentVersion)
}

// wait signal to do upgrade
func upgradeAgent() {
	for {
		info := <-upgradeSignal
		if agentVersion == info.Version {
			logs.GetLogger().Debugf("Agent version is already:%s", agentVersion)
			continue
		}
		upgradeMutex.Lock()
		isUpgrading = true
		proc, err := upgrade.DoUpgrade(agentHome, agentTmpHome, agentName, daemonName, info, agentProcess)
		if err != nil {
			logs.GetLogger().Errorf("Upgrade failed, error is not nil:%s", err.Error())
		}
		if proc == nil {
			logs.GetLogger().Errorf("Upgrade failed, exit process, err:%s", err.Error())
			logs.GetLogger().Flush()
			upgradeMutex.Unlock()
			os.Exit(-1)
		}
		agentProcess = proc

		agentVersion, _ = process.GetAgentVersion(agentHome + "/" + agentName)
		logs.GetLogger().Infof("New agent version: %s", agentVersion)
		isUpgrading = false
		upgradeMutex.Unlock()
	}
}

// start agent
func startAgent() {
	proc, err := process.StartProcess(agentHome + "/" + agentName)
	if err != nil {
		logs.GetLogger().Errorf("Start process failed, err:%s", err)
		logs.GetLogger().Flush()
		os.Exit(-1)
	}
	logs.GetLogger().Infof("Start process successfully and PID is: %d", proc.Pid)
	agentProcess = proc
}

// stop agent
func stopAgent() {
	if agentProcess != nil {
		isStopping = true
		err := process.StopProcess(agentHome+"/"+agentName, agentProcess)
		logs.GetLogger().Info("Agent has been shutdown by telescope.")
		if err != nil {
			logs.GetLogger().Info("Shutdown error." + err.Error())
		}
	}
}

// agent monitor
func startMonitor() {
	var (
		count                   uint = 0
		maxContinuousRetryCount uint = 3
	)
	for {
		logs.GetLogger().Infof("Start to check whether agent is alive and current go routine number is: %d", runtime.NumGoroutine())
		//wait old process finished
		select {
		case oldProc := <-process.AgentProcDoneChan:
			upgradeMutex.Lock()
			count = count + 1
			logs.GetLogger().Error("Agent process had been exited, try to check if need to start again.")
			if isUpgrading || isStopping {
				upgradeMutex.Unlock()
				continue
			}
			if agentProcess != nil && oldProc != agentProcess {
				logs.GetLogger().Warn("Agent process has been changed by other, continue to check.")
				upgradeMutex.Unlock()
				continue
			}

			if count > maxContinuousRetryCount {
				logs.GetLogger().Warn("Agent process has been started %d times continuously. Sleep for one minute.",
					maxContinuousRetryCount)
				go func(oldProc *os.Process) {
					time.Sleep(time.Minute)
					logs.GetLogger().Info("Send old process to process done channel.Try to restart agent process")
					process.AgentProcDoneChan <- oldProc
				}(oldProc)
				upgradeMutex.Unlock()
				continue
			}
			logs.GetLogger().Error("Agent process does not existed, start agent again.")
			startAgent()
			upgradeMutex.Unlock()
		case <-time.After(time.Minute):
			count = 0
			logs.GetLogger().Debug("Wait for agent process timeout, next loop.")
		}
	}
}

func isAgentFileExisted() bool {
	_, err := os.Stat(agentHome + "/" + agentName)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		return true
	}
	fmt.Println(agentName + " does not exist.")
	return false
}

func main() {
	switch runtime.GOOS {
	case "windows":
		logs.GetLogger().Infof("OS type is windows")
		winRun()
	default:
		logs.GetLogger().Infof("OS type is NOT windows")
		run()
	}
}

func run() {
	start := make(chan bool)
	startAgent()
	go startMonitor()
	go upgrade.ScanAgentTmpDir(agentTmpHome, agentName, upgradeSignal)
	go upgradeAgent()
	<-start
}

func winRun() {
	svcConfig := &service.Config{
		Name:        "telescoped",
		DisplayName: "telescoped",
		Description: "Telescoped",
	}

	prg := &program{}
	telescopeService, err := service.New(prg, svcConfig)

	if err != nil {
		logs.GetLogger().Errorf("Get service error: %s", err.Error())
		return
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			serviceInstall(telescopeService)
			return
		}

		if os.Args[1] == "uninstall" {
			serviceUninstall(telescopeService)
			return
		}
	}
	err = telescopeService.Run()
	if err != nil {
		logs.GetLogger().Errorf("Run service error: %s", err.Error())
	}
}

func serviceInstall(telescopeService service.Service) {
	if !isAgentFileExisted() {
		fmt.Println("Failed to install service.")
		return
	}

	telescopeService.Install()
	telescopeService.Start()

	fmt.Println("Install service success.")
}

func serviceUninstall(telescopeService service.Service) {
	telescopeService.Stop()
	telescopeService.Uninstall()
	fmt.Println("Uninstall service success.")
}

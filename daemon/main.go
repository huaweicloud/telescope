package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/logs"
	agent "github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/huaweicloud/telescope/daemon/process"
	"github.com/huaweicloud/telescope/daemon/upgrade"
	"github.com/kardianos/service"
)

var AgentHome string
var AgentTmpHome string
var AgentVersion string
var AgentName string
var DaemonName string
var AgentProcess *os.Process

var upgradeSignal = make(chan *agent.Info, 1)
var isUpgrading = false
var isStopping = false

func init() {
	AgentHome = logs.GetCurrentDirectory()
	AgentTmpHome = AgentHome + "/.tmp"

	err := utils.CreateDir(AgentTmpHome)
	if err != nil {
		logs.GetLogger().Errorf("Create agent tmp dir failed, err:%s", err)
		os.Exit(-1)
	}

	osName := runtime.GOOS
	if osName == "windows" {
		AgentName = utils.AgentNameWin
		DaemonName = utils.DaemonNameWin
	} else if osName == "linux" {
		AgentName = utils.AgentNameLinux
		DaemonName = utils.DaemonNameLinux
	} else {
		logs.GetLogger().Errorf("Unsupport os type: %s.", osName)
		os.Exit(-1)
	}

	AgentVersion, _ = process.GetAgentVersion(AgentHome + "/" + AgentName)
	logs.GetLogger().Infof("Current agent version: %s", AgentVersion)
}

// wait signal to do upgrade
func upgradeAgent() {
	for {
		info := <-upgradeSignal
		if AgentVersion == info.Version {
			logs.GetLogger().Debugf("Agent version is already:%s", AgentVersion)
			continue
		}
		isUpgrading = true
		proc, err := upgrade.DoUpgrade(AgentHome, AgentTmpHome, AgentName, DaemonName, info, AgentProcess)
		if err != nil {
			logs.GetLogger().Errorf("Upgrade failed, error is not nil:%s", err.Error())
		}
		if proc == nil {
			logs.GetLogger().Errorf("Upgrade failed, exit process, err:%s", err.Error())
			os.Exit(-1)
		}
		AgentProcess = proc

		AgentVersion, _ = process.GetAgentVersion(AgentHome + "/" + AgentName)
		logs.GetLogger().Infof("New agent version: %s", AgentVersion)
		isUpgrading = false
	}
}

// start agent
func startAgent() {
	proc, err := process.StartProcess(AgentHome + "/" + AgentName)
	if err != nil {
		logs.GetLogger().Errorf("Start process failed, err:%s", err)
		os.Exit(-1)
	}
	AgentProcess = proc
}

// stop agent
func stopAgent() {
	if AgentProcess != nil {
		isStopping = true
		err := process.StopProcess(AgentHome+"/"+AgentName, AgentProcess)
		logs.GetLogger().Info("Agent has been shutdown by telescope.")
		if err != nil {
			logs.GetLogger().Info("Shutdown error." + err.Error())
		}
	}
}

// agent monitor
func startMonitor() {
	for {
		//wait old process finished
		stat, err := AgentProcess.Wait()
		if (err != nil || stat.Exited()) && !isUpgrading && !isStopping {
			logs.GetLogger().Error("Agent process is not existed, start agent again.")
			startAgent()
		}
	}
}

func isAgentFileExisted() bool {
	_, err := os.Stat(AgentHome + "/" + AgentName)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		return true
	}
	fmt.Println(AgentName + " does not exist.")
	return false
}

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

func main() {

	switch runtime.GOOS {
	case "windows":
		winRun()
	default:
		run()
	}
}

func run() {

	start := make(chan bool)
	startAgent()
	go startMonitor()
	go upgrade.ScanAgentTmpDir(AgentTmpHome, AgentName, upgradeSignal)
	go upgradeAgent()
	<-start
}

func winRun() {

	svcConfig := &service.Config{
		Name:        "telescoped",
		DisplayName: "telescoped",
		Description: "telescoped",
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

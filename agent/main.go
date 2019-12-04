package main

import (
	"fmt"
	cb "github.com/huaweicloud/telescope/agent/core/circuitbreaker"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/manager"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"os"
	"time"
)

func usage() {
	fmt.Println("Usage: agent [OPTION]")
	fmt.Println("  -v, --version      print agent version")
	fmt.Println("  -h, --help         display this help and exit")
	fmt.Println("  stop               send stop signal")
	fmt.Println("  upgrade            send upgrade signal")
}

func main() {
	args := os.Args
	if len(args) == 2 {
		switch args[1] {
		case "-version", "-v", "--version":
			fmt.Println(utils.AgentVersion)
		case "stop":
			logs.GetCesLogger().Info("Agent will send stop signal")
			manager.HandleSignalDirect(manager.SIG_STOP)
		case "upgrade":
			logs.GetCesLogger().Info("This will send upgrade signal")
			manager.HandleSignalDirect(upgrade.SIG_UPGRADE)
		case "?", "-h", "--help":
			fallthrough
		default:
			usage()
		}
		return
	}

	err := cb.SingleInstanceCheck()
	if err != nil {
		logs.GetCesLogger().Warn("Agent will sleep one minute after single instance check error.Pid: %d", os.Getpid())
		logs.GetCesLogger().Flush()
		time.Sleep(time.Minute)
		panic(err)
	}

	if cb.IsSleepNeeded() {
		cb.DelFile()
		cb.Sleep()
	}

	start := make(chan bool)
	serviceManager := manager.NewServicemanager()
	serviceManager.Init()
	serviceManager.RegisterService()
	serviceManager.InitService()
	serviceManager.HeartBeat()
	serviceManager.StartService()
	go cb.Start()
	fmt.Println("Agent starts successfully.")
	logs.GetCesLogger().Info("This is agent main")

	<-start
}

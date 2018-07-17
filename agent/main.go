package main

import (
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	"fmt"
	"os"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/manager"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

func main() {
	arsWithProg := os.Args
	if len(arsWithProg) == 2 {
		switch arsWithProg[1] {
		case "-version":
			fmt.Print(utils.AGENT_VERSION)
			return
		case "stop":
			logs.GetLogger().Info("Agent will send stop signal")
			manager.HandleSignalDirect(manager.SIG_STOP)
			return
		case "upgrade":
			logs.GetLogger().Info("This will send upgrade signal")
			manager.HandleSignalDirect(upgrade.SIG_UPGRADE)
			return
		}

	}

	start := make(chan bool)
	serviceManager := manager.NewServicemanager()
	serviceManager.Init()
	serviceManager.RegisterService()
	serviceManager.InitService()
	serviceManager.HeartBeat()
	serviceManager.StartService()
	fmt.Println("Agent starts successfully.")
	logs.GetLogger().Info("This is agent main")

	<-start
}

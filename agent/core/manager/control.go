package manager

import (
	"github.com/huaweicloud/telescope/agent/core/heartbeat"
	"os"
	"time"

	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_config "github.com/huaweicloud/telescope/agent/core/lts/config"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
)

func HandleOsSignal(osSignal chan os.Signal) {
	sig := <-osSignal

	HandleSignalDirect(sig)
}

func HandleSignalDirect(signal os.Signal) {
	shutdownTime := time.Now().UnixNano() / 1000000
	var sigHb *channel.HBEntity
	if signal == upgrade.SIG_UPGRADE {
		sigHb = channel.NewHBEntity(channel.Upgrading, shutdownTime, lts_config.GetConfig().Enable, "", "")
		logs.GetLogger().Infof("Agent is upgrading at: %v, signal is %v", time.Now(), signal)
	} else {
		sigHb = channel.NewHBEntity(channel.Shutdown, shutdownTime, lts_config.GetConfig().Enable, "", "")
		logs.GetLogger().Infof("Shutdown agent at: %v, signal is %v", time.Now(), signal)
	}

	heartbeat.SendSignalHeartBeat(sigHb)
	os.Exit(0)
}

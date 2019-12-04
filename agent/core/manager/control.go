package manager

import (
	"os"
	"time"

	"github.com/huaweicloud/telescope/agent/core/heartbeat"

	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
)

// HandleOsSignal ...
func HandleOsSignal(osSignal chan os.Signal) {
	sig := <-osSignal

	HandleSignalDirect(sig)
}

// HandleSignalDirect ...
func HandleSignalDirect(signal os.Signal) {
	shutdownTime := time.Now().UnixNano() / 1000000
	var sigHb *channel.HBEntity
	if signal == upgrade.SIG_UPGRADE {
		sigHb = channel.NewHBEntity(channel.Upgrading, shutdownTime, "")
		logs.GetCesLogger().Infof("Agent is upgrading at: %v, signal is %v", time.Now(), signal)
	} else {
		sigHb = channel.NewHBEntity(channel.Shutdown, shutdownTime, "")
		logs.GetCesLogger().Infof("Shutdown agent at: %v, signal is %v", time.Now(), signal)
	}

	heartbeat.SendSignalHeartBeat(sigHb)
	os.Exit(0)
}

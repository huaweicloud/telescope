package manager

import (
	"os"
	"syscall"
	"testing"

	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/heartbeat"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//HandleOsSignal
func TestHandleOsSignal(t *testing.T) {
	//go InitConfig()
	Convey("Test_HandleOsSignal", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HandleSignalDirect, func(signal os.Signal) {
				return
			})
			signals := make(chan os.Signal, 1)
			signals <- syscall.SIGQUIT
			HandleOsSignal(signals)
		})
	})
}

//HandleSignalDirect
func TestHandleSignalDirect(t *testing.T) {
	//go InitConfig()
	Convey("Test_HandleSignalDirect", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(channel.NewHBEntity, func(status channel.StatusEnum, time int64, cesDetails string) *channel.HBEntity {
				return nil
			})
			mockfn.Replace(heartbeat.SendSignalHeartBeat, func(hb *channel.HBEntity) {
				return
			})
			mockfn.Replace(os.Exit, func(code int) {
				return
			})
			HandleSignalDirect(upgrade.SIG_UPGRADE)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(channel.NewHBEntity, func(status channel.StatusEnum, time int64, cesDetails string) *channel.HBEntity {
				return nil
			})
			mockfn.Replace(heartbeat.SendSignalHeartBeat, func(hb *channel.HBEntity) {
				return
			})
			mockfn.Replace(os.Exit, func(code int) {
				return
			})
			HandleSignalDirect(syscall.SIGQUIT)
		})
	})
}

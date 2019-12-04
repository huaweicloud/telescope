package manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
)

func HandleSignal() {
	logs.GetCesLogger().Info("Begin to handle os signal.")
	sigHandler := NewSignalSet()
	sigHandler.Register(syscall.SIGKILL, HandleSignalDirect)
	sigHandler.Register(syscall.SIGTERM, HandleSignalDirect)
	sigHandler.Register(upgrade.SIG_UPGRADE, HandleSignalDirect)
	sigHandler.Register(SigUserStop, ExitBySigUserStop)
	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan)

	for {
		select {
		case sig := <-sigChan:
			logs.GetCesLogger().Debugf("Begin to handle signal: %v", sig)
			err := sigHandler.Handle(sig)
			if err != nil {
				logs.GetCesLogger().Debugf("Handle sig:%v error: %s", sig, err.Error())
			}
		case <-time.After(time.Duration(3) * time.Second):
		}
	}
}

func ExitBySigUserStop(sig os.Signal) {
	if sig == SigUserStop {
		logs.GetCesLogger().Warnf("Agent will stop by SigUserStop. Bye")
		logs.GetCesLogger().Flush()
		os.Exit(0)
	}
	logs.GetCesLogger().Warn("Handle signal is not SigUserStop")
}

type SignalHandler func(s os.Signal)

type SignalSet struct {
	m map[os.Signal]SignalHandler
}

func NewSignalSet() *SignalSet {
	return &SignalSet{
		m: make(map[os.Signal]SignalHandler),
	}
}

func (set *SignalSet) Register(s os.Signal, handler SignalHandler) {
	if _, ok := set.m[s]; !ok {
		set.m[s] = handler
	}
}

func (set *SignalSet) Handle(sig os.Signal) (err error) {
	if _, ok := set.m[sig]; ok {
		set.m[sig](sig)
		return nil
	}

	return fmt.Errorf("no handler available for signal %v", sig)
}

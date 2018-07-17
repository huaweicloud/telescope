package manager

import (
	"os/signal"
	"syscall"

	"github.com/huaweicloud/telescope/agent/core/ces"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/heartbeat"
	"github.com/huaweicloud/telescope/agent/core/lts"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type servicemanager struct {
	serviceMap map[string]Service
}

func NewServicemanager() *servicemanager {
	sMap := make(map[string]Service)
	return &servicemanager{serviceMap: sMap}
}

func (sm *servicemanager) Init() {
	//init conf.json
	utils.InitConfig()
	//register and listen kill signal
	signal.Notify(getchOsSignal(), syscall.SIGKILL, syscall.SIGTERM, upgrade.SIG_UPGRADE)
	go HandleOsSignal(getchOsSignal())

}

func (sm *servicemanager) RegisterService() {
	sm.serviceMap["cesService"] = &ces.Service{}
	sm.serviceMap["ltsService"] = &lts.LtsService{}
}

func (sm *servicemanager) InitService() {
	for _, service := range sm.serviceMap {
		service.Init()
	}
}

func (sm *servicemanager) StartService() {
	for _, service := range sm.serviceMap {
		service.Start()
	}
}

func (sm *servicemanager) HeartBeat() {
	hb := heartbeat.HeartBeat{LtsDetails: "", CesDetails: ""}
	go hb.LoadHbServicesDetails(channel.GetServicesChData())
	go hb.ProduceHeartBeat(channel.GetHeartBeatChan())
	go hb.ConsumeHeartBeat(channel.GetHeartBeatChan())
}

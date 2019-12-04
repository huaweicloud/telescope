package manager

import (
	"github.com/huaweicloud/telescope/agent/core/assistant"
	"github.com/huaweicloud/telescope/agent/core/ces"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/heartbeat"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type servicemanager struct {
	serviceMap map[string]Service
}

// NewServicemanager ...
func NewServicemanager() *servicemanager {
	sMap := make(map[string]Service)
	return &servicemanager{serviceMap: sMap}
}

func (sm *servicemanager) Init() {
	//init conf.json
	utils.InitConfig()
	//register and listen kill signal
	go HandleSignal()

}

func (sm *servicemanager) RegisterService() {
	sm.serviceMap["cesService"] = &ces.Service{}
	sm.serviceMap["assistant"] = &assistant.Assistant{}
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
	hb := heartbeat.HeartBeat{CesDetails: ""}
	go hb.LoadHbServicesDetails(channel.GetServicesChData())
	go hb.ProduceHeartBeat(channel.GetHeartBeatChan())
	go hb.ConsumeHeartBeat(channel.GetHeartBeatChan())
}

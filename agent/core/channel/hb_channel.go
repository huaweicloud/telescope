package channel

import (
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/utils"
)

var heartBeatChan chan *HBEntity

// AgentStatus ...
type AgentStatus string

// StatusEnum ...
type StatusEnum int

const (
	// Running ...
	Running StatusEnum = iota
	Shutdown
	Upgrading
)

// MetaData ...
type MetaData struct {
	Version string      `json:"version"`
	OSName  string      `json:"os_name"`
	OSArch  string      `json:"os_arch"`
	CES     CesMetaData `json:"CES"`
}

// CesMetaData ...
type CesMetaData struct {
	Enable bool   `json:"enable"`
	Detail string `json:"detail"`
}

// HBEntity ...
type HBEntity struct {
	InstanceId string      `json:"instance_id"`
	Status     AgentStatus `json:"status"`
	TimeStamp  int64       `json:"time"`
	MetaData   MetaData    `json:"meta_data"`
}

// HBResponse ...
type HBResponse struct {
	Version     string `json:"version"`
	DownloadUrl string `json:"download_url"`
	CesConfig   string `json:"ces_config"`
	Md5         string `json:"md5"`
}

// Initialize the heartbeat channel
func init() {
	heartBeatChan = make(chan *HBEntity, 1)
}

// GetHeartBeatChan Get the heartbeat channel
func GetHeartBeatChan() chan *HBEntity {
	return heartBeatChan
}

// NewHBEntity ...
func NewHBEntity(status StatusEnum, time int64, cesDetails string) *HBEntity {
	cesMetaData := CesMetaData{Detail: cesDetails}
	cesMetaData.Enable = true

	hbEntity := &HBEntity{
		InstanceId: utils.GetConfig().InstanceId,
		Status:     status.MapStatus(),
		TimeStamp:  time,
		MetaData: MetaData{
			Version: utils.AgentVersion,
			OSName:  runtime.GOOS,
			OSArch:  runtime.GOARCH,
			CES:     cesMetaData,
		},
	}
	return hbEntity
}

// MapStatus ...
func (status StatusEnum) MapStatus() AgentStatus {
	switch status {
	case Running:
		return "running"
	case Shutdown:
		return "stopped"
	case Upgrading:
		return "upgrading"
	default:
		return "unknown"
	}
}

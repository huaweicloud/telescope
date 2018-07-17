package channel

import (
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/utils"
)

var heartBeatChan chan *HBEntity

type AgentStatus string
type StatusEnum int

const (
	Running StatusEnum = iota
	Shutdown
	Upgrading
)

type MetaData struct {
	Version string      `json:"version"`
	OSName  string      `json:"os_name"`
	OSArch  string      `json:"os_arch"`
	LTS     LtsMetaData `json:"LTS"`
	CES     CesMetaData `json:"CES"`
}

type LtsMetaData struct {
	Enable bool   `json:"enable"`
	Detail string `json:"detail"`
}

type CesMetaData struct {
	Enable bool   `json:"enable"`
	Detail string `json:"detail"`
}

type HBEntity struct {
	InstanceId string      `json:"instance_id"`
	Status     AgentStatus `json:"status"`
	TimeStamp  int64       `json:"time"`
	MetaData   MetaData    `json:"meta_data"`
}

type HBResponse struct {
	Version     string `json:"version"`
	DownloadUrl string `json:"download_url"`
	LtsConfig   string `json:"lts_config"`
	CesConfig   string `json:"ces_config"`
	Md5         string `json:"md5"`
}

// Initialize the heartbeat channel
func init() {
	heartBeatChan = make(chan *HBEntity, 1)
}

// Get the heartbeat channel
func GetHeartBeatChan() chan *HBEntity {
	return heartBeatChan
}

// New HBEntity
func NewHBEntity(status StatusEnum, time int64, ltsEnable bool, ltsDetails string, cesDetails string) *HBEntity {
	ltsMetaData := LtsMetaData{Enable: ltsEnable, Detail: ltsDetails}
	cesMetaData := CesMetaData{Detail: cesDetails}
	cesMetaData.Enable = true
	ltsMetaData.Enable = true

	metaData := MetaData{Version: utils.AGENT_VERSION, OSName: runtime.GOOS, OSArch: runtime.GOARCH,
		LTS: ltsMetaData, CES: cesMetaData}
	hbEntity := &HBEntity{
		InstanceId: utils.GetConfig().InstanceId,
		Status:     status.MapStatus(),
		TimeStamp:  time,
		MetaData:   metaData,
	}
	return hbEntity
}

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

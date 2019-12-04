package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	coreUtils "github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/net"
)

// NetStates is the type for store net state
type NetStates struct {
	byteSent        float64
	byteRecv        float64
	packetSent      float64
	packetRecv      float64
	errIn           float64
	errOut          float64
	dropIn          float64
	dropOut         float64
	collectTime     int64
	uptimeInSeconds int64
}

// NetCollector is the collector type for net metric
type NetCollector struct {
	LastStates *NetStates
}

// Collect implement the net Collector
func (n *NetCollector) Collect(collectTime int64) *model.InputMetric {
	var deltaTime = float64(coreUtils.DefaultMetricDeltaDataTimeInSecond)
	var packetErrInRate, packetErrOutRate, packetDropInRate, packetDropOutRate float64
	netStates, err := net.IOCounters(false)
	if nil != err {
		logs.GetCesLogger().Errorf("get net io count error %v", err)
		return nil
	}

	currStates := &NetStates{
		byteSent:    float64(netStates[0].BytesSent),
		byteRecv:    float64(netStates[0].BytesRecv),
		packetSent:  float64(netStates[0].PacketsSent),
		packetRecv:  float64(netStates[0].PacketsRecv),
		errIn:       float64(netStates[0].Errin),
		errOut:      float64(netStates[0].Errout),
		dropIn:      float64(netStates[0].Dropin),
		dropOut:     float64(netStates[0].Dropout),
		collectTime: collectTime,
	}
	currStates.uptimeInSeconds, _ = utils.GetUptimeInSeconds()

	if n.LastStates == nil {
		n.LastStates = currStates
		return nil
	}

	deltaTimeUsingCT := float64(currStates.collectTime-n.LastStates.collectTime) / 1000
	if currStates.uptimeInSeconds != -1 && n.LastStates.uptimeInSeconds != -1 {
		deltaTime = float64(currStates.uptimeInSeconds - n.LastStates.uptimeInSeconds)
	} else if deltaTimeUsingCT > 0 {
		deltaTime = deltaTimeUsingCT
	}

	totalSentPacket := utils.Float64From32Bits(currStates.packetSent - n.LastStates.packetSent)
	totalRecvPacket := utils.Float64From32Bits(currStates.packetRecv - n.LastStates.packetRecv)
	bitRecvRate := utils.Float64From32Bits(currStates.byteRecv-n.LastStates.byteRecv) / float64(deltaTime) * 8
	bitSentRate := utils.Float64From32Bits(currStates.byteSent-n.LastStates.byteSent) / float64(deltaTime) * 8
	packetSentRate := utils.Float64From32Bits(currStates.packetSent-n.LastStates.packetSent) / float64(deltaTime)
	packetRecvRate := utils.Float64From32Bits(currStates.packetRecv-n.LastStates.packetRecv) / float64(deltaTime)

	if totalRecvPacket != 0 {
		packetErrInRate = 100 * utils.Float64From32Bits(currStates.errIn-n.LastStates.errIn) / totalRecvPacket / float64(deltaTime)
		packetDropInRate = 100 * utils.Float64From32Bits(currStates.dropIn-n.LastStates.dropIn) / totalRecvPacket / float64(deltaTime)
	}
	if totalSentPacket != 0 {
		packetErrOutRate = 100 * utils.Float64From32Bits(currStates.errOut-n.LastStates.errOut) / totalSentPacket / float64(deltaTime)
		packetDropOutRate = 100 * utils.Float64From32Bits(currStates.dropOut-n.LastStates.dropOut) / totalSentPacket / float64(deltaTime)
	}

	logs.GetCesLogger().Debugf("bitRecvRate: %v bits/s, bitSentRate: %v bits/s, packetSentRate: %v Counts/s, packetRecvRate: %v Counts/s, collectime: %v", bitRecvRate, bitSentRate, packetSentRate, packetRecvRate, collectTime)

	n.LastStates = currStates

	fieldsG := []model.Metric{
		{
			MetricName:  "net_bitSent",
			MetricValue: bitRecvRate,
		},
		{
			MetricName:  "net_bitRecv",
			MetricValue: bitSentRate,
		},
		{
			MetricName:  "net_packetSent",
			MetricValue: packetSentRate,
		},
		{
			MetricName:  "net_packetRecv",
			MetricValue: packetRecvRate,
		},
		{
			MetricName:  "net_errin",
			MetricValue: packetErrInRate,
		},
		{
			MetricName:  "net_errout",
			MetricValue: packetErrOutRate,
		},
		{
			MetricName:  "net_dropin",
			MetricValue: packetDropInRate,
		},
		{
			MetricName:  "net_dropout",
			MetricValue: packetDropOutRate,
		},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "net",
		CollectTime: collectTime,
	}
}

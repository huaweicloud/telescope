package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/net"

	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
)

// NetStates is the type for store net state
type NetStates struct {
	byteSent    float64
	byteRecv    float64
	packetSent  float64
	packetRecv  float64
	errin       float64
	errout      float64
	dropin      float64
	dropout     float64
	collectTime int64
	uptimeInSeconds int64
}

// NetCollector is the collector type for net metric
type NetCollector struct {
	LastStates *NetStates
}

const UINT32_MAX = ^uint32(0)

// Collect implement the net Collector
func (n *NetCollector) Collect(collectTime int64) *model.InputMetric {

	var deltaTime = float64(ces_utils.DEFAULT_DELTA_TIME_IN_SECONDS)
	var result model.InputMetric
	var packetErrInRate float64 = 0.0
	var packetErrOutRate float64 = 0.0
	var packetDropInRate float64 = 0.0
	var packetDropOutRate float64 = 0.0
	netStates, _ := net.IOCounters(false)

	currStates := new(NetStates)
	currStates.byteRecv = float64(netStates[0].BytesRecv)
	currStates.byteSent = float64(netStates[0].BytesSent)
	currStates.packetSent = float64(netStates[0].PacketsSent)
	currStates.packetRecv = float64(netStates[0].PacketsRecv)

	currStates.errin = float64(netStates[0].Errin)
	currStates.errout = float64(netStates[0].Errout)
	currStates.dropin = float64(netStates[0].Dropin)
	currStates.dropout = float64(netStates[0].Dropout)

	currStates.collectTime = collectTime
	currStates.uptimeInSeconds, _ = ces_utils.GetUptimeInSeconds()

	if n.LastStates == nil {
		n.LastStates = currStates
		return nil
	}

	totalSentPacket := currStates.packetSent - n.LastStates.packetSent
	totalRecvPacket := currStates.packetRecv - n.LastStates.packetRecv

	deltaTimeUsingCT := float64(currStates.collectTime - n.LastStates.collectTime) / 1000
	if currStates.uptimeInSeconds != -1 && n.LastStates.uptimeInSeconds != -1{
		deltaTime = float64(currStates.uptimeInSeconds - n.LastStates.uptimeInSeconds)
	}else if deltaTime > 0{
		deltaTime = deltaTimeUsingCT
	}

	// windows can only support UINT32 for bytesRecv and bytesSent
	deltaByteRecv := currStates.byteRecv - n.LastStates.byteRecv
	deltaByteSent := currStates.byteSent - n.LastStates.byteSent
	if deltaByteRecv < 0 {
		deltaByteRecv = deltaByteRecv + float64(UINT32_MAX)
	}
	if deltaByteSent < 0 {
		deltaByteSent = deltaByteSent + float64(UINT32_MAX)
	}

	bitRecvRate := deltaByteRecv / float64(deltaTime) * 8
	bitSentRate := deltaByteSent / float64(deltaTime) * 8
	packetSentRate := (currStates.packetSent - n.LastStates.packetSent) / float64(deltaTime)
	packetRecvRate := (currStates.packetRecv - n.LastStates.packetRecv) / float64(deltaTime)

	if totalRecvPacket != 0 {
		packetErrInRate = 100 * (currStates.errin - n.LastStates.errin) / totalRecvPacket / float64(deltaTime)
		packetDropInRate = 100 * (currStates.dropin - n.LastStates.dropin) / totalRecvPacket / float64(deltaTime)
	}
	if totalSentPacket != 0 {
		packetErrOutRate = 100 * (currStates.errout - n.LastStates.errout) / totalSentPacket / float64(deltaTime)
		packetDropOutRate = 100 * (currStates.dropout - n.LastStates.dropout) / totalSentPacket / float64(deltaTime)
	}

	logs.GetCesLogger().Debugf("bitRecvRate: %v bits/s, bitSentRate: %v bits/s, packetSentRate: %v Counts/s, packetRecvRate: %v Counts/s, collectime: %v", bitRecvRate, bitSentRate, packetSentRate, packetRecvRate, collectTime)

	n.LastStates = currStates

	fieldsG := []model.Metric{
		model.Metric{MetricName: "net_bitSent", MetricValue: bitRecvRate},
		model.Metric{MetricName: "net_bitRecv", MetricValue: bitSentRate},
		model.Metric{MetricName: "net_packetSent", MetricValue: packetSentRate},
		model.Metric{MetricName: "net_packetRecv", MetricValue: packetRecvRate},
		model.Metric{MetricName: "net_errin", MetricValue: packetErrInRate},
		model.Metric{MetricName: "net_errout", MetricValue: packetErrOutRate},
		model.Metric{MetricName: "net_dropin", MetricValue: packetDropInRate},
		model.Metric{MetricName: "net_dropout", MetricValue: packetDropOutRate},
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

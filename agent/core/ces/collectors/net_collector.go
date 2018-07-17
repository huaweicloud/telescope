package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/net"
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
}

// NetCollector is the collector type for net metric
type NetCollector struct {
	LastStates *NetStates
}

const UINT32_MAX = ^uint32(0)

// Collect implement the net Collector
func (n *NetCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric
	var packetErrInRate float64 = 0.0
	var packetErrOutRate float64 = 0.0
	var packetDropInRate float64 = 0.0
	var packetDropOutRate float64 = 0.0
	netStates, _ := net.IOCounters(false)

	nowStates := new(NetStates)
	nowStates.byteRecv = float64(netStates[0].BytesRecv)
	nowStates.byteSent = float64(netStates[0].BytesSent)
	nowStates.packetSent = float64(netStates[0].PacketsSent)
	nowStates.packetRecv = float64(netStates[0].PacketsRecv)

	nowStates.errin = float64(netStates[0].Errin)
	nowStates.errout = float64(netStates[0].Errout)
	nowStates.dropin = float64(netStates[0].Dropin)
	nowStates.dropout = float64(netStates[0].Dropout)

	nowStates.collectTime = collectTime

	if n.LastStates == nil {
		n.LastStates = nowStates
		return nil
	}

	totalSentPacket := nowStates.packetSent - n.LastStates.packetSent
	totalRecvPacket := nowStates.packetRecv - n.LastStates.packetRecv

	secondDuration := (nowStates.collectTime - n.LastStates.collectTime) / 1000

	// windows can only support UINT32 for bytesRecv and bytesSent
	deltaByteRecv := nowStates.byteRecv - n.LastStates.byteRecv
	deltaByteSent := nowStates.byteSent - n.LastStates.byteSent
	if deltaByteRecv < 0 {
		deltaByteRecv = deltaByteRecv + float64(UINT32_MAX)
	}
	if deltaByteSent < 0 {
		deltaByteSent = deltaByteSent + float64(UINT32_MAX)
	}

	bitRecvRate := deltaByteRecv / float64(secondDuration) * 8
	bitSentRate := deltaByteSent / float64(secondDuration) * 8
	packetSentRate := (nowStates.packetSent - n.LastStates.packetSent) / float64(secondDuration)
	packetRecvRate := (nowStates.packetRecv - n.LastStates.packetRecv) / float64(secondDuration)

	if totalRecvPacket != 0 {
		packetErrInRate = (nowStates.errin - n.LastStates.errin) / totalRecvPacket / float64(secondDuration)
		packetDropInRate = (nowStates.dropin - n.LastStates.dropin) / totalRecvPacket / float64(secondDuration)
	}
	if totalSentPacket != 0 {
		packetErrOutRate = (nowStates.errout - n.LastStates.errout) / totalSentPacket / float64(secondDuration)
		packetDropOutRate = (nowStates.dropout - n.LastStates.dropout) / totalSentPacket / float64(secondDuration)
	}

	logs.GetCesLogger().Debugf("bitRecvRate: %v bits/s, bitSentRate: %v bits/s, packetSentRate: %v Counts/s, packetRecvRate: %v Counts/s, collectime: %v", bitRecvRate, bitSentRate, packetSentRate, packetRecvRate, collectTime)

	n.LastStates = nowStates

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

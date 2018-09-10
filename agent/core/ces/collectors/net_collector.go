// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
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
	return nil
}

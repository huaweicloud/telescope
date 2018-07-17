// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// MemCollector is the collector type for memory metric
type MemCollector struct {
}

// Collect implement the memory Collector
func (m *MemCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

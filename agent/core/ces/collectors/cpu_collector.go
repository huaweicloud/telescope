// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// CPUCollector is the collector type for cpu metric
type CPUCollector struct {
}

// Collect implement the cpu Collector
func (c *CPUCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

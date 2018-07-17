// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// LoadCollector is the collector type for cpu load metric
type LoadCollector struct {
}

// Collect implement the load Collector
func (l *LoadCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

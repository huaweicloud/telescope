// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// ProcessCollector is the collector type for process metric
type ProcessCollector struct {
}

// Collect implement the process Collector
func (p *ProcessCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

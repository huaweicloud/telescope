// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// ProcStatusCollector is the collector type for memory metric
type ProcStatusCollector struct {
}

// Collect implement the process status count Collector
func (p *ProcStatusCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

// +build !linux,!windows

package collectors

import "github.com/huaweicloud/telescope/agent/core/ces/model"

// ProcStatusCollector is the collector type for memory metric
type SpeProcCountCollector struct {
	CmdLines []string
}

// Collect implement the process status count Collector
func (p *SpeProcCountCollector) Collect(collectTime int64) *model.InputMetric {
	return &model.InputMetric{
		Type:        "cmdline",
		CollectTime: collectTime,
	}
}

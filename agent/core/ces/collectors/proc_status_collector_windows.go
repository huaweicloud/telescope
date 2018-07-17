package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/shirou/gopsutil/process"
)

// ProcStatusCollector is the collector type for memory metric
type ProcStatusCollector struct {
}

// Collect implement the process status count Collector
func (p *ProcStatusCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric
	allProcesses, _ := process.Processes()

	fieldsG := []model.Metric{
		model.Metric{MetricName: "proc_total_count", MetricValue: float64(len(allProcesses))},
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

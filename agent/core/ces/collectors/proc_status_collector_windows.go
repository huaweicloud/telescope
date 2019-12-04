package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/process"
)

// ProcStatusCollector is the collector type for memory metric
type ProcStatusCollector struct {
}

// Collect implement the process status count Collector
func (p *ProcStatusCollector) Collect(collectTime int64) *model.InputMetric {
	allPids, err := process.Pids()
	if nil != err {
		logs.GetCesLogger().Infof("get process status error %v", err)
		return nil
	}

	fieldsG := []model.Metric{
		{MetricName: "proc_total_count", MetricValue: float64(len(allPids))},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "process_total",
		CollectTime: collectTime,
	}
}

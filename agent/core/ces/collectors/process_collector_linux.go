package collectors

import (
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/shirou/gopsutil/process"
)

// ProcessCollector is the collector type for process metric
type ProcessCollector struct {
	Process *process.Process
}

// Collect implement the process Collector
func (p *ProcessCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	process := p.Process
	pName, _ := process.Name()
	pHashID := model.GenerateHashID(pName, process.Pid)

	processCPU, _ := process.Percent(time.Second)
	processMem, _ := process.MemoryPercent()
	openfiles, _ := process.OpenFiles()
	processOpenfiles := len(openfiles)

	fieldsG := []model.Metric{
		model.Metric{MetricName: "proc_cpu", MetricValue: float64(processCPU), MetricPrefix: pHashID},
		model.Metric{MetricName: "proc_mem", MetricValue: float64(processMem), MetricPrefix: pHashID},
		model.Metric{MetricName: "proc_file", MetricValue: float64(processOpenfiles), MetricPrefix: pHashID},
	}

	result.Data = fieldsG
	result.CollectTime = collectTime
	return &result
}

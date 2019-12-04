package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/process"
	"time"
)

// ProcessCollector is the collector type for process metric
type ProcessCollector struct {
	Process *process.Process
}

// Collect implement the process Collector
func (p *ProcessCollector) Collect(collectTime int64) *model.InputMetric {
	proc := p.Process
	pName, err := proc.Name()
	if nil != err {
		logs.GetCesLogger().Errorf("get process name error %v", err)
		return nil
	}
	pHashID := model.GenerateHashID(pName, proc.Pid)

	processCPU, err := proc.Percent(time.Second)
	if nil != err {
		logs.GetCesLogger().Errorf("get process cpu percent error %v", err)
		return nil
	}
	processMem, err := proc.MemoryPercent()
	if nil != err {
		logs.GetCesLogger().Errorf("get process memory percent error %v", err)
		return nil
	}

	fieldsG := []model.Metric{
		{MetricName: "proc_cpu", MetricValue: float64(processCPU), MetricPrefix: pHashID},
		{MetricName: "proc_mem", MetricValue: float64(processMem), MetricPrefix: pHashID},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "process",
		CollectTime: collectTime,
	}
}

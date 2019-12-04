package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/mem"
)

// MemCollector is the collector type for memory metric
type MemCollector struct {
}

// Collect implement the memory Collector
func (m *MemCollector) Collect(collectTime int64) *model.InputMetric {
	vm, err := mem.VirtualMemory()
	if nil != err {
		logs.GetCesLogger().Infof("get memory status error %v", err)
		return nil
	}

	fieldsG := []model.Metric{
		{MetricName: "mem_available", MetricValue: float64(vm.Available) / model.GBConversion},
		{MetricName: "mem_usedPercent", MetricValue: float64(vm.UsedPercent)},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "mem",
		CollectTime: collectTime,
	}
}

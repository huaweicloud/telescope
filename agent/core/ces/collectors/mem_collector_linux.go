package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/shirou/gopsutil/mem"
)

// MemCollector is the collector type for memory metric
type MemCollector struct {
}

// Collect implement the memory Collector
func (m *MemCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	vm, _ := mem.VirtualMemory()

	fieldsG := []model.Metric{
		model.Metric{MetricName: "mem_available", MetricValue: float64(vm.Available) / model.GBConversion},
		model.Metric{MetricName: "mem_usedPercent", MetricValue: float64(vm.UsedPercent)},
		model.Metric{MetricName: "mem_free", MetricValue: float64(vm.Free) / model.GBConversion},
		model.Metric{MetricName: "mem_buffers", MetricValue: float64(vm.Buffers) / model.GBConversion},
		model.Metric{MetricName: "mem_cached", MetricValue: float64(vm.Cached) / model.GBConversion},
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

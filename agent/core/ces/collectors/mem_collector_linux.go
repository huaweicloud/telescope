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
		logs.GetCesLogger().Errorf("Get memory stats failed and error is:%v", err)
		return nil
	}

	fieldsG := []model.Metric{
		{
			MetricName:  "mem_available",
			MetricValue: float64(vm.Available) / model.GBConversion,
		},
		{
			MetricName:  "mem_usedPercent",
			MetricValue: float64(vm.Total-vm.Available) / float64(vm.Total) * 100,
		},
		{
			MetricName:  "mem_free",
			MetricValue: float64(vm.Free) / model.GBConversion,
		},
		{
			MetricName:  "mem_buffers",
			MetricValue: float64(vm.Buffers) / model.GBConversion,
		},
		{
			MetricName:  "mem_cached",
			MetricValue: float64(vm.Cached) / model.GBConversion,
		},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "mem",
		CollectTime: collectTime,
	}
}

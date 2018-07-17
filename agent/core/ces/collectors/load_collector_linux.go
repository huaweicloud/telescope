package collectors

import (
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/shirou/gopsutil/load"
)

// LoadCollector is the collector type for cpu load metric
type LoadCollector struct {
}

// Collect implement the load Collector
func (l *LoadCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	loadAvg, _ := load.Avg()

	numCPU := float64(runtime.NumCPU())

	fieldsG := []model.Metric{
		model.Metric{MetricName: "load_average1", MetricValue: float64(loadAvg.Load1 / numCPU)},
		model.Metric{MetricName: "load_average5", MetricValue: float64(loadAvg.Load5 / numCPU)},
		model.Metric{MetricName: "load_average15", MetricValue: float64(loadAvg.Load15 / numCPU)},
	}
	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

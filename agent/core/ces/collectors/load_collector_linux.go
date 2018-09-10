package collectors

import (
	"runtime"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/load"
)

// LoadCollector is the collector type for cpu load metric
type LoadCollector struct {
}

// Collect implement the load Collector
func (l *LoadCollector) Collect(collectTime int64) *model.InputMetric {
	loadAvg, err := load.Avg()
	if nil != err {
		logs.GetCesLogger().Infof("get load error %v", err)
		return nil
	}

	numCPU := float64(runtime.NumCPU())

	fieldsG := []model.Metric{
		{MetricName: "load_average1", MetricValue: float64(loadAvg.Load1 / numCPU)},
		{MetricName: "load_average5", MetricValue: float64(loadAvg.Load5 / numCPU)},
		{MetricName: "load_average15", MetricValue: float64(loadAvg.Load15 / numCPU)},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "load",
		CollectTime: collectTime,
	}
}

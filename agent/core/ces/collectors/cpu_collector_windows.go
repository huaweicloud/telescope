package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/cpu"
)

// CPUStates is the type for store cpu state
type CPUStates struct {
	user         float64
	guest        float64
	system       float64
	idle         float64
	other        float64
	totalCPUTime float64
}

// CPUCollector is the collector type for cpu metric
type CPUCollector struct {
	LastStates *CPUStates
}

func getTotalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Idle
	return total
}

// Collect implement the cpu Collector
func (c *CPUCollector) Collect(collectTime int64) *model.InputMetric {
	cpuStats, err := cpu.Times(false)
	if nil != err || len(cpuStats) == 0 {
		logs.GetCesLogger().Errorf("get cpu stat error %v", err)
		return nil
	}

	stat := cpuStats[0]
	nowStates := &CPUStates{
		user:         stat.User,
		guest:        stat.Guest,
		system:       stat.System,
		idle:         stat.Idle,
		other:        1 - (stat.User - stat.Guest) - stat.System - stat.Idle,
		totalCPUTime: getTotalCPUTime(stat),
	}

	if c.LastStates == nil {
		c.LastStates = nowStates
		return nil
	}

	totalCPUTime := getTotalCPUTime(stat)
	totalDelta := totalCPUTime - c.LastStates.totalCPUTime

	cpuUsagUser := 100 * (nowStates.user - c.LastStates.user - (nowStates.guest - c.LastStates.guest)) / totalDelta
	cpuUsagSystem := 100 * (nowStates.system - c.LastStates.system) / totalDelta
	cpuUsagIdle := 100 * (nowStates.idle - c.LastStates.idle) / totalDelta

	c.LastStates = nowStates

	fieldsG := []model.Metric{
		{MetricName: "cpu_usage_user", MetricValue: cpuUsagUser},
		{MetricName: "cpu_usage_system", MetricValue: cpuUsagSystem},
		{MetricName: "cpu_usage_idle", MetricValue: cpuUsagIdle},
		{MetricName: "cpu_usage_other", MetricValue: utils.GetNonNegative(100 - cpuUsagUser - cpuUsagSystem - cpuUsagIdle)},
		{MetricName: "cpu_usage", MetricValue: utils.GetNonNegative(100 - cpuUsagIdle)},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "cpu",
		CollectTime: collectTime,
	}
}

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
	nice         float64
	iowait       float64
	irq          float64
	softirq      float64
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
		nice:         stat.Nice,
		iowait:       stat.Iowait,
		irq:          stat.Irq,
		softirq:      stat.Softirq,
		totalCPUTime: getTotalCPUTime(stat),
	}

	if c.LastStates == nil {
		c.LastStates = nowStates
		return nil
	}

	totalCPUTime := getTotalCPUTime(stat)
	totalDelta := totalCPUTime - c.LastStates.totalCPUTime

	cpuUsageUser := 100 * (nowStates.user - c.LastStates.user - (nowStates.guest - c.LastStates.guest)) / totalDelta
	cpuUsageSystem := 100 * (nowStates.system - c.LastStates.system) / totalDelta
	cpuUsageIdle := 100 * (nowStates.idle - c.LastStates.idle) / totalDelta

	cpuUsageNice := 100 * (nowStates.nice - c.LastStates.nice) / totalDelta
	cpuUsageIOWait := 100 * (nowStates.iowait - c.LastStates.iowait) / totalDelta
	cpuUsageIrq := 100 * (nowStates.irq - c.LastStates.irq) / totalDelta
	cpuUsageSoftIrq := 100 * (nowStates.softirq - c.LastStates.softirq) / totalDelta

	c.LastStates = nowStates

	fieldsG := []model.Metric{
		{MetricName: "cpu_usage_user", MetricValue: cpuUsageUser},
		{MetricName: "cpu_usage_system", MetricValue: cpuUsageSystem},
		{MetricName: "cpu_usage_idle", MetricValue: cpuUsageIdle},
		{MetricName: "cpu_usage_other", MetricValue: utils.GetNonNegative(100 - cpuUsageUser - cpuUsageSystem - cpuUsageIdle)},
		{MetricName: "cpu_usage", MetricValue: utils.GetNonNegative(100 - cpuUsageIdle)},

		{MetricName: "cpu_usage_nice", MetricValue: cpuUsageNice},
		{MetricName: "cpu_usage_iowait", MetricValue: cpuUsageIOWait},
		{MetricName: "cpu_usage_irq", MetricValue: cpuUsageIrq},
		{MetricName: "cpu_usage_softirq", MetricValue: cpuUsageSoftIrq},
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "cpu",
		CollectTime: collectTime,
	}
}

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
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
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal +
		t.Idle
	return total
}

// Collect implement the cpu Collector
func (c *CPUCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	cpuTimes, _ := cpu.Times(false)

	totalCPUTime := getTotalCPUTime(cpuTimes[0])

	nowStates := new(CPUStates)

	nowStates.user = cpuTimes[0].User
	nowStates.guest = cpuTimes[0].Guest
	nowStates.system = cpuTimes[0].System
	nowStates.idle = cpuTimes[0].Idle
	nowStates.other = 1 - (cpuTimes[0].User - cpuTimes[0].Guest) - cpuTimes[0].System - cpuTimes[0].Idle

	nowStates.nice = cpuTimes[0].Nice
	nowStates.iowait = cpuTimes[0].Iowait
	nowStates.irq = cpuTimes[0].Irq
	nowStates.softirq = cpuTimes[0].Softirq

	nowStates.totalCPUTime = getTotalCPUTime(cpuTimes[0])

	if c.LastStates == nil {
		c.LastStates = nowStates
		return nil
	}

	totalDelta := totalCPUTime - c.LastStates.totalCPUTime

	cpuUsagUser := 100 * (nowStates.user - c.LastStates.user - (nowStates.guest - c.LastStates.guest)) / totalDelta
	cpuUsagSystem := 100 * (nowStates.system - c.LastStates.system) / totalDelta
	cpuUsagIdle := 100 * (nowStates.idle - c.LastStates.idle) / totalDelta

	cpuUsagNice := 100 * (nowStates.nice - c.LastStates.nice) / totalDelta
	cpuUsagIOWait := 100 * (nowStates.iowait - c.LastStates.iowait) / totalDelta
	cpuUsagIrq := 100 * (nowStates.irq - c.LastStates.irq) / totalDelta
	cpuUsagSoftirq := 100 * (nowStates.softirq - c.LastStates.softirq) / totalDelta

	c.LastStates = nowStates

	fieldsG := []model.Metric{
		model.Metric{MetricName: "cpu_usage_user", MetricValue: cpuUsagUser},
		model.Metric{MetricName: "cpu_usage_system", MetricValue: cpuUsagSystem},
		model.Metric{MetricName: "cpu_usage_idle", MetricValue: cpuUsagIdle},
		model.Metric{MetricName: "cpu_usage_other", MetricValue: 100 - cpuUsagUser - cpuUsagSystem - cpuUsagIdle},
		model.Metric{MetricName:"cpu_usage", MetricValue: 100 - cpuUsagIdle},

		model.Metric{MetricName: "cpu_usage_nice", MetricValue: cpuUsagNice},
		model.Metric{MetricName: "cpu_usage_iowait", MetricValue: cpuUsagIOWait},
		model.Metric{MetricName: "cpu_usage_irq", MetricValue: cpuUsagIrq},
		model.Metric{MetricName: "cpu_usage_softirq", MetricValue: cpuUsagSoftirq},
	}
	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

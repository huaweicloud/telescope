package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/process"
)

// ProcStatusCollector is the collector type for memory metric
type ProcStatusCollector struct {
}

// ProcessNum the type for process info
type ProcessNum struct {
	runningProcNum  int
	idleProcNum     int
	zombieProcNum   int
	blockedProcNum  int
	sleepingProcNum int
}

// Collect implement the process status count Collector
func (p *ProcStatusCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric
	allProcesses, _ := process.Processes()
	procs := &ProcessNum{}

	for _, eachProcess := range allProcesses {
		status, _ := eachProcess.Status()
		switch status {
		case "R":
			procs.runningProcNum++
		case "S":
			procs.sleepingProcNum++
		case "Z":
			procs.zombieProcNum++
		case "I":
			procs.idleProcNum++
		case "B":
			procs.blockedProcNum++
		default:
			logs.GetCesLogger().Infof("Other status of process :%v,status :%v", eachProcess.Pid, status)
		}
	}

	fieldsG := []model.Metric{
		model.Metric{MetricName: "proc_running_count", MetricValue: float64(procs.runningProcNum)},
		model.Metric{MetricName: "proc_idle_count", MetricValue: float64(procs.idleProcNum)},
		model.Metric{MetricName: "proc_zombie_count", MetricValue: float64(procs.zombieProcNum)},
		model.Metric{MetricName: "proc_blocked_count", MetricValue: float64(procs.blockedProcNum)},
		model.Metric{MetricName: "proc_sleeping_count", MetricValue: float64(procs.sleepingProcNum)},
		model.Metric{MetricName: "proc_total_count", MetricValue: float64(len(allProcesses))},
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

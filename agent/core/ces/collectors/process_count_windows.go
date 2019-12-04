package collectors

import (
	"strings"

	"github.com/huaweicloud/telescope/agent/core/ces/gopsutil/process"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// ProcStatusCollector is the collector type for memory metric
type SpeProcCountCollector struct {
	CmdLines []string
}

// Collect implement the process status count Collector
func (p *SpeProcCountCollector) Collect(collectTime int64) *model.InputMetric {
	if len(p.CmdLines) == 0 {
		logs.GetCesLogger().Errorf("input process command line is empty")
		return nil
	}

	allProcesses, err := process.GetWin32Proc()
	if nil != err {
		logs.GetCesLogger().Errorf("Get all processes error %v", err)
		return nil
	}

	cmdCount, cmdID := make(map[string]int), make(map[string]string)
	for _, cmd := range p.CmdLines {
		cmdCount[cmd] = 0
		cmdID[cmd] = model.GenerateHashIDByPname(cmd)
	}

	logs.GetCesLogger().Errorf("allProcesses len is %v", len(allProcesses))
	for _, proc := range allProcesses {
		str := *proc.CommandLine
		for _, cmd := range p.CmdLines {
			if strings.Contains(str, cmd) {
				if v, ok := cmdCount[cmd]; ok {
				
					cmdCount[cmd] = v + 1
				}
			}
		}
	}
	var fieldsG []model.Metric
	for k, v := range cmdCount {
		fieldsG = append(fieldsG, model.Metric{MetricName: "proc_specified_count", MetricValue: float64(v), MetricPrefix: cmdID[k], CustomProcName: k})
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "cmdline",
		CollectTime: collectTime,
	}
}

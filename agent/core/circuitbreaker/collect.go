package circuitbreaker

import (
	"errors"
	"os"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/process"
)

func getCurrentProcess() (*process.Process, error) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		logs.GetCesLogger().Errorf("NewProcess failed and error is: %v", err)
	}

	return p, err
}

// getCPUPercent returns cpu percent like linux, using single processor as unit
func getCPUPercent(p *process.Process) (float64, error) {
	if p == nil {
		logs.GetCesLogger().Error("Get CPU percent from nil pointer")
		return 0, errors.New("nil pointer")
	}

	cpuPercent, err := p.Percent(0)
	if err != nil {
		logs.GetCesLogger().Errorf("Get CPU percent failed and error is: %v", err)
		return 0, err
	}

	logs.GetCesLogger().Debugf("Get CPU percent successfully and CPU percent is: %f", cpuPercent)
	return cpuPercent, nil
}

// getMemory returns RSS(resident size set)
func getMemory(p *process.Process) (uint64, error) {
	if p == nil {
		logs.GetCesLogger().Error("Get memory info from nil pointer")
		return 0, errors.New("nil pointer")
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		logs.GetCesLogger().Errorf("Get memory info failed and error is: %v", err)
		return 0, err
	}

	logs.GetCesLogger().Debugf("Get memory info successfully and memory info is: %v", memInfo)
	return memInfo.RSS, nil
}

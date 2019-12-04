package model

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/process"
)

// CPUProcess the type for top5 cpu process
type CPUProcess struct {
	Pid        int32   `json:"pid"`
	Pname      string  `json:"name"`
	CPU        float64 `json:"cpu"`
	Cmdline    string  `json:"cmdline"`
	CreateTime int64   `json:"create_time"`
	Process    *process.Process
}

// CPUProcessList the type for CPUProcess slice
type CPUProcessList []*CPUProcess

func (c CPUProcessList) String() string {
	var result string
	for _, v := range c {
		result += fmt.Sprintf(" {PID is %d, Pname is %s, CPU percent is %f, Cmdline is %s, CreateTime is %d}", v.Pid, v.Pname, v.CPU, v.Cmdline, v.CreateTime)
	}

	return result
}

// ProcessInfo the type for process info
type ProcessInfo struct {
	Pid        int32  `json:"pid"`
	Pname      string `json:"name"`
	Status     bool   `json:"status"`
	Cmdline    string `json:"cmdline"`
	CreateTime int64  `json:"create_time"`
}

// HashProcess the type for hash process
type HashProcess struct {
	ProcessHashID string      `json:"process_hashid"`
	Info          ProcessInfo `json:"info"`
}

// ProcessInfoDB the type for process info in request
type ProcessInfoDB struct {
	InstanceID string        `json:"instance_id"`
	Processes  []HashProcess `json:"processes"`
}

// ChProcessList used in channel
type ChProcessList []*ProcessInfo

func (c ChProcessList) String() string {
	var result string
	for _, v := range c {
		result += fmt.Sprintf(" {PID is %d, Pname is %s, Cmdline is %s, CreateTime is %d}", v.Pid, v.Pname, v.Cmdline, v.CreateTime)
	}

	return result
}

// GetTop5CpuProcessList get top5 cpu process list for channel
func GetTop5CpuProcessList() ChProcessList {
	logs.GetCesLogger().Debug("Enter GetTop5CpuProcessList")
	var top5CpuProcessList CPUProcessList
	allProcesses, err := process.Processes()
	if err != nil {
		logs.GetCesLogger().Errorf("GetTop5CpuProcessList get all process failed and error is: %v", err)
		return []*ProcessInfo{}
	}
	allProcessesNum := len(allProcesses)
	if allProcessesNum == 0 {
		logs.GetCesLogger().Warnf("GetTop5CpuProcessList returns all process as empty")
		return []*ProcessInfo{}
	}
	top5ChProcessList := make(ChProcessList, 0, allProcessesNum)

	logs.GetCesLogger().Debugf("GetTop5CpuProcessList get all process successfully, processes are: %s", func() string {
		var result string
		for _, p := range allProcesses {
			result += p.String()
		}
		return result
	}())
	logs.GetCesLogger().Debug("GetTop5CpuProcessList starts getting all process cpu percent")
	cpuProcessChan := make(chan *CPUProcess, allProcessesNum)
	wg := &sync.WaitGroup{}
	wg.Add(allProcessesNum)
	for _, eachProcess := range allProcesses {
		go func(p *process.Process) {
			defer wg.Done()

			pid := p.Pid
			logs.GetCesLogger().Debugf("GetTop5CpuProcessList in for and pid is: %d", pid)
			eachCPUProcessChan := make(chan *CPUProcess, 1)
			go func() {
				cpuPercent, cpuErr := p.Percent(time.Duration(cesUtils.Top5ProcessSamplePeriodInSeconds) * time.Second)

				if cpuErr != nil {
					logs.GetCesLogger().Errorf("GetTop5CpuProcessList get cpu percent failed(PID:%d) and error is: %v. Maybe the process has died in most cases.", pid, err)
					eachCPUProcessChan <- nil
					return
				}

				cpuProcess := &CPUProcess{
					Pid:     pid,
					CPU:     cpuPercent,
					Process: p,
				}
				eachCPUProcessChan <- cpuProcess
				logs.GetCesLogger().Debugf("GetTop5CpuProcessList get cpu percent successfully(PID:%d)", pid)
			}()
			select {
			case cpuProcess := <-eachCPUProcessChan:
				cpuProcessChan <- cpuProcess
				logs.GetCesLogger().Debugf("GetTop5CpuProcessList send cpuProcess to channel successfully(PID:%d).", pid)
			case <-time.After(time.Duration(cesUtils.Top5ProcessSamplePeriodInSeconds+3) * time.Second):
				logs.GetCesLogger().Errorf("GetTop5CpuProcessList get cpu percent timeout(PID:%d)", pid)
			}
		}(eachProcess)
	}
	wg.Wait()
	close(cpuProcessChan)
	for p := range cpuProcessChan {
		// 兼容采集进程CPU信息时进程已经退出的情况
		if p != nil {
			top5CpuProcessList = append(top5CpuProcessList, p)
		}
	}

	logs.GetCesLogger().Debugf("GetTop5CpuProcessList finish getting all process cpu percent, all process list is: %s", top5CpuProcessList.String())
	sort.Sort(top5CpuProcessList)
	var resultList CPUProcessList
	for _, p := range top5CpuProcessList {
		createTime, createTimeErr := p.Process.CreateTime()
		if createTimeErr != nil {
			continue
		}
		p.CreateTime = createTime
		resultList = append(resultList, p)
		if len(resultList) >= 5 {
			break
		}
	}

	for _, oneCPUProcess := range resultList {
		eachProcessInfo := new(ProcessInfo)
		eachProcessInfo.Pid = oneCPUProcess.Pid
		eachProcessInfo.Status = true
		eachProcessInfo.CreateTime = oneCPUProcess.CreateTime
		eachProcessInfo.Pname, _ = oneCPUProcess.Process.Name()
		eachProcessInfo.Cmdline, _ = oneCPUProcess.Process.Cmdline()

		if len(eachProcessInfo.Cmdline) > cesUtils.MaxCmdlineLen {
			eachProcessInfo.Cmdline = utils.SubStr(eachProcessInfo.Cmdline, cesUtils.MaxCmdlineLen-len(cesUtils.CmdlineSuffix))
			eachProcessInfo.Cmdline = utils.ConcatStr(eachProcessInfo.Cmdline, cesUtils.CmdlineSuffix)
		}

		top5ChProcessList = append(top5ChProcessList, eachProcessInfo)
	}

	logs.GetCesLogger().Debugf("Top 5 process list is: %v", top5ChProcessList.String())
	return top5ChProcessList
}

// used for sort by cpu desc
func (c CPUProcessList) Len() int {
	return len(c)
}

func (c CPUProcessList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CPUProcessList) Less(i, j int) bool {
	return c[i].CPU > c[j].CPU
}

// BuildProcessInfoByList used to build process info for api request
func BuildProcessInfoByList(processList ChProcessList) ProcessInfoDB {
	var processInfoDB ProcessInfoDB
	var processArr []HashProcess
	processInfoDB.InstanceID = utils.GetConfig().InstanceId

	for _, eachProcess := range processList {
		eachProcessInfo := new(HashProcess)

		eachProcessInfo.ProcessHashID = GenerateHashID(eachProcess.Pname, eachProcess.Pid)
		eachProcessInfo.Info = ProcessInfo{}
		eachProcessInfo.Info.Pid = eachProcess.Pid
		eachProcessInfo.Info.Pname = eachProcess.Pname
		eachProcessInfo.Info.Status = eachProcess.Status
		eachProcessInfo.Info.CreateTime = eachProcess.CreateTime
		eachProcessInfo.Info.Cmdline = eachProcess.Cmdline
		processArr = append(processArr, *eachProcessInfo)
	}
	processInfoDB.Processes = processArr

	return processInfoDB
}

// GenerateHashID generate hashid string by pname and pid
func GenerateHashID(pname string, pid int32) string {
	processStr := []byte(pname + strconv.Itoa(int(pid)))
	return fmt.Sprintf("%x", md5.Sum(processStr))
}

// GenerateHashIDByPname generate hashid string by pname
func GenerateHashIDByPname(pname string) string {
	processStr := []byte(pname)
	return fmt.Sprintf("%x", md5.Sum(processStr))
}

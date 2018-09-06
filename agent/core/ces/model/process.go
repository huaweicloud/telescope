package model

import (
	"sort"

	"github.com/shirou/gopsutil/process"

	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// CPUProcess the type for top5 cpu process
type CPUProcess struct {
	Pid        int32   `json:"pid"`
	Pname      string  `json:"name"`
	CPU        float64 `json:"cpu"`
	Cmdline    string  `json:"cmdline"`
	CreateTime int64   `json:"create_time"`
}

// CPUProcessList the type for CPUProcess slice
type CPUProcessList []*CPUProcess

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

// GetTop5CpuProcessList get top5 cpu process list for channel
func GetTop5CpuProcessList() ChProcessList {
	var top5CpuProcessList CPUProcessList
	var top5ChProcessList ChProcessList
	var createTimeErr error
	allProcesses, _ := process.Processes()

	for _, eachProcess := range allProcesses {
		eachProcessCPU, cpuErr := eachProcess.Percent(time.Second)
		if cpuErr != nil {
			continue
		}
		eachProcessID := eachProcess.Pid
		eachCPUProcess := new(CPUProcess)
		eachCPUProcess.Pid = eachProcessID
		eachCPUProcess.CPU = eachProcessCPU
		eachCPUProcess.Cmdline, _ = eachProcess.Cmdline()
		eachCPUProcess.CreateTime, createTimeErr = eachProcess.CreateTime()
		if createTimeErr != nil {
			continue
		}

		eachCPUProcess.Pname, _ = eachProcess.Name()
		top5CpuProcessList = append(top5CpuProcessList, eachCPUProcess)
	}
	sort.Sort(top5CpuProcessList)
	top5CpuProcessList = top5CpuProcessList[:5]

	for _, oneCPUProcess := range top5CpuProcessList {
		eachProcessInfo := new(ProcessInfo)
		eachProcessInfo.Pid = oneCPUProcess.Pid
		eachProcessInfo.Pname = oneCPUProcess.Pname
		eachProcessInfo.CreateTime = oneCPUProcess.CreateTime
		eachProcessInfo.Status = true
		eachProcessInfo.Cmdline = oneCPUProcess.Cmdline

		if len(eachProcessInfo.Cmdline) > cesUtils.MaxCmdlineLen {
			eachProcessInfo.Cmdline = utils.SubStr(eachProcessInfo.Cmdline, cesUtils.MaxCmdlineLen-len(cesUtils.CmdlineSuffix))
			eachProcessInfo.Cmdline = utils.ConcatStr(eachProcessInfo.Cmdline, cesUtils.CmdlineSuffix)
		}

		top5ChProcessList = append(top5ChProcessList, eachProcessInfo)
	}
	return top5ChProcessList
}

// used for sort by cpu desc
func (s CPUProcessList) Len() int {
	return len(s)
}

func (s CPUProcessList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s CPUProcessList) Less(i, j int) bool {
	return s[i].CPU > s[j].CPU
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

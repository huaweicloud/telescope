package circuitbreaker

import (
	"container/list"
	"fmt"
	"os"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/process"
)

// CircuitBreaker
type CircuitBreaker struct {
	//consecutiveFailures uint
	StateList *list.List
}

// State represents the state of cb
type State struct {
	CPUPct float64
	Memory uint64
	Time   time.Time
}

func (s *State) String() string {
	return fmt.Sprintf("{CPUPct: %f, Memory: %d, Time: %s}", s.CPUPct, s.Memory, s.Time.String())
}

var (
	p  *process.Process
	cb *CircuitBreaker
)

func init() {
	initProcess()
	initCB()
}

func initProcess() {
	var err error
	p, err = getCurrentProcess()
	if err != nil {
		p = nil
		logs.GetCesLogger().Error("Init process instance failed.")
		return
	}
	logs.GetCesLogger().Debugf("Init process instance(agent self) successfully and process instance is: %v", p)
}

func initCB() {
	cb = &CircuitBreaker{
		StateList: list.New(),
	}
}

// Start
func Start() {
	ticker := time.NewTicker(time.Duration(CBCheckIntervalInSecond) * time.Second)
	for range ticker.C {
		logs.GetCesLogger().Info("Circuit breaker checks fuse starts")
		go checkTheFuse()
	}
}

func checkTheFuse() {
	if p == nil {
		initProcess()
		logs.GetCesLogger().Error("Process(Agent self) is nil, checkTheFuse return")
		return
	}

	cpuPct, err := getCPUPercent(p)
	if err != nil {
		logs.GetCesLogger().Error("Get CPU pct failed, checkTheFuse return")
		return
	}

	memory, err := getMemory(p)
	if err != nil {
		logs.GetCesLogger().Error("Get memory failed, checkTheFuse return")
		return
	}

	cpuPctThreshold2, memoryThreshold2 := get2ndThreshold()
	if shouldTrip(cpuPct, cpuPctThreshold2, memory, memoryThreshold2) {
		logs.GetCesLogger().Warnf("Trip(2nd) occurs, cpu percent: %f(threshold: %f), memory: %d(threshold: %d)", cpuPct, cpuPctThreshold2, memory, memoryThreshold2)
		exec2ndTrip(cpuPct, memory, SecondTripType)
		return
	}

	cpuPctThreshold1, memoryThreshold1 := get1stThreshold()
	if shouldTrip(cpuPct, cpuPctThreshold1, memory, memoryThreshold1) {
		logs.GetCesLogger().Warnf("Trip(1st) occurs, cpu percent: %f(threshold: %f), memory: %d(threshold: %d)", cpuPct, cpuPctThreshold1, memory, memoryThreshold1)
		exec1stTrip(cpuPct, memory)
		return
	}

	logs.GetCesLogger().Infof("The fuse is ok, cpu percent: %f(threshold1: %f, threshold2: %f), memory: %d(threshold1: %d, threshold2: %d)", cpuPct, cpuPctThreshold1, cpuPctThreshold2, memory, memoryThreshold1, memoryThreshold2)
}

func shouldTrip(cpuPct, cpuPctThreshold float64, memory, memoryThreshold uint64) bool {
	return cpuPct >= cpuPctThreshold || memory >= memoryThreshold
}

func get1stThreshold() (float64, uint64) {
	var (
		cpuPctThreshold = utils.CPU1stPctThreshold
		memoryThreshold = utils.Memory1stThreshold
	)

	if utils.GetConfig() != nil && utils.GetConfig().CPU1stPctThreshold != 0 {
		cpuPctThreshold = utils.GetConfig().CPU1stPctThreshold
		logs.GetCesLogger().Infof("CPU1stPctThreshold is set to user defined value(%f)", cpuPctThreshold)
	} else {
		logs.GetCesLogger().Debugf("CPU1stPctThreshold is set to default value(%f)", cpuPctThreshold)
	}

	if utils.GetConfig() != nil && utils.GetConfig().Memory1stThreshold != 0 {
		memoryThreshold = utils.GetConfig().Memory1stThreshold
		logs.GetCesLogger().Infof("Memory1stThreshold is set to user defined value(%f)", memoryThreshold)
	} else {
		logs.GetCesLogger().Debugf("Memory1stThreshold is set to default value(%f)", memoryThreshold)
	}

	return cpuPctThreshold, memoryThreshold
}

func get2ndThreshold() (float64, uint64) {
	var (
		cpuPctThreshold = utils.CPU2ndPctThreshold
		memoryThreshold = utils.Memory2ndThreshold
	)

	if utils.GetConfig() != nil && utils.GetConfig().CPU2ndPctThreshold != 0 {
		cpuPctThreshold = utils.GetConfig().CPU2ndPctThreshold
		logs.GetCesLogger().Infof("CPU2ndPctThreshold is set to user defined value(%f)", cpuPctThreshold)
	} else {
		logs.GetCesLogger().Debugf("CPU2ndPctThreshold is set to default value(%f)", cpuPctThreshold)
	}

	if utils.GetConfig() != nil && utils.GetConfig().Memory2ndThreshold != 0 {
		memoryThreshold = utils.GetConfig().Memory2ndThreshold
		logs.GetCesLogger().Infof("Memory2ndThreshold is set to user defined value(%f)", memoryThreshold)
	} else {
		logs.GetCesLogger().Debugf("Memory2ndThreshold is set to default value(%f)", memoryThreshold)
	}

	return cpuPctThreshold, memoryThreshold
}

func exec1stTrip(cpuPct float64, memory uint64) {
	if cb.StateList.Len() == TripCountFor1stThreshold-1 {
		logs.GetCesLogger().Warnf("Trip(1st) has occurred for %s times, prepare for exit", TripCountFor1stThreshold)
		appendFile(cpuPct, memory, FirstTripType)
		programExit()
	}

	s := State{
		CPUPct: cpuPct,
		Memory: memory,
		Time:   time.Now(),
	}
	cb.StateList.PushBack(s)
	logs.GetCesLogger().Debugf("Trip(1st) add State(%s) and len is:%d", s.String(), cb.StateList.Len())

	dueTime := time.Now().Add(-time.Duration(TripCheckWindowInSecond) * time.Second)
	delOutdatedState(dueTime)
}

func exec2ndTrip(cpuPct float64, memory uint64, tripType uint) {
	appendFile(cpuPct, memory, tripType)
	programExit()
}

func programExit() {
	logs.GetCesLogger().Warn("Agent will exit due to fuse protection")
	logs.GetCesLogger().Flush()
	os.Exit(-1)
}

func resetList() {
	var next *list.Element
	for e := cb.StateList.Front(); e != nil; e = next {
		next = e.Next()
		cb.StateList.Remove(e)
	}
}

func delOutdatedState(dueTime time.Time) {
	var next *list.Element
	for e := cb.StateList.Front(); e != nil; e = next {
		next = e.Next()
		s, ok := e.Value.(State)
		if ok && !s.Time.Before(dueTime) {
			continue
		}
		cb.StateList.Remove(e)
		logs.GetCesLogger().Debugf("Delete state(%s) and dutTime is:%s", s.String(), dueTime.String())
	}
}

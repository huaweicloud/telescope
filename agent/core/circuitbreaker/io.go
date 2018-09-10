package circuitbreaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// Record is one item/row in record file
type Record struct {
	Time       string  `json:"time"`
	PID        int     `json:"pid"`
	CPUPercent float64 `json:"cpu_percent"`
	Memory     uint64  `json:"memory"`
	TripType   uint    `json:"trip_type"`
}

func (r *Record) String() string {
	return fmt.Sprintf("{Time: %s, PID: %d, CPUPercent: %f, Memory: %d, TripType: %d", r.Time, r.PID, r.CPUPercent, r.Memory, r.TripType)
}

var (
	cbFilePath string
)

func init() {
	workingPath := utils.GetWorkingPath()
	cbFilePath = filepath.Join(workingPath, AgentCircuitBreakerFileName)
}

func appendFile(cpuPct float64, memory uint64, tripType uint) {
	file, err := os.OpenFile(cbFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		logs.GetCesLogger().Errorf("Open/Create %s failed and error is: %v", cbFilePath, err)
		return
	}
	defer file.Close()

	if _, err = file.Write(getRecordBytes(cpuPct, memory, tripType)); err != nil {
		logs.GetCesLogger().Errorf("Write record to circuit breaker file(%s) failed, error is: %v", cbFilePath, err)
		return
	}
	logs.GetCesLogger().Infof("Write record to circuit breaker file(%s) successfully.", cbFilePath)
}

func getRecord(cpuPct float64, memory uint64, tripType uint) string {
	r := fmt.Sprintf("%-37s || %-7d || %-12f || %-11d || %-2d\n", time.Now().Format(TimeDefaultLayout), os.Getpid(), cpuPct, memory, tripType)
	logs.GetCesLogger().Debugf("Circuit breaker record is: %s", r)
	return r
}

func getRecordBytes(cpuPct float64, memory uint64, tripType uint) []byte {
	lf := []byte("\n")

	r := Record{
		Time:       time.Now().Format(TimeDefaultLayout),
		PID:        os.Getpid(),
		CPUPercent: cpuPct,
		Memory:     memory,
		TripType:   tripType,
	}

	recordBytes, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Marshall failed and error is: %v", err)
		return lf
	}

	recordBytes = append(recordBytes, lf...)
	return recordBytes
}

func readFile() []Record {
	var records []Record

	p, err := ioutil.ReadFile(cbFilePath)
	if err != nil {
		logs.GetCesLogger().Errorf("Read circuit breaker file(%s) failed and error is: %v", cbFilePath, err)
		return records
	}
	for _, line := range bytes.Split(p, []byte{'\n'}) {
		// ignore the last empty line
		if len(line) == 0 {
			continue
		}

		var r Record
		if err := json.Unmarshal(line, &r); err != nil {
			logs.GetCesLogger().Errorf("Unmarshal record failed, line is: %s, error is: %v", string(line), err)
			continue
		}
		records = append(records, r)
	}

	logs.GetCesLogger().Debugf("Read record(%v) successfully", records)
	return records
}

// DelFile
func DelFile() {
	if !utils.IsFileExist(cbFilePath) {
		logs.GetCesLogger().Warnf("Circuit breaker file(%s) does not exist", cbFilePath)
		return
	}

	err := os.Remove(cbFilePath)
	if err != nil {
		logs.GetCesLogger().Errorf("Delete circuit breaker file(%s) failed and error is: %v", cbFilePath, err)
	} else {
		logs.GetCesLogger().Infof("Delete circuit breaker file(%s) successfully.", cbFilePath)
	}
}

func getFileSize() int64 {
	if !utils.IsFileExist(cbFilePath) {
		logs.GetCesLogger().Warnf("Circuit breaker file(%s) does not exist", cbFilePath)
		return 0
	}

	fileInfo, err := os.Stat(cbFilePath)
	if err != nil {
		logs.GetCesLogger().Errorf("Get circuit breaker file(%s) stat failed and error is: %v", cbFilePath, err)
		return 0
	}

	logs.GetCesLogger().Debugf("Get circuit breaker file(%s) stat successfully and size is: %d bytes", fileInfo.Size())
	return fileInfo.Size()
}

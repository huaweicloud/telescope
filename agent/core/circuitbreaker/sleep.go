package circuitbreaker

import (
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"time"
)

// IsSleepNeeded
func IsSleepNeeded() bool {
	if !utils.IsFileExist(cbFilePath) {
		logs.GetCesLogger().Warnf("Circuit breaker record file(%s) does not exist", cbFilePath)
		return false
	}

	records := readFile()
	recordsLen := len(records)
	if recordsLen == 0 {
		logs.GetCesLogger().Warn("Circuit breaker record is empty")
		return false
	}

	var count int
	dueTime := time.Now().Add(-time.Duration(SleepCheckWindowInSecond) * time.Second)
	for i := recordsLen - 1; i >= 0; i-- {
		r := records[i]
		t, err := time.Parse(TimeDefaultLayout, r.Time)
		if err != nil {
			logs.GetCesLogger().Errorf("Record time parse failed and error is: %v", err)
			continue
		}
		if t.After(dueTime) {
			logs.GetCesLogger().Debugf("Record(%s) hits the due time", r.String())
			count ++
		}
	}
	if count >= SleepCheckCount {
		logs.GetCesLogger().Warnf("Record count(%d) exceeds or equals to the SleepCheckCount(%d)", count, SleepCheckCount)
		return true
	}

	logs.GetCesLogger().Debugf("Record count(%d) LESS THAN SleepCheckCount(%d)", count, SleepCheckCount)
	return false
}

func Sleep() {
	logs.GetCesLogger().Warnf("Agent will sleep for %d seconds due to fuse protection.", SleepTimeInSecond)
	logs.GetCesLogger().Flush()
	time.Sleep(time.Duration(SleepTimeInSecond) * time.Second)
}

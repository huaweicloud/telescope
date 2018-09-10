package task

import (
	"github.com/huaweicloud/telescope/agent/core/logs"
	"time"
)

/*
 * getTimeLocation returns the new cron string according to the local timezone
 * the cron in parameter express using UTC +08:00
 */
func getCSTLocation() *time.Location {
	location, err := time.LoadLocation(TIME_LOCATION_CST)
	if err != nil {
		logs.GetAssistantLogger().Errorf("Get time location failed and error is: ; just return local time location", err.Error())
		return time.Now().Location()
	}

	return location
}

func isFinalState(cronFlag bool, status string) bool {
	if cronFlag && status == STATE_CANCELED {
		return true
	}

	if !cronFlag && (status == STATE_CANCELED ||
		status == STATE_FAILED ||
		status == STATE_SUCCEEDED ||
		status == STATE_TIMEOUT) {
		return true
	}

	return false
}

package utils

import (
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"

	os_uptime "github.com/go-osstat/uptime"
)

func GetUptimeInSeconds()(int64, error){
	osUptime, err := os_uptime.Get()
	if err != nil{
		logs.GetCesLogger().Errorf("Get uptime failed and error is: %s", err.Error())
		return -1, err
	} else {
		return int64(osUptime / time.Second), nil
	}
}
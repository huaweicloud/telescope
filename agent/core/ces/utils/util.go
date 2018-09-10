package utils

import (
	"math"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/mackerelio/go-osstat/uptime"
)

// GetUptimeInSeconds ...
func GetUptimeInSeconds() (int64, error) {
	osUptime, err := uptime.Get()
	if err != nil {
		logs.GetCesLogger().Errorf("Get uptime failed and error is: %s", err.Error())
		return -1, err
	} else {
		return int64(osUptime / time.Second), nil
	}
}

//If x<0, returns the number of negative
func GetNonNegative(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

// Float64From32Bits 返回由于翻转导致出现负数时的实际差值 uint32
// a -> a' && a' < a ==> (a'-a)<0
func Float64From32Bits(f float64) float64 {
	if f < 0 {
		return float64(f + math.MaxUint32 + 1)
	}

	return f
}

// Float64From64Bits 返回由于翻转导致出现负数时的实际差值 uint64
func Float64From64Bits(f float64) float64 {
	if f < 0 {
		return float64(f + math.MaxUint64 + 1)
	}

	return f
}

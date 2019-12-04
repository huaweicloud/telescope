package circuitbreaker

const (
	FirstTripType uint = iota
	SecondTripType
)

const (
	TimeDefaultLayout           = "2006-01-02 15:04:05.999999999 -0700 MST"
	AgentCircuitBreakerFileName = "cesAgentCircuitBreaker"
	AgentPIDFileName            = "cesAgent.pid"
	CBCheckIntervalInSecond     = 60
	TripCountFor1stThreshold    = 3
	TripCheckWindowInSecond     = (TripCountFor1stThreshold-1)*CBCheckIntervalInSecond + 11 // 11 just for tolerance
	SleepCheckCount             = 3
	SleepCheckWindowInSecond    = (SleepCheckCount + 1) * TripCountFor1stThreshold * CBCheckIntervalInSecond
	SleepTimeInSecond           = 20 * 60
	CBMaxSizeInBytes            = 1024
)

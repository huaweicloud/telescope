package utils

const (
	// NameSpace const namespace
	NameSpace = "AGT.ECS"

	// ExternalServiceBMS const namespace of BMS
	ExternalServiceBMS = "SERVICE.BMS"

	// TagBMS const prefix to distinguish BMS from meta_data
	TagBMS = "physical"

	// TagServiceBMS const field of BMS from meta_data
	TagServiceBMS = "metering.resourcespeccode"

	// TTLOneHour const TTL of one hour
	TTLOneHour = 3600

	// TTLTwoDay const TTL of two days
	TTLTwoDay = 172800 // 2 day

	// DimensionName const dimension name
	DimensionName = "instance_id"

	// PostAggregatedMetricDataURI const URI for post aggregated metric data (1 min)
	PostAggregatedMetricDataURI = "/metric-data"

	// PostRawMetricDataURI const URI for post raw metric data (10s)
	PostRawMetricDataURI = "/detailed-metric-data"

	// PostProcessInfo const URI for post process info
	PostProcessInfo = "/process-info"

	// PostCustomMonitorMetricDataURI ...
	PostCustomMonitorMetricDataURI = PostAggregatedMetricDataURI

	// PostEventDataURI
	PostEventDataURI = "/events"

	// Service const for CES agent service
	Service = "CES"

	// PluginConf const for agent plugin config
	PluginConf = "./plugins/conf.json"

	// MaxPluginNum const for max plugin
	MaxPluginNum = 2

	// DefaultPluginCronTime const for default plugin cron time, seconds
	DefaultPluginCronTime = 60
	// DefaultCustomMonitorPluginCronTime ...
	DefaultCustomMonitorPluginCronTime = 60
	// DefaultEventPluginCronTime ...
	DefaultEventPluginCronTime = 60

	// MaxCmdlineLen const for max process cmdline length
	MaxCmdlineLen = 4096

	// CmdlineSuffix const for process cmdline suffix
	CmdlineSuffix = "..."

	// VolumePrefix ...
	// CmdlineSuffix const for process cmdline suffix
	VolumePrefix = "volumeSlAsH"

	// DefaultPluginType AgtPluginType  CustomMonitorPluginType EventPluginType Plugin types
	DefaultPluginType       = ""
	AgtPluginType           = "Agent"
	CustomMonitorPluginType = "Custom Monitor"
	EventPluginType         = "Event"

	// DefaultMaxTimeoutProcNum is the max plugin process number while waiting for output
	DefaultMaxTimeoutProcNum = 5

	// EnvInstanceID ...
	EnvInstanceID = "CES_EVN_INSTANCE_ID"

	ExecutionPerMinute = "0 * * * * *"

	// SendTotal ...
	//metric collection frequency
	SendTotal = 6

	Top5ProcessCollectPeriodInSeconds = 5 * 60
	Top5ProcessSamplePeriodInSeconds  = 60
)

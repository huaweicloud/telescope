package utils

// NameSpace const namespace
const NameSpace = "AGT.ECS"

// ExternalServiceNS const namespace of BMS
const ExternalServiceBMS = "SERVICE.BMS"

// TagBMS const prefix to distinguish BMS from meta_data
const TagBMS = "physical"

// TagServiceBMS const field of BMS from meta_data
const TagServiceBMS = "metering.resourcespeccode"

// TTLOneHour const TTL of one hour
const TTLOneHour = 3600

// TTLTwoDay const TTL of two days
const TTLTwoDay = 172800 // 2 day

// DimensionName const dimension name
const DimensionName = "instance_id"

// PostAggregatedMetricDataURI const URI for post aggregated metric data (1 min)
const PostAggregatedMetricDataURI = "/metric-data"

// PostRawMetricDataURI const URI for post raw metric data (10s)
const PostRawMetricDataURI = "/detailed-metric-data"

// PostProcessInfo const URI for post process info
const PostProcessInfo = "/process-info"

// Service const for CES agent service
const Service = "CES"

// PluginConf const for agent plugin config
const PluginConf = "./plugins/conf.json"

// MaxPluginNum const for max plugin
const MaxPluginNum = 2

// DefaultPluginCronTime const for default plugin cron time, seconds
const DefaultPluginCronTime = 10

// MaxCmdlineLen const for max process cmdline length
const MaxCmdlineLen = 4096

// CmdlineSuffix const for process cmdline suffix
const CmdlineSuffix = "..."
package model

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// GBConversion the multiple of GB --> Byte
const GBConversion = 1024 * 1024 * 1024

// Metric the type for metric data
type Metric struct {
	MetricName   string  `json:"metric_name"`
	MetricValue  float64 `json:"metric_value"`
	MetricPrefix string  `json:"metric_prefix,omitempty"`
}

// InputMetric the type for input metric
type InputMetric struct {
	CollectTime int64    `json:"collect_time"`
	Data        []Metric `json:"data"`
}

// InputMetricSlice the type for input metric sclice
type InputMetricSlice []*InputMetric

// DimensionType the type for dimension
type DimensionType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CesMeticExtraInfo struct {
	OriginMetricName string `json:"origin_metric_name"`
	MetricPrefix     string `json:"metric_prefix,omitempty"`
}

// MetricType the type for metric
type MetricType struct {
	Namespace       string             `json:"namespace"`
	Dimensions      []DimensionType    `json:"dimensions"`
	MetricName      string             `json:"metric_name"`
	MetricExtraInfo *CesMeticExtraInfo `json:"extra_info,omitempty"`
}

// CesMetricData the type for post metric data
type CesMetricData struct {
	Metric      MetricType `json:"metric"`
	TTL         int        `json:"ttl"`
	CollectTime int64      `json:"collect_time"`
	Value       float64    `json:"value"`
	Unit        string     `json:"unit"`
}

// CesMetricDataArr the type for metric data array
type CesMetricDataArr []CesMetricData

var metricUnitMap = map[string]string{
	"cpu_usage_user":               "%",
	"cpu_usage_system":             "%",
	"cpu_usage_idle":               "%",
	"cpu_usage_other":              "%",
	"cpu_usage_nice":               "%",
	"cpu_usage_iowait":             "%",
	"cpu_usage_irq":                "%",
	"cpu_usage_softirq":            "%",
	"cpu_usage_steal":              "%",
	"cpu_usage_guest":              "%",
	"cpu_usage_guest_nice":         "%",
	"mem_total":                    "GB",
	"mem_available":                "GB",
	"mem_used":                     "GB",
	"mem_free":                     "GB",
	"mem_usedPercent":              "%",
	"mem_buffers":                  "GB",
	"mem_cached":                   "GB",
	"net_bitSent":                  "bits/s",
	"net_bitRecv":                  "bits/s",
	"net_packetSent":               "Counts/s",
	"net_packetRecv":               "Counts/s",
	"net_errin":                    "%",
	"net_errout":                   "%",
	"net_dropin":                   "%",
	"net_dropout":                  "%",
	"net_fifoin":                   "Bytes",
	"net_fifoout":                  "Bytes",
	"disk_total":                   "GB",
	"disk_free":                    "GB",
	"disk_used":                    "GB",
	"disk_usedPercent":             "%",
	"disk_inodesTotal":             "",
	"disk_inodesUsed":              "",
	"disk_inodesFree":              "",
	"disk_inodesUsedPercent":       "%",
	"disk_writeBytes":              "Bytes",
	"disk_readBytes":               "Bytes",
	"disk_iopsInProgress":          "Bytes",
	"disk_agt_read_bytes_rate":     "Byte/s",
	"disk_agt_read_requests_rate":  "Requests/Second",
	"disk_agt_write_bytes_rate":    "Byte/s",
	"disk_agt_write_requests_rate": "Requests/Second",
	"disk_writeTime":               "ms/Count",
	"disk_readTime":                "ms/Count",
	"disk_ioUtils":                 "%",
	"proc_cpu":                     "%",
	"proc_mem":                     "%",
	"proc_file":                    "Count",
	"gpu_performance_state":        "",
	"gpu_usage_gpu":                "%",
	"gpu_usage_mem":                "%",
}

// BuildMetric build metric as input metric
func BuildMetric(collectTime int64, data []Metric) *InputMetric {
	return &InputMetric{
		CollectTime: collectTime,
		Data:        data,
	}
}

// BuildCesMetricData build ces metric data
func BuildCesMetricData(inputMetric *InputMetric, isAggregated bool) CesMetricDataArr {
	var dimension DimensionType
	var metricTTL int
	var cesMetricDataArr CesMetricDataArr

	dimension.Name = ces_utils.DimensionName
	dimension.Value = utils.GetConfig().InstanceId
	dimensions := make([]DimensionType, 1)
	dimensions[0] = dimension
	collectTime := inputMetric.CollectTime
	namespace := ces_utils.NameSpace

	externalNamespace := utils.GetConfig().ExternalService
	if externalNamespace == ces_utils.ExternalServiceBMS {
		namespace = externalNamespace
	}

	if isAggregated {
		metricTTL = ces_utils.TTLTwoDay
	} else {
		metricTTL = ces_utils.TTLOneHour
	}

	for _, metric := range inputMetric.Data {

		var newMetricData CesMetricData
		newMetricData.Metric.Dimensions = dimensions

		newMetricData.Metric.MetricName = metric.MetricName
		//metric name has two info ; use hashid and {metricname,MetricPrefix} replace it
		if metric.MetricPrefix != "" {
			newMetricData.Metric.MetricName = generateHashID(metric.MetricPrefix + metric.MetricName)
			newMetricData.Metric.MetricExtraInfo = &CesMeticExtraInfo{OriginMetricName: metric.MetricName, MetricPrefix: metric.MetricPrefix}
		} else {
			// almost for metric name is too long, no scene now, if metric get in the follow logic, it's a new metric
			aliasName := AliasMetricName(metric.MetricName)
			if aliasName != "" {
				newMetricData.Metric.MetricName = aliasName
				newMetricData.Metric.MetricExtraInfo = &CesMeticExtraInfo{OriginMetricName: metric.MetricName}
			}
		}

		newMetricData.Metric.Namespace = namespace
		newMetricData.CollectTime = collectTime
		newMetricData.TTL = metricTTL
		newMetricData.Value = utils.Limit2Decimal(metric.MetricValue)
		newMetricData.Unit = getUnitByMetric(metric.MetricName)

		cesMetricDataArr = append(cesMetricDataArr, newMetricData)

		cesMetricDataArr = setOldMetricData(cesMetricDataArr, newMetricData, metric)
	}

	return cesMetricDataArr

}

func getOldMetricName(metricName, MetricPrefix string) string {

	if MetricPrefix == "" {
		return metricName
	}

	//disk metric
	if strings.HasPrefix(metricName, "disk_") {
		diskPrefix := GetMountPrefix(MetricPrefix)
		return diskPrefix + metricName
	}

	//proc metric
	if strings.HasPrefix(metricName, "proc_") {
		metricSuffix := strings.Split(metricName, "proc")[1]
		return "proc_" + MetricPrefix + metricSuffix
	}

	//gpu metric
	if strings.HasPrefix(metricName, "gpu_") {
		return "slot" + MetricPrefix + "_" + metricName
	}

	//raid metric
	if strings.HasSuffix(metricName, "_device") && strings.HasPrefix(MetricPrefix, "md") {
		return MetricPrefix + "_" + metricName
	}

	return ""
}

//if needed set an old metric data for transition
func setOldMetricData(cesMetricDataArr CesMetricDataArr, originMetricData CesMetricData, metric Metric) CesMetricDataArr {
	// disk:slAsH
	// gpu:slot
	// raid:md
	// proc
	// 以上为 oldMetricName，则需要按照原来格式发送
	// 发送时如果有拼接的出现特殊字符或者超长的   则发送id的指标

	oldMetricName := getOldMetricName(metric.MetricName, metric.MetricPrefix)

	// in old strategy ,the metric has an the other metric name, now resume it and send it
	if oldMetricName != "" && oldMetricName != metric.MetricName {

		oldMetricData := originMetricData
		//init for old metric data
		oldMetricData.Metric.MetricExtraInfo = nil
		oldMetricData.Metric.MetricName = oldMetricName

		aliasName := AliasMetricName(oldMetricData.Metric.MetricName)
		//like disk metric name, need to add extro_info
		if aliasName != "" {
			oldMetricData.Metric.MetricName = aliasName
			oldMetricData.Metric.MetricExtraInfo = &CesMeticExtraInfo{OriginMetricName: oldMetricName}
		}

		cesMetricDataArr = append(cesMetricDataArr, oldMetricData)
	}

	return cesMetricDataArr
}

func AliasMetricName(metricName string) string {

	//more char: . and - and  ~ and /
	pattern := "^([0-9A-Za-z]|_|/)*(-|~|\\.|/){1,}([0-9A-Za-z]|_|/|-|~|\\.)*$"
	match, _ := regexp.MatchString(pattern, metricName)

	if match {
		return generateHashID(metricName)
	}

	pattern = "^([a-z]|[A-Z]){1}([a-z]|[A-Z]|[0-9]|_)*$"
	match, _ = regexp.MatchString(pattern, metricName)
	if match && len(metricName) > 64 {
		return generateHashID(metricName)
	}

	return ""
}

func generateHashID(name string) string {
	keyStr := []byte(name)
	return "id_" + fmt.Sprintf("%x", md5.Sum(keyStr))
}

func getUnitByMetric(metricName string) string {

	return metricUnitMap[metricName]

}

func GetMountPrefix(name string) string {
	slashFlag := "SlAsH"
	// for linux
	diskPrefix := strings.Replace(name, "/", slashFlag, -1)
	// for windows
	diskPrefix = strings.Replace(diskPrefix, ":", "", -1)
	// "_" used to seperate metricName from mountPoint
	diskPrefix = diskPrefix + "_"

	return diskPrefix
}

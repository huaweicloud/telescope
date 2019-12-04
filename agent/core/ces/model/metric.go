package model

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// GBConversion the multiple of GB --> Byte
const GBConversion = 1024 * 1024 * 1024

// Metric the type for metric data
type Metric struct {
	MetricName     string  `json:"metric_name"`
	MetricValue    float64 `json:"metric_value"`
	MetricPrefix   string  `json:"metric_prefix,omitempty"`
	CustomProcName string  `json:"custom_proc_name,omitempty"`
}

// InputMetric the type for input metric
type InputMetric struct {
	CollectTime int64    `json:"collect_time"`
	Type        string   `json:"-"`
	Data        []Metric `json:"data"`
}

// InputMetricSlice the type for input metric slice
type InputMetricSlice []*InputMetric

// DimensionType the type for dimension
type DimensionType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CesMetricExtraInfo ...
type CesMetricExtraInfo struct {
	OriginMetricName string `json:"origin_metric_name"`
	MetricPrefix     string `json:"metric_prefix,omitempty"`
	MetricType       string `json:"metric_type,omitempty"`
	CustomProcName   string `json:"custom_proc_name,omitempty"`
}

// MetricType the type for metric
type MetricType struct {
	Namespace       string              `json:"namespace"`
	Dimensions      []DimensionType     `json:"dimensions"`
	MetricName      string              `json:"metric_name"`
	MetricExtraInfo *CesMetricExtraInfo `json:"extra_info,omitempty"`
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

// CesEventData  ...
type CesEventData struct {
	EventName   string         `json:"event_name"`
	EventSource string         `json:"event_source,omitempty"`
	Time        int64          `json:"time"`
	Detail      CesEventDetail `json:"detail"`
}

// CesEventDetail ...
type CesEventDetail struct {
	Content      string `json:"content,omitempty"`
	GroupID      string `json:"group_id,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	ResourceName string `json:"resource_name"`
	EventState   string `json:"event_state"`
	EventLevel   string `json:"event_level"`
	EventUser    string `json:"event_user"`
	EventType    string `json:"event_type,omitempty"`
}

// CesEventDataArr ...
type CesEventDataArr []CesEventData

var metricUnitMap = map[string]string{
	//cpu
	"cpu_usage":                      "%",//linux  windows
	"cpu_usage_user":                 "%",//linux  windows
	"cpu_usage_system":               "%",//linux  windows
	"cpu_usage_idle":                 "%",//linux  windows
	"cpu_usage_other":                "%",//linux  windows
	"cpu_usage_nice":                 "%",//linux
	"cpu_usage_iowait":               "%",//linux
	"cpu_usage_irq":                  "%",//linux
	"cpu_usage_softirq":              "%",//linux
	//	"cpu_usage_steal":                "%",
	//	"cpu_usage_guest":                "%",
	//	"cpu_usage_guest_nice":           "%",

	//mem
	"mem_available":                  "GB",//linux   windows
	"mem_usedPercent":                "%", //linux   windows
	"mem_free":                       "GB",//linux
	"mem_buffers":                    "GB",//linux
	"mem_cached":                     "GB",//linux
	//	"mem_total":                      "GB",
	//	"mem_used":                       "GB",

	//net
	"net_bitSent":                    "bit/s",    //linux   windows
	"net_bitRecv":                    "bit/s",    //linux   windows
	"net_packetSent":                 "Count/s",  //linux   windows
	"net_packetRecv":                 "Count/s",  //linux   windows
	"net_errin":                      "%",        //linux   windows
	"net_errout":                     "%",        //linux   windows
	"net_dropin":                     "%",        //linux   windows
	"net_dropout":                    "%",        //linux   windows
	//	"net_fifoin":                     "Bytes",
	//	"net_fifoout":                    "Bytes",

	//disk
	"disk_total":                     "GB",       //linux   windows
	"disk_free":                      "GB",       //linux   windows
	"disk_used":                      "GB",       //linux   windows
	"disk_usedPercent":               "%",        //linux   windows
	"disk_agt_read_bytes_rate":       "Byte/s",   //linux   windows
	"disk_agt_read_requests_rate":    "Request/s",//linux   windows
	"disk_agt_write_bytes_rate":      "Byte/s",   //linux   windows
	"disk_agt_write_requests_rate":   "Request/s",//linux   windows
	"disk_inodesTotal":               "",         //linux
	"disk_inodesUsed":                "",         //linux
	"disk_inodesUsedPercent":         "%",        //linux
	"disk_writeTime":                 "ms/Count", //linux
	"disk_readTime":                  "ms/Count", //linux
	"disk_ioUtils":                   "%",        //linux
	"disk_fs_rwstate":                "",         //linux
	"disk_queue_length":              "Count",    //linux
	"disk_write_bytes_per_operation": "KB/op",    //linux
	"disk_read_bytes_per_operation":  "KB/op",    //linux
	"disk_io_svctm":                  "ms/op",    //linux
	"disk_writeBytes":                "Bytes",
	"disk_readBytes":                 "Bytes",
	//	"disk_inodesFree":                "",
	//	"disk_iopsInProgress":            "Bytes",

	//proc
	"proc_cpu":                       "%",    //linux   windows
	"proc_mem":                       "%",    //linux   windows
	"proc_specified_count":           "Count",//linux   windows
	"proc_total_count":               "Count",//linux   windows
	"proc_file":                      "Count",//linux
	//	"proc_count_spicified":           "Count",

	//gpu
	"gpu_performance_state":          "",
	"gpu_usage_gpu":                  "%",
	"gpu_usage_mem":                  "%",

}

// BuildMetric build metric as input metricpro
func BuildMetric(collectTime int64, data []Metric) *InputMetric {
	return &InputMetric{
		CollectTime: collectTime,
		Data:        data,
	}
}

// BuildCesMetricData build ces metric data
func BuildCesMetricData(inputMetric *InputMetric, isAggregated bool) CesMetricDataArr {
	var (
		dimension        DimensionType
		metricTTL        int
		cesMetricDataArr CesMetricDataArr
	)

	dimension.Name = cesUtils.DimensionName
	dimension.Value = utils.GetConfig().InstanceId
	dimensions := make([]DimensionType, 1)
	dimensions[0] = dimension
	collectTime := inputMetric.CollectTime
	namespace := cesUtils.NameSpace

	externalNamespace := utils.GetConfig().ExternalService
	if utils.GetConfig().BmsFlag || externalNamespace == cesUtils.ExternalServiceBMS {
		namespace = cesUtils.ExternalServiceBMS
	}

	if isAggregated {
		metricTTL = cesUtils.TTLTwoDay
	} else {
		metricTTL = cesUtils.TTLOneHour
	}

	for _, metric := range inputMetric.Data {

		var newMetricData CesMetricData
		newMetricData.Metric.Dimensions = dimensions

		newMetricData.Metric.MetricName = metric.MetricName
		// metric name has two info ; use hashid and {metricname,MetricPrefix} replace it
		if metric.MetricPrefix != "" && strings.HasPrefix(metric.MetricPrefix, cesUtils.VolumePrefix) {
			newMetricData.Metric.MetricName = generateHashID(metric.MetricPrefix + metric.MetricName)
			newMetricData.Metric.MetricExtraInfo = &CesMetricExtraInfo{OriginMetricName: metric.MetricName, MetricPrefix: strings.Replace(metric.MetricPrefix, cesUtils.VolumePrefix, "", -1), MetricType: "volume"}
		} else if metric.MetricPrefix != "" {
			newMetricData.Metric.MetricName = generateHashID(metric.MetricPrefix + metric.MetricName)
			newMetricData.Metric.MetricExtraInfo = &CesMetricExtraInfo{OriginMetricName: metric.MetricName, MetricPrefix: metric.MetricPrefix, CustomProcName: metric.CustomProcName}
		} else {
			// almost for metric name is too long, no scene now, if metric get in the follow logic, it's a new metric
			aliasName := AliasMetricName(metric.MetricName)
			if aliasName != "" {
				newMetricData.Metric.MetricName = aliasName
				newMetricData.Metric.MetricExtraInfo = &CesMetricExtraInfo{OriginMetricName: metric.MetricName}
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
	switch {
	case MetricPrefix == "":
		return metricName
	case strings.HasPrefix(metricName, "disk_"):
		diskPrefix := GetMountPrefix(MetricPrefix)
		return diskPrefix + metricName
	case strings.HasPrefix(metricName, "proc_"):
		metricSuffix := strings.Split(metricName, "proc")[1]
		return "proc_" + MetricPrefix + metricSuffix
	case strings.HasPrefix(metricName, "gpu_"):
		return "slot" + MetricPrefix + "_" + metricName
	case strings.HasPrefix(MetricPrefix, "md") && strings.HasSuffix(metricName, "_device"):
		return MetricPrefix + "_" + metricName
	default:
		return ""
	}
}

// if needed set an old metric data for transition
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
		// init for old metric data
		oldMetricData.Metric.MetricExtraInfo = nil
		oldMetricData.Metric.MetricName = oldMetricName

		aliasName := AliasMetricName(oldMetricData.Metric.MetricName)
		// like disk metric name, need to add extro_info
		if aliasName != "" {
			oldMetricData.Metric.MetricName = aliasName
			oldMetricData.Metric.MetricExtraInfo = &CesMetricExtraInfo{OriginMetricName: oldMetricName}
		}

		cesMetricDataArr = append(cesMetricDataArr, oldMetricData)
	}

	return cesMetricDataArr
}

// AliasMetricName ...
func AliasMetricName(metricName string) string {

	// more char: . and - and  ~ and /
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

// GetMountPrefix ...
func GetMountPrefix(name string) string {
	slashFlag := "SlAsH"
	// for linux
	diskPrefix := strings.Replace(name, utils.SLASH, slashFlag, -1)
	// for windows
	diskPrefix = strings.Replace(diskPrefix, ":", "", -1)
	// "_" used to separate metricName from mountPoint
	diskPrefix = diskPrefix + "_"

	return diskPrefix
}

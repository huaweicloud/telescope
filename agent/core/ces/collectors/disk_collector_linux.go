package collectors

import (
	"errors"
	"strings"
	"sync"

	cesCommon "github.com/huaweicloud/telescope/agent/core/ces/gopsutil/common"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/disk"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct {
	DiskMap sync.Map
}

// DiskIOCountersStat is the type for store disk IO data
type DiskIOCountersStat struct {
	collectTime     int64
	uptimeInSeconds int64
	readBytes       float64
	readCount       float64
	writeBytes      float64
	writeCount      float64
	readTime        float64
	writeTime       float64
	ioTime          float64
	weightedIO      float64
}

// Collect implement the disk Collector
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {
	var (
		result  model.InputMetric
		fieldsG []model.Metric
	)

	// /proc/self/mounts & /proc/filesystems
	partitions, err := disk.Partitions(false)
	if err != nil {
		logs.GetCesLogger().Errorf("Exec disk.Partitions(false) failed and error is: %v", err)
		return &result
	}

	// /proc/diskstats
	diskStats, err := disk.IOCounters()
	if err != nil {
		logs.GetCesLogger().Warnf("Exec disk.IOCounters() failed and error is: %v", err)
	}

	devMap := cesCommon.GetDeviceMap()
	mountMapRecord := make(map[string]int)
	for _, p := range partitions {
		// type of mount (ro or rw)
		metric, err := getRWStateMetric(p)
		if err == nil {
			fieldsG = append(fieldsG, metric)
		}

		// discard the duplicate mount point
		diskMountPoint := p.Mountpoint
		if mountMapRecord[diskMountPoint] == 1 {
			continue
		}
		mountMapRecord[diskMountPoint] = 1

		// file system usage
		fsUsageMetrics := getFSUsageMetrics(diskMountPoint)
		fieldsG = append(fieldsG, fsUsageMetrics...)

		// get the proper device name by mount point
		deviceName := getDeviceNameInDiskStats(p, devMap, diskStats)
		if deviceName == "" {
			continue
		}

		// get current stats for the device name above
		currStats := getStats(collectTime, deviceName, diskStats)

		// calculate disk metrics
		if lastStatsData, ok := d.DiskMap.Load(deviceName); ok {
			if lastStats, ok := lastStatsData.(*DiskIOCountersStat); ok {
				diskIOMetrics := getDiskIOMetrics(diskMountPoint, currStats, lastStats)
				logs.GetCesLogger().Debugf("Get disk IO metrics finished, metrics are: %v", diskIOMetrics)
				fieldsG = append(fieldsG, diskIOMetrics...)
			} else {
				logs.GetCesLogger().Errorf("Disk stats found in map for device(%s), but convert failed")
			}
		} else {
			logs.GetCesLogger().Warnf("Disk stats NOT found in map for device(%s)", deviceName)
		}

		d.DiskMap.Store(deviceName, currStats)

	}

	// volume metric collector
	volumeMetrics := getVolumeMetrics(diskStats, d, collectTime)
	fieldsG = append(fieldsG, volumeMetrics...)

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

// getVolumeMetrics returns metrics for device which labeled as 'disk'
// using lsblk to get types(disk/part/lvm)
func getVolumeMetrics(diskStats map[string]disk.IOCountersStat, dc *DiskCollector, collectTime int64) []model.Metric {
	var (
		metrics      []model.Metric
		diskTypeList = cesCommon.GetDeviceTypeMap()["disk"]
	)

	for _, deviceName := range diskTypeList {
		currStats := getStats(collectTime, deviceName, diskStats)
		deviceName4CES := cesUtils.VolumePrefix + deviceName
		if lastStatesData, ok := dc.DiskMap.Load(deviceName4CES); ok {
			if lastStats, ok := lastStatesData.(*DiskIOCountersStat); ok {
				// 老的数据格式卷指标以"volumeSlAsH"开头,即volumeSlAsHxvda_disk....
				// 新的数据格式 metric_prefix:xvda
				diskIOMetrics := getDiskIOMetrics(deviceName4CES, currStats, lastStats)
				logs.GetCesLogger().Debugf("Get disk IO metrics finished, metrics are: %v", diskIOMetrics)
				metrics = append(metrics, diskIOMetrics...)
			} else {
				logs.GetCesLogger().Errorf("Disk stats found in map for device(%s), but convert failed")
			}
		} else {
			logs.GetCesLogger().Warnf("Disk stats NOT found in map for device(%s)", deviceName)
		}

		dc.DiskMap.Store(deviceName4CES, currStats)
	}

	return metrics
}

// getRWStateMetric returns metric for fs readable and writable
func getRWStateMetric(stat disk.PartitionStat) (model.Metric, error) {
	var (
		state float64 = -1
		err   error
	)
	// http://man7.org/linux/man-pages/man5/fstab.5.html
	optsArray := strings.Split(stat.Opts, ",")

	switch {
	case utils.StrArrContainsStr(optsArray, "ro"):
		state = 1
		logs.GetCesLogger().Debugf("Type of mount is ro for partition %v", stat)
	case utils.StrArrContainsStr(optsArray, "rw"):
		state = 0
		logs.GetCesLogger().Debugf("Type of mount is rw for partition %v", stat)
	default:
		err = errors.New("type of mount is not ro neither rw")
		logs.GetCesLogger().Warnf("Type of mount is not ro neither rw for partition %v", stat)
	}

	return model.Metric{
		MetricName:   "disk_fs_rwstate",
		MetricValue:  state,
		MetricPrefix: stat.Mountpoint,
	}, err
}

// getDeltaTime returns the delta time;
// priority: uptime > collect time > default
func getDeltaTime(current, last *DiskIOCountersStat) float64 {
	var deltaTime = float64(utils.DefaultMetricDeltaDataTimeInSecond)
	deltaTimeUsingCT := float64(current.collectTime-last.collectTime) / 1000
	if current.uptimeInSeconds != -1 && last.uptimeInSeconds != -1 {
		deltaTime = float64(current.uptimeInSeconds - last.uptimeInSeconds)
	} else if deltaTimeUsingCT > 0 {
		deltaTime = deltaTimeUsingCT
	}

	return deltaTime
}

// getFSUsageMetrics returns the metric for fs usage
func getFSUsageMetrics(diskMountPoint string) []model.Metric {
	var metrics []model.Metric

	usageStats, err := disk.Usage(diskMountPoint)
	if err != nil {
		logs.GetCesLogger().Errorf("Exec disk.Usage(diskMountPoint) failed and error is: %v", err)
		return metrics
	}

	metrics = append(metrics, model.Metric{
		MetricName:   "disk_total",
		MetricValue:  float64(usageStats.Total) / model.GBConversion,
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_free",
		MetricValue:  float64(usageStats.Free) / model.GBConversion,
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_used",
		MetricValue:  float64(usageStats.Used) / model.GBConversion,
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_usedPercent",
		MetricValue:  float64(usageStats.UsedPercent),
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_inodesTotal",
		MetricValue:  float64(usageStats.InodesTotal),
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_inodesUsed",
		MetricValue:  float64(usageStats.InodesUsed),
		MetricPrefix: diskMountPoint,
	}, model.Metric{
		MetricName:   "disk_inodesUsedPercent",
		MetricValue:  float64(usageStats.InodesUsedPercent),
		MetricPrefix: diskMountPoint,
	})

	return metrics
}

// getDeviceNameInDiskStats returns the device name in /proc/diskstats
// by mount point in /proc/self/mounts using the relationship
// lsblk indicates
func getDeviceNameInDiskStats(p disk.PartitionStat, devMap map[string]string, diskStats map[string]disk.IOCountersStat) string {
	// block device
	deviceName := strings.TrimPrefix(p.Device, "/dev/")
	if _, ok := diskStats[deviceName]; ok {
		logs.GetCesLogger().Debugf("Device name is: %s", deviceName)
		return deviceName
	}

	// lvm device filtering by label
	for _, v := range diskStats {
		label := strings.TrimPrefix(p.Device, "/dev/mapper/")
		vLabel := strings.TrimSuffix(v.Label, "\n") // /sys/block/{dm-0}/dm/name contains \n
		if label == vLabel {
			deviceName = v.Name
			logs.GetCesLogger().Debugf("Device name is: %s, when p.Device is: %s and diskstats is: %v", deviceName, p.Device, v)
			return deviceName
		}

		logs.GetCesLogger().Warnf("Device name dose NOT match label, when p.Device is: %s, diskstats is: %v, label(%s)!=vLabel(%s)", p.Device, v, label, vLabel)
	}

	// lvm device filtering by device number
	if _, ok := devMap[p.Mountpoint]; ok {
		logs.GetCesLogger().Debugf("Mount point has been found in devMap(mount point: %s)", p.Mountpoint)
		// TODO complete it if label filter failed
	}

	logs.GetCesLogger().Warnf("Device name does NOT in diskStats, p.Device name is: %s", p.Device)
	return ""
}

// getStats returns DiskIOCountersStat using diskstats for certain device
func getStats(collectTime int64, deviceName string, diskStats map[string]disk.IOCountersStat) *DiskIOCountersStat {
	if _, ok := diskStats[deviceName]; ok {
		uptimeInSeconds, _ := cesUtils.GetUptimeInSeconds()
		return &DiskIOCountersStat{
			collectTime:     collectTime,
			uptimeInSeconds: uptimeInSeconds,
			readBytes:       float64(diskStats[deviceName].ReadBytes),
			readCount:       float64(diskStats[deviceName].ReadCount),
			writeBytes:      float64(diskStats[deviceName].WriteBytes),
			writeCount:      float64(diskStats[deviceName].WriteCount),
			readTime:        float64(diskStats[deviceName].ReadTime),
			writeTime:       float64(diskStats[deviceName].WriteTime),
			ioTime:          float64(diskStats[deviceName].IoTime),
			weightedIO:      float64(diskStats[deviceName].WeightedIO),
		}
	}

	return &DiskIOCountersStat{}
}

// getDiskIOMetrics returns the metrics for certain fs(mount point) OR disk(/dev/sda)
func getDiskIOMetrics(diskPrefix string, c, l *DiskIOCountersStat) []model.Metric {
	var (
		fieldsG         []model.Metric
		deltaReadBytes  = cesUtils.Float64From32Bits(c.readBytes - l.readBytes)
		deltaReadReq    = cesUtils.Float64From32Bits(c.readCount - l.readCount)
		deltaWriteBytes = cesUtils.Float64From32Bits(c.writeBytes - l.writeBytes)
		deltaWriteReq   = cesUtils.Float64From32Bits(c.writeCount - l.writeCount)
		// ms
		deltaIOTime      = cesUtils.Float64From32Bits(c.ioTime - l.ioTime)
		deltaWriteTime   = cesUtils.Float64From32Bits(c.writeTime - l.writeTime)
		deltaReadTime    = cesUtils.Float64From32Bits(c.readTime - l.readTime)
		deltaQueueLength = cesUtils.Float64From32Bits(c.weightedIO - l.weightedIO)
		// second
		deltaTime = getDeltaTime(c, l)
	)
	if deltaTime > 0 {
		fieldsG = append(fieldsG, model.Metric{
			MetricName:   "disk_agt_read_bytes_rate",
			MetricValue:  float64(deltaReadBytes) / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "disk_agt_read_requests_rate",
			MetricValue:  float64(deltaReadReq) / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "disk_agt_write_bytes_rate",
			MetricValue:  float64(deltaWriteBytes) / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "disk_agt_write_requests_rate",
			MetricValue:  float64(deltaWriteReq) / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "disk_ioUtils",
			MetricValue:  100 * deltaIOTime / (deltaTime * 1000),
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "disk_queue_length",
			MetricValue:  deltaQueueLength / deltaTime,
			MetricPrefix: diskPrefix,
		})
	}

	var (
		diskWriteTime           = 0.0
		diskReadTime            = 0.0
		diskWriteBytesPerSecond = 0.0
		diskReadBytesPerSecond  = 0.0
		diskIOSvctm             = 0.0
	)

	if deltaWriteReq != 0 {
		diskWriteTime = deltaWriteTime / deltaWriteReq
		diskWriteBytesPerSecond = deltaWriteBytes / deltaWriteReq
	}
	if deltaReadReq != 0 {
		diskReadTime = deltaReadTime / deltaReadReq
		diskReadBytesPerSecond = deltaReadBytes / deltaReadReq
	}

	deltaWRReq := deltaReadReq + deltaWriteReq
	if deltaWRReq != 0 {
		diskIOSvctm = deltaIOTime / deltaWRReq
	}

	fieldsG = append(fieldsG, model.Metric{
		MetricName:   "disk_writeTime",
		MetricValue:  diskWriteTime,
		MetricPrefix: diskPrefix,
	}, model.Metric{
		MetricName:   "disk_readTime",
		MetricValue:  diskReadTime,
		MetricPrefix: diskPrefix,
	}, model.Metric{
		MetricName:   "disk_write_bytes_per_operation",
		MetricValue:  diskWriteBytesPerSecond,
		MetricPrefix: diskPrefix,
	}, model.Metric{
		MetricName:   "disk_read_bytes_per_operation",
		MetricValue:  diskReadBytesPerSecond,
		MetricPrefix: diskPrefix,
	}, model.Metric{
		MetricName:   "disk_io_svctm",
		MetricValue:  diskIOSvctm,
		MetricPrefix: diskPrefix,
	})

	return fieldsG
}

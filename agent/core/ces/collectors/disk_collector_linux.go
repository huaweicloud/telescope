package collectors

import (
	"strings"
	"sync"

	cesdisk "github.com/huaweicloud/telescope/agent/core/ces/gopsutil/disk"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/disk"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct {
	DiskMap sync.Map
}

// DiskIOCountersStat is the type for store disk IO data
type DiskIOCountersStat struct {
	collectTime int64
	readBytes   float64
	readCount   float64
	writeBytes  float64
	writeCount  float64
	readTime    float64
	writeTime   float64
	ioTime      float64
}

// Collect implement the disk Collector
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	fieldsG := []model.Metric{}

	diskPartitions, _ := disk.Partitions(false)
	diskInfo, _ := disk.IOCounters()

	if fsState, fsStateErr := cesdisk.GetFileSystemStatus(); fsStateErr != nil{
		logs.GetCesLogger().Errorf("Failed to get filesystem state, error is: %v", fsStateErr)
	}else{
		for _, eachDisk := range diskPartitions{
			diskMountpoint := eachDisk.Mountpoint
			if fsState[disMountPoint].State != -1{
				fieldsG = append(fieldsG, model.Metric{MetricName:"disk_fs_rwstate", MetricValue:float64(fsState[diskMountpoint].State), MetricPrefix : diskMountpoint})
			}
		}
	}

	for _, eachDisk := range diskPartitions {

		diskMountpoint := eachDisk.Mountpoint

		diskStats, _ := disk.Usage(diskMountpoint)

		diskName := strings.TrimPrefix(eachDisk.Device, "/dev/")

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_total", MetricValue: float64(diskStats.Total) / model.GBConversion, MetricPrefix: diskMountpoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_free", MetricValue: float64(diskStats.Free) / model.GBConversion, MetricPrefix: diskMountpoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_used", MetricValue: float64(diskStats.Used) / model.GBConversion, MetricPrefix: diskMountpoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_usedPercent", MetricValue: float64(diskStats.UsedPercent), MetricPrefix: diskMountpoint})

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesTotal", MetricValue: float64(diskStats.InodesTotal), MetricPrefix: diskMountpoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesUsed", MetricValue: float64(diskStats.InodesUsed), MetricPrefix: diskMountpoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesUsedPercent", MetricValue: float64(diskStats.InodesUsedPercent), MetricPrefix: diskMountpoint})

		if diskInfo[diskName].Name == "" {
			logs.GetCesLogger().Infof("No IO data for the disk : %v", diskName)
			continue
		}

		nowStatesInfo := new(DiskIOCountersStat)
		nowStatesInfo.collectTime = collectTime
		nowStatesInfo.readBytes = float64(diskInfo[diskName].ReadBytes)
		nowStatesInfo.readCount = float64(diskInfo[diskName].ReadCount)
		nowStatesInfo.writeBytes = float64(diskInfo[diskName].WriteBytes)
		nowStatesInfo.writeCount = float64(diskInfo[diskName].WriteCount)
		nowStatesInfo.ioTime = float64(diskInfo[diskName].IoTime)
		nowStatesInfo.writeTime = float64(diskInfo[diskName].WriteTime)
		nowStatesInfo.readTime = float64(diskInfo[diskName].ReadTime)

		lastStatesData, ok := d.DiskMap.Load(diskName)
		if ok {
			lastStatesInfo, _ := lastStatesData.(*DiskIOCountersStat)

			DeltaReadBytes := nowStatesInfo.readBytes - lastStatesInfo.readBytes
			DeltaReadReq := nowStatesInfo.readCount - lastStatesInfo.readCount
			DeltaWriteBytes := nowStatesInfo.writeBytes - lastStatesInfo.writeBytes
			DeltaWriteReq := nowStatesInfo.writeCount - lastStatesInfo.writeCount
			DeltaIOTime := nowStatesInfo.ioTime - lastStatesInfo.ioTime
			DeltaWriteTime := nowStatesInfo.writeTime - lastStatesInfo.writeTime
			DeltaReadTime := nowStatesInfo.readTime - lastStatesInfo.readTime

			secondDuration := float64(nowStatesInfo.collectTime-lastStatesInfo.collectTime) / 1000
			if secondDuration != 0 {
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_bytes_rate", MetricValue: float64(DeltaReadBytes) / secondDuration, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_requests_rate", MetricValue: float64(DeltaReadReq) / secondDuration, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_bytes_rate", MetricValue: float64(DeltaWriteBytes) / secondDuration, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_requests_rate", MetricValue: float64(DeltaWriteReq) / secondDuration, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_ioUtils", MetricValue: 100 * DeltaIOTime / (secondDuration * 1000), MetricPrefix: diskMountpoint})
			}

			var diskWriteTime float64 = 0.0
			var diskReadTime float64 = 0.0
			if DeltaWriteReq != 0 {
				diskWriteTime = DeltaWriteTime / DeltaWriteReq
			}
			if DeltaReadReq != 0 {
				diskReadTime = DeltaReadTime / DeltaReadReq
			}
			fieldsG = append(fieldsG, model.Metric{MetricName: "disk_writeTime", MetricValue: diskWriteTime, MetricPrefix: diskMountpoint})
			fieldsG = append(fieldsG, model.Metric{MetricName: "disk_readTime", MetricValue: diskReadTime, MetricPrefix: diskMountpoint})

		}

		d.DiskMap.Store(diskName, nowStatesInfo)
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

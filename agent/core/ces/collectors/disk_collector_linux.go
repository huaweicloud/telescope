package collectors

import (
	"strings"
	"sync"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/disk"
	ces_disk "github.com/huaweicloud/telescope/agent/core/ces/gopsutil/disk"
	ces_utils "github.com/huaweicloud/telescope/agent/core/utils"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct {
	DiskMap sync.Map
}

// DiskIOCountersStat is the type for store disk IO data
type DiskIOCountersStat struct {
	collectTime int64
	uptimeInSeconds int64
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

	if fsState, fsStateErr := ces_disk.GetFileSystemStatus(); fsStateErr != nil{
		logs.GetCesLogger().Errorf("Failed to get filesystem state, error is: %v", fsStateErr)
	}else{
		for _, eachDisk := range diskPartitions{
			diskMountPoint := eachDisk.Mountpoint
			if fsState[diskMountPoint].State != -1{
				fieldsG = append(fieldsG, model.Metric{MetricName:"disk_fs_rwstate", MetricValue:float64(fsState[diskMountPoint].State), MetricPrefix : diskMountPoint})
			}
		}
	}

	for _, eachDisk := range diskPartitions {
		var deltaTime = float64(ces_utils.DEFAULT_DELETA_TIME_IN_SECONDS)
		diskMountPoint := eachDisk.Mountpoint
		diskState, _ := disk.Usage(diskMountPoint)
		diskName := strings.TrimPrefix(eachDisk.Device, "/dev/")

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_total", MetricValue: float64(diskStats.Total) / model.GBConversion, MetricPrefix: diskMountPoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_free", MetricValue: float64(diskStats.Free) / model.GBConversion, MetricPrefix: diskMountPoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_used", MetricValue: float64(diskStats.Used) / model.GBConversion, MetricPrefix: diskMountPoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_usedPercent", MetricValue: float64(diskStats.UsedPercent), MetricPrefix: diskMountPoint})

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesTotal", MetricValue: float64(diskStats.InodesTotal), MetricPrefix: diskMountPoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesUsed", MetricValue: float64(diskStats.InodesUsed), MetricPrefix: diskMountPoint})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_inodesUsedPercent", MetricValue: float64(diskStats.InodesUsedPercent), MetricPrefix: diskMountPoint})

		if diskInfo[diskName].Name == "" {
			logs.GetCesLogger().Infof("No IO data for the disk : %v", diskName)
			continue
		}

		currStatesInfo := new(DiskIOCountersStat)
		currStatesInfo.collectTime = collectTime
		currStatesInfo.uptimeInSeconds, _ = ces_utils.GetUptimeInSeconds()
		currStatesInfo.readBytes = float64(diskInfo[diskName].ReadBytes)
		currStatesInfo.readCount = float64(diskInfo[diskName].ReadCount)
		currStatesInfo.writeBytes = float64(diskInfo[diskName].WriteBytes)
		currStatesInfo.writeCount = float64(diskInfo[diskName].WriteCount)
		currStatesInfo.ioTime = float64(diskInfo[diskName].IoTime)
		currStatesInfo.writeTime = float64(diskInfo[diskName].WriteTime)
		currStatesInfo.readTime = float64(diskInfo[diskName].ReadTime)

		lastStatesData, ok := d.DiskMap.Load(diskName)
		if ok {
			lastStatesInfo, _ := lastStatesData.(*DiskIOCountersStat)

			DeltaReadBytes := currStatesInfo.readBytes - lastStatesInfo.readBytes
			DeltaReadReq := currStatesInfo.readCount - lastStatesInfo.readCount
			DeltaWriteBytes := currStatesInfo.writeBytes - lastStatesInfo.writeBytes
			DeltaWriteReq := currStatesInfo.writeCount - lastStatesInfo.writeCount
			DeltaIOTime := currStatesInfo.ioTime - lastStatesInfo.ioTime
			DeltaWriteTime := currStatesInfo.writeTime - lastStatesInfo.writeTime
			DeltaReadTime := currStatesInfo.readTime - lastStatesInfo.readTime

			deltaTimeUsingCT := float64(currStatesInfo.collectTime - lastStatesInfo.collectTime) / 1000
			if currStatesInfo.uptimeInSeconds != -1 && lastStatesInfo.uptimeInSeconds != -1{
				deltaTime = float64(currStatesInfo.uptimeInSeconds - lastStatesInfo.uptimeInSeconds)
			}else if (deltaTimeUsingCT > 0){
				deltaTime = deltaTimeUsingCT
			}

			if deltaTime != 0 {
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_bytes_rate", MetricValue: float64(DeltaReadBytes) / deltaTime, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_requests_rate", MetricValue: float64(DeltaReadReq) / deltaTime, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_bytes_rate", MetricValue: float64(DeltaWriteBytes) / deltaTime, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_requests_rate", MetricValue: float64(DeltaWriteReq) / deltaTime, MetricPrefix: diskMountpoint})
				fieldsG = append(fieldsG, model.Metric{MetricName: "disk_ioUtils", MetricValue: 100 * DeltaIOTime / (deltaTime * 1000), MetricPrefix: diskMountpoint})
			}

			var diskWriteTime float64 = 0.0
			var diskReadTime float64 = 0.0
			if DeltaWriteReq != 0 {
				diskWriteTime = DeltaWriteTime / DeltaWriteReq
			}
			if DeltaReadReq != 0 {
				diskReadTime = DeltaReadTime / DeltaReadReq
			}
			fieldsG = append(fieldsG, model.Metric{MetricName: "disk_writeTime", MetricValue: diskWriteTime, MetricPrefix: diskMountPoint})
			fieldsG = append(fieldsG, model.Metric{MetricName: "disk_readTime", MetricValue: diskReadTime, MetricPrefix: diskMountPoint})

		}

		d.DiskMap.Store(diskName, currStatesInfo)
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}

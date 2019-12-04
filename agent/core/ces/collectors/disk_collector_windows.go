// Copyright (c) 2014, WAKAYAMA Shirou
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
// * Neither the name of the gopsutil authors nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package collectors

import (
	"context"
	"strings"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/gopsutil/process"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/disk"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct{}

type IOCountersStat struct {
	disk.IOCountersStat
	DiskReadsPerSec      uint64 `json:"DiskReadsPerSec"`      // 磁盘读操作速率,MR287
	DiskWritesPerSec     uint64 `json:"DiskWritesPerSec"`     // 磁盘读速率,MR287
	DiskReadBytesPerSec  uint64 `json:"DiskReadBytesPerSec"`  // 磁盘写操作速率,MR287
	DiskWriteBytesPerSec uint64 `json:"DiskWriteBytesPerSec"` // 磁盘写速率,MR287
}

type Win32PerfFormattedData struct {
	Name                    string
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskWriteQueueLength uint64
	AvgDisksecPerRead       uint64
	AvgDisksecPerWrite      uint64
	DiskReadBytesPerSec     uint64 // 磁盘读操作速率
	DiskReadsPerSec         uint64 // 磁盘读速率
	DiskWriteBytesPerSec    uint64 // 磁盘写操作速率
	DiskWritesPerSec        uint64 // 磁盘写速率
}

// Collect implements DiskCollector
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {
	diskPartitions, err := disk.Partitions(false)
	if nil != err {
		logs.GetCesLogger().Errorf("Get disk partitions error: %s", err.Error())
		return nil
	}
	diskInfo, err := ioCounters()
	if nil != err {
		logs.GetCesLogger().Warnf("Get IO counters stats error: %s", err.Error())
	}

	var fieldsG []model.Metric
	for _, p := range diskPartitions {
		mountPoint := p.Mountpoint
		deviceName := p.Device
		fieldsG = append(fieldsG, getDiskUsageMetrics(mountPoint)...)
		fieldsG = append(fieldsG, getDiskIOMetrics(diskInfo, deviceName, mountPoint)...)
	}

	return &model.InputMetric{
		Data:        fieldsG,
		Type:        "disk",
		CollectTime: collectTime,
	}
}

func ioCounters(names ...string) (map[string]IOCountersStat, error) {
	ret := make(map[string]IOCountersStat, 0)
	var diskIOStats []Win32PerfFormattedData

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := process.WMIQueryWithContext(ctx, "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk", &diskIOStats)
	if err != nil {
		return ret, err
	}
	for _, diskStat := range diskIOStats {
		if len(diskStat.Name) > 3 { // skip *_Total or Hard driver
			continue
		}

		if len(names) > 0 && !StringsHas(names, diskStat.Name) {
			continue
		}

		ret[diskStat.Name] = IOCountersStat{
			IOCountersStat: disk.IOCountersStat{
				Name:       diskStat.Name,
				ReadCount:  uint64(diskStat.AvgDiskReadQueueLength),
				WriteCount: diskStat.AvgDiskWriteQueueLength,
				ReadBytes:  uint64(diskStat.AvgDiskBytesPerRead),
				WriteBytes: uint64(diskStat.AvgDiskBytesPerWrite),
				ReadTime:   diskStat.AvgDisksecPerRead,
				WriteTime:  diskStat.AvgDisksecPerWrite,
			},
			// 新增指标：
			// DiskReadsPerSec   磁盘读操作速率
			DiskReadsPerSec: uint64(diskStat.DiskReadsPerSec),
			// DiskWritesPerSec   磁盘写操作速率
			DiskWritesPerSec: uint64(diskStat.DiskWritesPerSec),
			// DiskReadBytesPerSec   磁盘读速率
			DiskReadBytesPerSec: uint64(diskStat.DiskReadBytesPerSec),
			// DiskWriteBytesPerSec   磁盘写速率
			DiskWriteBytesPerSec: uint64(diskStat.DiskWriteBytesPerSec),
		}
	}
	return ret, nil
}

// StringsHas checks the target string slice contains src or not
func StringsHas(target []string, src string) bool {
	for _, t := range target {
		if strings.TrimSpace(t) == src {
			return true
		}
	}
	return false
}

func getDiskUsageMetrics(mountPoint string) []model.Metric {
	diskStats, err := disk.Usage(mountPoint)
	if err != nil {
		logs.GetCesLogger().Errorf("Get disk usage for %s error: %s", mountPoint, err.Error())
		return []model.Metric{}
	}

	return []model.Metric{
		{
			MetricName:   "disk_total",
			MetricValue:  float64(diskStats.Total) / model.GBConversion,
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_free",
			MetricValue:  float64(diskStats.Free) / model.GBConversion,
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_used",
			MetricValue:  float64(diskStats.Used) / model.GBConversion,
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_usedPercent",
			MetricValue:  float64(diskStats.UsedPercent),
			MetricPrefix: mountPoint,
		},
	}
}

func getDiskIOMetrics(diskInfo map[string]IOCountersStat, deviceName, mountPoint string) []model.Metric {
	d, ok := diskInfo[deviceName]
	if !ok || d.Name == "" {
		logs.GetCesLogger().Errorf("No IO data for disk(%s) or disk.IOCountersStat.Name(%s) is empty", deviceName, d.Name)
		return []model.Metric{}
	}

	return []model.Metric{
		{
			MetricName:   "disk_agt_read_bytes_rate",
			MetricValue:  float64(d.DiskReadBytesPerSec),
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_agt_read_requests_rate",
			MetricValue:  float64(d.DiskReadsPerSec),
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_agt_write_bytes_rate",
			MetricValue:  float64(d.DiskWriteBytesPerSec),
			MetricPrefix: mountPoint,
		},
		{
			MetricName:   "disk_agt_write_requests_rate",
			MetricValue:  float64(d.DiskWritesPerSec),
			MetricPrefix: mountPoint,
		},
	}
}

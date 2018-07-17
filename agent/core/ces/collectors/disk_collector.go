// +build !linux,!windows

package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct {
}

// Collect implement the disk Collector
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {
	return nil
}

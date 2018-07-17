package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// CollectorInterface for raw metric collect
type CollectorInterface interface {
	Collect(collectTime int64) *model.InputMetric
}

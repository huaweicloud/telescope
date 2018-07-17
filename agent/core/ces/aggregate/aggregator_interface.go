package aggregate

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// AggregatorInterface for metric aggregate
type AggregatorInterface interface {
	Aggregate(metricSlice model.InputMetricSlice) *model.InputMetric
}

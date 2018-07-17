package aggregate

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// MinValue is the min result type for Aggregate
type MinValue struct {
}

// Aggregate implement the min aggregator
func (minValue *MinValue) Aggregate(input model.InputMetricSlice) *model.InputMetric {

	if input == nil || len(input) == 0 {
		logs.GetCesLogger().Error("Input slice is nil or empty")
		return nil
	}
	minMetric := *input[0]
	metricNameKeyMap := GenerateMetricNameKeyMap(&minMetric.Data)
	for _, metricData := range input {

		for _, metric := range metricData.Data {
			if metric.MetricValue < metricNameKeyMap[metric.MetricPrefix+metric.MetricName].MetricValue {
				metricNameKeyMap[metric.MetricPrefix+metric.MetricName].MetricValue = metric.MetricValue
			}
		}

	}

	return &minMetric

}

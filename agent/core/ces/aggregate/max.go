package aggregate

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// MaxValue is the max result type for Aggregate
type MaxValue struct {
}

// Aggregate implement the max aggregator
func (maxValue *MaxValue) Aggregate(input model.InputMetricSlice) *model.InputMetric {
	if input == nil || len(input) == 0 {
		logs.GetCesLogger().Error("Input slice is nil or empty")
		return nil
	}
	metric := input[0]
	if nil == metric {
		logs.GetCesLogger().Error("Metric in slice is nil")
		return nil
	}
	maxMetric := *metric
	metricNameKeyMap := GenerateMetricNameKeyMap(&maxMetric.Data)
	for _, metricData := range input {
		if nil == metricData || len(metricData.Data) == 0 {
			continue
		}
		for _, metric := range metricData.Data {
			if metric.MetricValue > metricNameKeyMap[metric.MetricPrefix+metric.MetricName].MetricValue {
				metricNameKeyMap[metric.MetricPrefix+metric.MetricName].MetricValue = metric.MetricValue
			}
		}

	}
	return &maxMetric
}

// GenerateMetricNameKeyMap ...
func GenerateMetricNameKeyMap(metrics *[]model.Metric) map[string]*model.Metric {
	metricNameKeyMap := make(map[string]*model.Metric, 0)
	for index := range *metrics {
		metricNameKeyMap[(*metrics)[index].MetricPrefix+(*metrics)[index].MetricName] = &(*metrics)[index]
	}
	return metricNameKeyMap
}

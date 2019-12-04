package aggregate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// AvgValue is the average result type for Aggregate
type AvgValue struct {
}

// Aggregate implement the average aggregator
func (averageValue *AvgValue) Aggregate(input model.InputMetricSlice) *model.InputMetric {
	if input == nil || len(input) == 0 {
		logs.GetCesLogger().Error("Input slice is nil or empty")
		return nil
	}

	dataCount := len(input)
	metric := input[0]
	if nil == metric {
		logs.GetCesLogger().Error("Metric in slice is nil")
		return nil
	}
	avgMetric := *metric
	// aggregate collectTime Round to Minute
	avgMetric.CollectTime = time.Unix(avgMetric.CollectTime/1000, 0).Truncate(time.Minute).Unix() * 1000
	metricNameKeyMap := GenerateMetricNameKeyMap(&avgMetric.Data)

	metricCount := len((*metric).Data)
	sum := make(map[string]float64, metricCount)
	for _, metricData := range input {
		if nil == metricData || len(metricData.Data) == 0 {
			continue
		}
		for _, metric := range metricData.Data {
			sum[metric.MetricPrefix+metric.MetricName] = sum[metric.MetricPrefix+metric.MetricName] + metric.MetricValue
		}
	}

	for _, metric := range avgMetric.Data {
		avg := sum[metric.MetricPrefix+metric.MetricName] / float64(dataCount)
		metricNameKeyMap[metric.MetricPrefix+metric.MetricName].MetricValue, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", avg), 64)
	}
	return &avgMetric
}

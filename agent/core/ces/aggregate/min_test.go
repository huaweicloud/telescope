package aggregate

import (
	"testing"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

func TestAggregateMin(t *testing.T) {
	var testMetricSlice model.InputMetricSlice

	var testMetric1, testMetric2, testMetric3 model.InputMetric

	var minMetric *model.InputMetric

	testMetric1.CollectTime = 1496643900000
	testMetric1.Data = []model.Metric{
		model.Metric{MetricName: "mem_free", MetricValue: 1.1},
		model.Metric{MetricName: "mem_used", MetricValue: 76.4},
	}

	testMetric2.CollectTime = 1496643910000
	testMetric2.Data = []model.Metric{
		model.Metric{MetricName: "mem_free", MetricValue: 1.11},
		model.Metric{MetricName: "mem_used", MetricValue: 80},
	}

	testMetric3.CollectTime = 1496643920000
	testMetric3.Data = []model.Metric{
		model.Metric{MetricName: "mem_free", MetricValue: 1.105},
		model.Metric{MetricName: "mem_used", MetricValue: 70},
	}

	testMetricSlice = append(testMetricSlice, &testMetric1)
	testMetricSlice = append(testMetricSlice, &testMetric2)
	testMetricSlice = append(testMetricSlice, &testMetric3)

	minMetric = new(MinValue).Aggregate(testMetricSlice)
	keyMap := GenerateMetricNameKeyMap(&minMetric.Data)
	if (*minMetric).CollectTime != 1496643900000 || keyMap["mem_used"].MetricValue != 70 || keyMap["mem_free"].MetricValue != 1.1 {
		t.Error("min error")
	}
}

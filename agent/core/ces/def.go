package ces

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// common variables (chans and vars)
var (
	// Channels
	chRawData, chAgResult, chPluginData chan *model.InputMetric
	chAgRawData                         chan model.InputMetricSlice
	chProcessInfo                       chan model.ChProcessList
)

// Initialize the original data channel
func initchRawData() {
	chRawData = make(chan *model.InputMetric, 100)
}

// Get the original data channel
func getchRawData() chan *model.InputMetric {
	if chRawData == nil {
		initchRawData()
	}

	return chRawData
}

// Initialize the original data channel
func initchPluginData() {
	chPluginData = make(chan *model.InputMetric, 100)
}

// Get the original data channel
func getchPluginData() chan *model.InputMetric {
	if chPluginData == nil {
		initchPluginData()
	}

	return chPluginData
}

// Initialize the aggregate data channel
func initchAgRawData() {
	chAgRawData = make(chan model.InputMetricSlice, 100)
}

// Get the data channel
func getchAgRawData() chan model.InputMetricSlice {
	if chAgRawData == nil {
		initchAgRawData()
	}

	return chAgRawData
}

// Initialize the agResult channel
func initchAgResult() {
	chAgResult = make(chan *model.InputMetric, 100)
}

// Get the agResult channel
func getchAgResult() chan *model.InputMetric {
	if chAgResult == nil {
		initchAgResult()
	}

	return chAgResult
}

func initchProcessInfo() {
	chProcessInfo = make(chan model.ChProcessList, 10)
}

func getchProcessInfo() chan model.ChProcessList {
	if chProcessInfo == nil {
		initchProcessInfo()
	}
	return chProcessInfo
}

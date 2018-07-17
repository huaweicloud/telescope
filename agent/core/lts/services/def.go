package services

import "github.com/huaweicloud/telescope/agent/core/lts/logdumper"

var chData chan logdumper.FileEvent

// Initialize the data channel
func InitchData() {
	chData = make(chan logdumper.FileEvent, 20)
}

// Get the data channel
func GetchData() chan logdumper.FileEvent {
	if chData == nil {
		InitchData()
	}

	return chData
}

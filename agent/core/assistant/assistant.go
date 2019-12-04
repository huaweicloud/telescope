package assistant

import (
	"github.com/huaweicloud/telescope/agent/core/assistant/config"
	"github.com/huaweicloud/telescope/agent/core/assistant/heartbeat"
	"github.com/huaweicloud/telescope/agent/core/assistant/task"
)

// Assistant ...
type Assistant struct {
	Switch chan bool
}

// Init ...
func (s *Assistant) Init() {
	config.InitConfig()
	s.Switch = make(chan bool, 1)
}

// Start ...
func (s *Assistant) Start() {
	go heartbeat.SendHBTicker(s.Switch)
	go task.PullTasksTicker(s.Switch)
	go task.ReplyTaskTicker()
}

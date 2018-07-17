package logdumper

import (
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type Event struct {
	Message    string `json:"message" required:"true" max:"256KB"`
	Time       uint64 `json:"time" required:"true"`
	Path       string `json:"path" required:"true"`
	Ip         string `json:"ip" required:"true"`
	HostName   string `json:"host_name" required:"true"`
	LineNumber uint64 `json:"line_no"`
}

type LogEventMessage struct {
	LogEvents  []Event `json:"log_events" required:"true" len:"20"`
	LogTopicId string  `json:"log_topic_id" required:"true"`
	LogGroupId string  `json:"log_group_id" required:"true"`
}

//validate Event
func (e *Event) Validate() error {
	fields := utils.ErrInvalidFields{Object: "Event"}
	if len(e.Message) == 0 {
		fields.Add(utils.NewErrFieldRequired("Message"))
	}
	if len(e.Message) > lts_utils.CONTENT_LENGTH_LIMIT_PER_LOG_TEXT {
		fields.Add(utils.NewErrFieldMaxLen("Message", lts_utils.CONTENT_LENGTH_LIMIT_PER_LOG_TEXT))
	}
	if e.Time <= 0 {
		fields.Add(utils.NewErrFieldRequired("Time"))
	}
	if len(e.HostName) <= 0 {
		fields.Add(utils.NewErrFieldRequired("HostName"))
	}
	if len(e.Ip) <= 0 {
		fields.Add(utils.NewErrFieldRequired("Ip"))
	}
	if len(e.Path) <= 0 {
		fields.Add(utils.NewErrFieldRequired("Path"))
	}
	if fields.Len() > 0 {
		return fields
	}
	return nil
}

//validate log event message
func (e *LogEventMessage) Validate() error {
	fields := utils.ErrInvalidFields{Object: "LogEventMessage"}
	if len(e.LogEvents) <= 0 {
		fields.Add(utils.NewErrFieldRequired("LogEvents"))
	}
	if len(e.LogEvents) > lts_utils.PER_FILE_EVENT_LOGS_MAX_NUMBER {
		fields.Add(utils.NewErrFieldMaxLen("LogEvents", lts_utils.PER_FILE_EVENT_LOGS_MAX_NUMBER))
	}
	if len(e.LogGroupId) <= 0 {
		fields.Add(utils.NewErrFieldRequired("LogGroupId"))
	}
	if len(e.LogTopicId) <= 0 {
		fields.Add(utils.NewErrFieldRequired("LogTopicId"))
	}

	if fields.Len() > 0 {
		return fields
	}
	return nil
}

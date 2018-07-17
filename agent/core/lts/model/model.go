package model

import "github.com/elastic/beats/winlogbeat/sys"

type LtsDetail struct {
	Errors []LtsError `json:"errors"`
}

type LtsError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type WindowsOsLogEventXml struct {
	sys.Event
}

type WindowsSystemEvent struct {
	ProviderName      string `json:"provider_name"`
	EventSourceName   string `json:"event_source_name"`
	EventId           string `json:"event_id"`
	Version           string `json:"version"`
	Level             string `json:"level"`
	Task              string `json:"task"`
	Opcode            string `json:"op_code"`
	TimeCreated       int64  `json:"time_created"`
	RecordId          string `json:"event_record_id"`
	ActivityId        string `json:"activity_id"`
	RelatedActivityID string `json:"related_activity_id"`
	ProcessId         string `json:"process_id"`
	ThreadId          string `json:"thread_id"`
	Channel           string `json:"channel""`
	HostName          string `json:"host_name"`
	UserId            string `json:"user_id"`
}

package logdumper

import (
	"testing"
)

func TestValidateEvent(t *testing.T) {
	event := Event{Message: "logs"}
	err := event.Validate()
	if err != nil {
		t.Log(err.Error())
	}
}

func TestValidateLogEventMessage(t *testing.T) {
	event := LogEventMessage{LogTopicName: "topic1"}
	err := event.Validate()
	if err != nil {
		t.Log(err.Error())
	}
}

package logdumper

import (
	"testing"
)

func TestGetStatesFromRecordFile(t *testing.T) {
	oldStates := GetStatesFromRecordFile()
	if len(oldStates.States) > 0 {
		for oldstateindex := range oldStates.States {
			t.Logf("the state is:%s", oldStates.States[oldstateindex].FilePath)
		}
	} else {
		t.Logf("there is no state entry")
	}

}

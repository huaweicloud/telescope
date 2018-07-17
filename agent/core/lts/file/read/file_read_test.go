package read

import (
	"testing"

	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
)

func TestReadFixedLengthFromOffset(t *testing.T) {
	logs, offset := ReadFixedLengthFromOffset("D:/logs/app.log", uint64(1628), lts_utils.PER_FILE_EVENT_LOGS_MAX_TOTAL_SIZE)
	for index := range logs {
		t.Logf(logs[index])
	}
	t.Logf("the offset is:%d\n", offset)
}

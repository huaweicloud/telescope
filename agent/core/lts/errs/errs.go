package errs

import (
	"encoding/json"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/model"
)

const MAX_ERRORS_COUNT = 100

var ltsErrorChan chan *model.LtsError = make(chan *model.LtsError, MAX_ERRORS_COUNT)

// get lts detail from lts error chan
func GetLtsDetail() string {
	ltsDetail := model.LtsDetail{}
Loop:
	for {
		select {
		case ltsError := <-ltsErrorChan:
			if len(ltsDetail.Errors) < MAX_ERRORS_COUNT {
				ltsDetail.Errors = append(ltsDetail.Errors, *ltsError)
			} else {
				break Loop
			}
		default:
			// if no data in chan, skip this for-select loop
			break Loop
		}
	}

	if len(ltsDetail.Errors) == 0 {
		ltsDetail.Errors = make([]model.LtsError, 0)
	}

	detail, err := json.Marshal(ltsDetail)
	if err != nil {
		logs.GetLtsLogger().Errorf("Marshal lts detail failed:%s.", err.Error())
		return ""
	}
	return string(detail)
}

// put error code and message to lts error chan
func PutLtsDetail(code, message string) {
	if len(ltsErrorChan) < MAX_ERRORS_COUNT {
		ltsError := model.LtsError{code, message}
		ltsErrorChan <- &ltsError
	}
}

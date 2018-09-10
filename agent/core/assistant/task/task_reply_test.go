package task

import (
	"errors"
	assistant_utils "github.com/huaweicloud/telescope/agent/core/assistant/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

//getTasks
func TestGetTasks(t *testing.T) {

	//go getTasks()
	Convey("Test_getTasks", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			// make(map[string]*Task)
			TaskMap["123"] = nil
			getTasks()
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			//getTaskState
			mockfn.Replace((*Task).getTaskState, func(*Task) string {
				return ""
			})
			mockfn.Replace((*Task).getTaskOutput, func(*Task) string {
				return ""
			})
			// make(map[string]*Task)
			TaskMap["123"] = &Task{
				TaskPulled: &TaskPulled{},
			}
			getTasks()
		})
	})

}

//replyTasks
func TestReplyTasks(te *testing.T) {

	Convey("Test_getTasks", te, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			//assistant_utils.GetMarshalledRequestBody
			mockfn.Replace(assistant_utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, errors.New("123")
			})
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})

			body := ReplyTaskRequestBody{}
			replyTasks(&body)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})

			body := ReplyTaskRequestBody{}
			replyTasks(&body)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 2001,
					Body:       ioutil.NopCloser(strings.NewReader("123")),
				}
				return response, errors.New("")
			})
			body := ReplyTaskRequestBody{}
			replyTasks(&body)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 2001,
					Body:       ioutil.NopCloser(strings.NewReader("123")),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), nil
			})
			//json.Unmarshal
			mockfn.Replace(json.Unmarshal, func(data []byte, v interface{}) error {
				return errors.New("123")
			})
			body := ReplyTaskRequestBody{}
			replyTasks(&body)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(assistant_utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("123")),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), nil
			})
			//json.Unmarshal
			mockfn.Replace(json.Unmarshal, func(data []byte, v interface{}) error {
				te.SkipNow()
				return errors.New("123")
			})
			/*	body := &ReplyTaskRequestBody{
					InstanceID: "123",
					Tasks: []ReplyTask{
						{TaskID: ""},
					},
				}
				replyTasks(body)*/
		})

	})
}

//ReplyTaskTicker
func TestReplyTaskTicker(te *testing.T) {
	Convey("Test_ReplyTaskTicker", te, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.NewTicker, func(d time.Duration) *time.Ticker {
				te.SkipNow()
				return nil
			})
			//replyTaskExec
			mockfn.Replace(replyTaskExec, func() {
				te.SkipNow()
			})
			ReplyTaskTicker()
		})
	})
}

//updateTasks
func TestUpdateTasks(te *testing.T) {

	Convey("Test_updateTasks", te, func() {
		Convey("test case 1", func() {
			//isFinalState
			defer mockfn.RevertAll()
			mockfn.Replace(isFinalState, func(cronFlag bool, status string) bool {
				return true
			})
			body := ReplyTaskRespBody{
				SucceededList: []string{"123", ""},
			}
			task := []ReplyTask{
				{
					TaskID:       "",
					InvocationID: "123",
					Status:       "",
					ErrNum:       "",
					Output:       "",
				},
			}
			TaskMap["123"] = &Task{
				TaskPulled: &TaskPulled{},
			}
			updateTasks(&body, task)
		})
		Convey("test case2", func() {
			body := ReplyTaskRespBody{
				SucceededList: []string{"", ""},
			}
			updateTasks(&body, nil)
		})
	})
}

//replyTaskExec
func TestReplyTaskExec(te *testing.T) {
	//replyTaskExec()
	Convey("Test_replyTaskExec", te, func() {
		Convey("test case 1", func() {
			//getTasks
			defer mockfn.RevertAll()
			mockfn.Replace(getTasks, func() *ReplyTaskRequestBody {
				body := &ReplyTaskRequestBody{
					Tasks: []ReplyTask{
						{TaskID: "12"},
					},
				}
				return body
			})
			//replyTasks
			mockfn.Replace(replyTasks, func(requestBody *ReplyTaskRequestBody) (*ReplyTaskRespBody, error) {
				return nil, nil
			})
			//updateTasks
			mockfn.Replace(updateTasks, func(replyTaskRespBody *ReplyTaskRespBody, taskReplySnapshot []ReplyTask) {

			})
			replyTaskExec()
		})
		Convey("test case2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getTasks, func() *ReplyTaskRequestBody {
				body := &ReplyTaskRequestBody{
					Tasks: []ReplyTask{
						{TaskID: "12"},
					},
				}
				return body
			})
			//replyTasks
			mockfn.Replace(replyTasks, func(requestBody *ReplyTaskRequestBody) (*ReplyTaskRespBody, error) {
				return nil, errors.New("123")
			})
			replyTaskExec()
		})
	})
}

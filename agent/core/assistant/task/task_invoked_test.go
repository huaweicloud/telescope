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
)

//PullTasksTicker
func TestPullTasksTicker(te *testing.T) {
	Convey("Test_PullTasksTicker", te, func() {
		Convey("test case 1", func() {
			bools := make(chan bool, 1)
			go PullTasksTicker(bools)
		})
	})
}

//pullTasks
func TestPullTasks(te *testing.T) {

	Convey("Test_pullTasks", te, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return "123"
			})
			//utils.GetConfig
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				return nil, errors.New("123")
			})
			mockfn.Replace((*TaskEventProcessor).Action, func(t *TaskEventProcessor, action string, fromState string, toState string, args []interface{}) {
			})
			/*	mockfn.Replace(io.ReadCloser.Close, func(io.ReadCloser) error {
				return nil
			})*/
			pullTasks()

		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return "123"
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace((*TaskEventProcessor).Action, func(t *TaskEventProcessor, action string, fromState string, toState string, args []interface{}) {
			})
			/*	mockfn.Replace((*TaskEventProcessor).Action, func(t *TaskEventProcessor, action string, fromState string, toState string, args []interface{}) {
				})
				mockfn.Replace(io.ReadCloser.Close, func(io.ReadCloser) error {
					return nil
				})*/
			//go pullTasks()
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), errors.New("123")
			})
			//io.ReadCloser
			/*	mockfn.Replace(io.ReadCloser.Close, func(io.ReadCloser) error {
				return nil
			})*/
			mockfn.Replace((*TaskEventProcessor).Action, func(t *TaskEventProcessor, action string, fromState string, toState string, args []interface{}) {
			})
			//go pullTasks()

		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(assistant_utils.BuildURL, func(destURI string) string {
				return ""
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("1231")),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			/*mockfn.Replace(logs.GetLogger, func() seelog.LoggerInterface {
				//syscall.Exit(0)
				//configFromReader
				loggerInterface, _ := seelog.LoggerFromConfigAsBytes([]byte("123"))
				return loggerInterface
			})
			mockfn.Replace(seelog.LoggerInterface.Debugf, func(s seelog.LoggerInterface, format string, params ...interface{}) {

			})*/
			mockfn.Replace((*TaskEventProcessor).Action, func(t *TaskEventProcessor, action string, fromState string, toState string, args []interface{}) {
			})
			//LoggerInterface
			//go pullTasks()
		})
		Convey("test case 5", func() {
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
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("12123")),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), nil
			})
			//json.Unmarshal
			mockfn.Replace(json.Unmarshal, func(data []byte, v interface{}) error {
				return nil
			})
			//go pullTasks()
		})
	})
}

//pullTasksExec
func TestPullTasksExec(te *testing.T) {
	Convey("Test_pullTasksExec", te, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(pullTasks, func() []TaskPulled {
				pulleds := []TaskPulled{
					{
						InvocationID: "23",
					},
				}
				return pulleds
			})
			//StartTaskExecEntity
			mockfn.Replace(StartTaskExecEntity, func(taskExecEntity *TaskExecEntity) {

			})
			//TaskExecEntity
			mockfn.Replace((*TaskExecEntity).RunAndListen, func(*TaskExecEntity) {
				te.SkipNow()
			})

			pullTasksExec()
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(pullTasks, func() []TaskPulled {
				pulleds := []TaskPulled{
					{
						InvocationID: "23",
					},
				}
				return pulleds
			})
			//StartTaskExecEntity
			mockfn.Replace(StartTaskExecEntity, func(taskExecEntity *TaskExecEntity) {

			})
			//TaskExecEntity
			mockfn.Replace((*TaskExecEntity).RunAndListen, func(*TaskExecEntity) {
				te.SkipNow()
			})
			//isCronTask
			mockfn.Replace((*Task).isCronTask, func(*Task) bool {
				return true
			})
			//pushTaskExec
			mockfn.Replace((*Task).pushTaskExec, func(t *Task, taskExecEntity *TaskExecEntity) {
				return
			})
			pullTasksExec()
		})
	})
}

//pullTasks

//InstantiateTaskExecEntity
func TestInstantiateTaskExecEntity(te *testing.T) {
	Convey("Test_InstantiateTaskExecEntity", te, func() {
		Convey("test case 1", func() {
			//InstantiateTaskInterface
			defer mockfn.RevertAll()
			mockfn.Replace(InstantiateTaskInterface, func(taskPulled *TaskPulled) TaskInterface {
				return nil
			})
			//ExecTaskAction
			mockfn.Replace(ExecTaskAction, func(action string, task *TaskExecEntity) error {
				return nil
			})
			pulled := TaskPulled{}
			entity := InstantiateTaskExecEntity(&pulled)
			So(entity, ShouldNotBeNil)
		})
	})
}

//InstantiateTaskInterface
func TestInstantiateTaskInterface(te *testing.T) {

	Convey("Test_InstantiateTaskInterface", te, func() {
		Convey("test case 1", func() {
			pulled := TaskPulled{
				Command: TASK_COMMAND_RUN_SHELL,
			}
			InstantiateTaskInterface(&pulled)
		})
		Convey("test case 2", func() {
			pulled := TaskPulled{
				Command: TASK_COMMAND_CANCEL,
			}
			InstantiateTaskInterface(&pulled)
		})
		Convey("test case 3", func() {
			pulled := TaskPulled{
				Command: "123",
			}
			InstantiateTaskInterface(&pulled)
		})
	})
}

//StartTaskExecEntity
func TestStartTaskExecEntity(te *testing.T) {

	Convey("Test_StartTaskExecEntity", te, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*TaskExecEntity).RunAndListen, func(*TaskExecEntity) {
				te.SkipNow()
			})
			StartTaskExecEntity(nil)
		})
	})
}

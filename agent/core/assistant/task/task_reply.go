package task

import (
	"bytes"
	"errors"
	assistant_utils "github.com/huaweicloud/telescope/agent/core/assistant/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func getTasks() *ReplyTaskRequestBody {
	replyTaskRequestBody := ReplyTaskRequestBody{
		InstanceID: utils.GetConfig().InstanceId,
		Tasks:      []ReplyTask{},
	}

	for invocationID, task := range TaskMap {
		if task == nil {
			logs.GetAssistantLogger().Warnf("Task is empty.(TaskInvocationID:%s)", invocationID)
			break
		}

		replyTask := ReplyTask{
			TaskID:       task.TaskPulled.TaskID,
			InvocationID: task.TaskPulled.InvocationID,
			// TODO Status和Output增加异常处理(status为failed、succeeded才有output)，ErrNum确定后再说
			Status: task.getTaskState(),
			ErrNum: "",
			Output: task.getTaskOutput(),
		}
		replyTaskRequestBody.Tasks = append(replyTaskRequestBody.Tasks, replyTask)
	}

	return &replyTaskRequestBody
}

func replyTasks(requestBody *ReplyTaskRequestBody) (*ReplyTaskRespBody, error) {
	replyTaskRespBody := ReplyTaskRespBody{}

	url := assistant_utils.BuildURL(assistant_utils.REPLY_TASK_URI)
	requestBodyInBytes, err := assistant_utils.GetMarshalledRequestBody(*requestBody, url)
	if err != nil {
		return &replyTaskRespBody, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyInBytes))
	if err != nil {
		logs.GetAssistantLogger().Errorf("[ReplyTasks]Create request Error: %s", err.Error())
		return &replyTaskRespBody, err
	}

	res, err := utils.HTTPSend(request, assistant_utils.SERVICE)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[ReplyTasks]Failed to request, url is %s, error is %s", url, err.Error())
		return &replyTaskRespBody, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logs.GetAssistantLogger().Errorf("[ReplyTasks]Request server succeeded but status code is %d(should be %d), url is %s, requestBody is %v", res.StatusCode, http.StatusOK, url, requestBody)
		return &replyTaskRespBody, errors.New("Status code does not match")
	}

	resBodyBytes, _ := ioutil.ReadAll(res.Body)
	logs.GetAssistantLogger().Debugf("[ReplyTasks]Response: %s", string(resBodyBytes))
	err = json.Unmarshal(resBodyBytes, &replyTaskRespBody)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[ReplyTasks]Failed to unmarshal response [%s], error is %s.", string(resBodyBytes), err.Error())
		return &replyTaskRespBody, err
	}

	return &replyTaskRespBody, nil
}

// ReplyTaskTicker ..
func ReplyTaskTicker() {
	// TODO 从配置文件或server端获取间隔时间
	ticker := time.NewTicker(time.Second * assistant_utils.REPLY_TASK_INTERVAL_IN_SECOND)
	for t := range ticker.C {
		if len(TaskMap) != 0 {
			logs.GetAssistantLogger().Infof("Reply task starts at %s", t.String())
			go replyTaskExec()
		}
	}
}

// updateTasks deletes the NON-cron task which is in final state
func updateTasks(replyTaskRespBody *ReplyTaskRespBody, taskReplySnapshot []ReplyTask) {
	for _, taskInvkID := range replyTaskRespBody.SucceededList {
		task, ok := TaskMap[taskInvkID]
		if !ok {
			logs.GetAssistantLogger().Warnf("Task invocation(ID:%s) does not exist in task map.", taskInvkID)
			break
		}

		for index := range taskReplySnapshot {
			taskReply := &taskReplySnapshot[index]
			if taskInvkID == taskReply.InvocationID &&
				isFinalState(task.isCronTask(), taskReply.Status) {
				delete(TaskMap, taskInvkID)
				logs.GetAssistantLogger().Infof("Task invocation(ID:%s) delete successfully.", taskInvkID)
			}
		}
	}
}

func replyTaskExec() {
	replyTaskRequestBodyPtr := getTasks()
	if len(replyTaskRequestBodyPtr.Tasks) != 0 {
		replyTaskRespBody, err := replyTasks(replyTaskRequestBodyPtr)
		if err == nil {
			updateTasks(replyTaskRespBody, replyTaskRequestBodyPtr.Tasks)
		}
	}
}

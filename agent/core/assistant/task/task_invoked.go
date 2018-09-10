package task

import (
	"container/list"
	"io/ioutil"
	"net/http"
	"time"
	"io"
	"bytes"
	assistant_utils "github.com/huaweicloud/telescope/agent/core/assistant/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/robfig/cron"
)

// PullTasksTicker ...
func PullTasksTicker(switchChan chan bool) {
	// TODO 从配置文件或server端获取间隔时间
	ticker := time.NewTicker(time.Second * assistant_utils.PULL_TASK_INTERVAL_IN_SECOND)
	switchFlag := false
	for {
		select {
		case t := <-ticker.C:
			if switchFlag {
				logs.GetAssistantLogger().Infof("Pull tasks starts at %s", t.String())
				go pullTasksExec()
			}
		}
		select {
		case switchFlag = <-switchChan:
			logs.GetAssistantLogger().Debugf("Pull tasks switch is: %v", switchFlag)
		default:
		}
	}
}

func pullTasksExec() {
	tasks := pullTasks()

	logs.GetAssistantLogger().Debugf("[pullTasksExec] task pulled is: %v", tasks)
	for index := range tasks {
		task := Task{
			TaskPulled:         &tasks[index],
			TaskExecEntityList: list.New(),
		}
		if _, ok := TaskMap[task.TaskPulled.InvocationID]; ok {
			logs.GetAssistantLogger().Infof("Task(%s) has been existed in TaskMap, continue", task.TaskPulled.InvocationID)
			continue
		}
		TaskMap[task.TaskPulled.InvocationID] = &task
		logs.GetAssistantLogger().Infof("I am set new to map %s", task.TaskPulled.InvocationID)

		taskExecEntity := InstantiateTaskExecEntity(&tasks[index])
		task.pushTaskExec(taskExecEntity)
		if task.isCronTask() {
			task.Cron = cron.NewWithLocation(getCSTLocation())
			task.Cron.AddFunc(task.TaskPulled.Cron, func() {
				logs.GetAssistantLogger().Debugf("Cron job in, invocation id: %s.", task.TaskPulled.InvocationID)
				taskExecEntity := InstantiateTaskExecEntity(&tasks[index])
				if taskExecEntity.TaskInterface == nil {
					logs.GetAssistantLogger().Debug("TaskExecEntity.TaskInterface is nil.")
				}
				task.pushTaskExec(taskExecEntity)
				go StartTaskExecEntity(taskExecEntity)
			})
			task.Cron.Start()
		} else {
			go StartTaskExecEntity(taskExecEntity)
		}
	}
}

func pullTasks() []TaskPulled {
	logs.GetAssistantLogger().Debug("[pullTasks]Enter pullTasks")
	tasks := []TaskPulled{}

	url := assistant_utils.BuildURL(assistant_utils.PULL_TASK_URI + utils.SLASH + utils.GetConfig().InstanceId)
	apptask := utils.GetConfig().AppTask
	var body io.Reader
	if apptask != "" {
		requestBody := PullTaskRequestBody{
			InvokeScope: "APPTASK",
			Resources:   utils.GetConfig().AppTask,
		}
		requestBodyInBytes, err := json.Marshal(requestBody)
		if err != nil {
			logs.GetAssistantLogger().Errorf("Failed marshall request body for pull task", requestBody)
		}
		body = bytes.NewBuffer(requestBodyInBytes)
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[pullTasks]Create request Error: %s", err.Error())
	}
	logs.GetAssistantLogger().Debug("[pullTasks]Create request successfully.")

	res, err := utils.HTTPSend(request, assistant_utils.SERVICE)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[pullTasks]Failed to request, url is %s, error is %s", url, err.Error())
		return tasks
	}
	logs.GetAssistantLogger().Debug("[pullTasks]Request send successfully.")

	if res.StatusCode != http.StatusOK {
		logs.GetAssistantLogger().Errorf("[pullTasks]Request server succeeded but status code is %d(should be %d), url is %s", res.StatusCode, http.StatusOK, url)
		return tasks
	}
	logs.GetAssistantLogger().Debugf("[pullTasks] res.StatusCode is: %d", res.StatusCode)

	defer res.Body.Close()
	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[pullTasks] ioutil.ReadAll failed, error is: %s, res.Body is: %v", err.Error(), res.Body)
		return tasks
	}
	logs.GetAssistantLogger().Debugf("[pullTasks]Response: %s", string(resBodyBytes))

	pullTaskRespBody := PullTaskRespBody{}
	err = json.Unmarshal(resBodyBytes, &pullTaskRespBody)
	if err != nil {
		logs.GetAssistantLogger().Errorf("[pullTasks]Failed to unmarshal response [%s], error is %s.", string(resBodyBytes), err.Error())
		return []TaskPulled{}
	}
	tasks = pullTaskRespBody.Tasks

	logs.GetAssistantLogger().Debugf("[pullTasks]Pull task is: %v", tasks)
	return tasks
}

// todo InstantiateTaskInterface 初始化失败场景需要考虑
// InstantiateTaskExecEntity ...
func InstantiateTaskExecEntity(taskPulled *TaskPulled) *TaskExecEntity {
	taskExecEntity := TaskExecEntity{
		TaskInterface: InstantiateTaskInterface(taskPulled),
		TaskPulled:    taskPulled,
		EventChan:     make(chan string, MAX_EVENT_CHAN_SIZE),
		ExitCodeChan:  make(chan int, MAX_RETURN_CODE_CHAN_SIZE),
		ExitChan:      make(chan bool, MAX_EXIT_CHAN_SIZE),
		OutputChan:    make(chan *string, MAX_OUTPUT_CHAN_SIZE),
		Output:        nil,
		State:         STATE_CREATED,
		States:        list.New(),
		Events:        list.New(),
	}

	logs.GetAssistantLogger().Debugf("Initialize TaskExecEntity is: %v", taskExecEntity)
	return &taskExecEntity
}

// InstantiateTaskInterface ...
func InstantiateTaskInterface(taskPulled *TaskPulled) TaskInterface {
	logs.GetAssistantLogger().Debugf("Initialize TaskInterface, command is: %s", taskPulled.Command)
	switch taskPulled.Command {
	case TASK_COMMAND_RUN_SHELL:
		return InstantiateTaskRunShell(taskPulled)
	case TASK_COMMAND_CANCEL:
		return InstantiateTaskCancel(taskPulled)
	default:
		logs.GetAssistantLogger().Errorf("Initialize TaskInterface failed, command is: %s", taskPulled.Command)
		return nil
	}
}

// StartTaskExecEntity ...
func StartTaskExecEntity(taskExecEntity *TaskExecEntity) {
	logs.GetAssistantLogger().Debug("[StartTaskExecEntity] Enter StartTaskExecEntity")
	taskExecEntity.RunAndListen()
}

package windowslog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
)

type WindowsOsLogRecorder struct {
	OsLogChannelStates []*WindowsOsLogChannelState
}

type WindowsOsLogChannelState struct {
	Channel  string
	RecordId uint64
	Finished bool
}

func (recorder *WindowsOsLogRecorder) InitRecord() {
	fileStates := getStatesFromLocalRecordFile()
	recorder.OsLogChannelStates = fileStates
}

//从本地读取windows_os_log_path_record.json文件
func getStatesFromLocalRecordFile() []*WindowsOsLogChannelState {
	dir, _ := os.Getwd()
	recordFilePath := dir + lts_utils.WINDOWS_OS_LOG_RECORD_FILE_PATH
	_, err := os.Stat(recordFilePath)
	states := []*WindowsOsLogChannelState{}
	if os.IsNotExist(err) {
		recordPath := filepath.Dir(recordFilePath)
		err := os.MkdirAll(recordPath, 0750)
		if err != nil {
			logs.GetLtsLogger().Errorf("Failed to created windows os log record file dir %s, error:%s", recordPath, err.Error())
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}
		f, createError := os.Create(recordFilePath) //create file
		defer f.Close()
		if createError != nil {
			logs.GetLtsLogger().Errorf("Failed to windows os log record record file:%s, error is %s ", recordFilePath, createError.Error())
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}
		return states
	} else {
		openfile, openError := os.Open(recordFilePath) //open file
		defer openfile.Close()
		if openError != nil {
			logs.GetLtsLogger().Errorf("Failed to open file: %s, error is %s", recordFilePath, openError.Error())
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}

		decoder := json.NewDecoder(openfile)
		err := decoder.Decode(&states)
		//每次重启加载json文件的时候，把所有的finished 置为true
		for stateIndex := range states {
			states[stateIndex].Finished = true
		}
		if err != nil {
			logs.GetLtsLogger().Warn("Failed to decode windows os log record and reset states to empty")
		}
		return states
	}
}

//从内存中的recorder根据channel name 获取相应的channel状态
func (recorder *WindowsOsLogRecorder) GetChannelStateInRecorder(channelName string) *WindowsOsLogChannelState {
	if len(recorder.OsLogChannelStates) == 0 {
		return nil
	}
	for _, channelState := range recorder.OsLogChannelStates {
		if strings.Compare(channelState.Channel, channelName) == 0 {
			return channelState
		}
	}
	return nil
}

func (recorder *WindowsOsLogRecorder) AddChannelState(channelState *WindowsOsLogChannelState) {
	recorder.OsLogChannelStates = append(recorder.OsLogChannelStates, channelState)
}

func (recorder *WindowsOsLogRecorder) UpdateChannelStateRecordId(channelName string, recordId uint64) {
	for _, channelState := range recorder.OsLogChannelStates {
		if channelState.Channel == channelName {
			channelState.RecordId = recordId
			channelState.Finished = true
		}
	}

	//持久化
	dir, _ := os.Getwd()
	path := dir + lts_utils.WINDOWS_OS_LOG_RECORD_FILE_PATH
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to open windows os log record file : %s, error is %s.", path, err.Error())
		logs.GetLtsLogger().Flush()
		os.Exit(1)
	}
	encoder := json.NewEncoder(f)
	err = encoder.Encode(recorder.OsLogChannelStates)
	if err != nil {
		f.Close()
		logs.GetLtsLogger().Errorf("Error when encoding the windows os log states: %s", err)
		return
	}

}

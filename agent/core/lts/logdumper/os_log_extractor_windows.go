package logdumper

import (
	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"

	"github.com/huaweicloud/telescope/agent/core/utils"

	"github.com/huaweicloud/telescope/agent/core/lts/windowslog"
)

func (e *extractor) readWindowsOsLog(recorder *windowslog.WindowsOsLogRecorder, data chan FileEvent) {
	//遍历收集多个channel的日志
	for _, channelName := range e.windowsOsLogChannels {
		channelState := recorder.GetChannelStateInRecorder(channelName)
		if channelState == nil {
			channelState = &windowslog.WindowsOsLogChannelState{Channel: channelName, RecordId: 0, Finished: true}
			recorder.AddChannelState(channelState)
		}
		//读取日志前先置为false，有且只有当本次收集的文本日志成功发送到服务端，该channel state的状态才能设置为true
		channelState.Finished = false
		collector, err := windowslog.NewWindowsLogCollector(channelState.Channel)
		if err != nil {
			logs.GetLtsLogger().Errorf("Failed to create windows log collector, error: %s", err.Error())
			continue
		}
		//从上一次保存下来的读取位置继续读取日志
		collector.ReadWindowsOsLogFromRecordId(channelState.RecordId)

		//封装服务端需要的data model，放到data channel供发送模块消费
		logArr := collector.WindowsOsLogs
		if logArr != nil && len(logArr) > 0 {
			localIp := utils.GetLocalIp()
			hostName := utils.GetHostName()
			events := make([]Event, 0, lts_utils.WINDOWS_OS_LOG_PER_COLLECT_MAX_NUMBER)
			for i := range logArr {
				if utils.GetCurrTSInMs()-logArr[i].TimeCreated <= lts_utils.LOG_File_VALID_DURATION {
					eventBytes, err := jsonx.Marshal(logArr[i])
					if err != nil {
						logs.GetLtsLogger().Errorf("Failed to marshal windows os log. error: %s", err.Error())
						continue
					}
					events = append(events, Event{Message: string(eventBytes), Time: uint64(logArr[i].TimeCreated), Path: "System", Ip: localIp, HostName: hostName})
				}
			}

			if len(events) > 0 {
				logEventMsg := LogEventMessage{LogEvents: events, LogGroupId: e.groupId, LogTopicId: e.topicId}
				data <- FileEvent{IsWindowsOsLog: true, WindowsOsLogChannelState: *channelState, LogEvent: logEventMsg, Offset: collector.LastRecordId}
				logs.GetLtsLogger().Debugf("Channel [%s], start record id is %v, end record is %v.", channelState.Channel, channelState.RecordId, collector.LastRecordId)
			} else {
				logs.GetLtsLogger().Warn("The logs are ignored due to time invalid.")
				recorder.UpdateChannelStateRecordId(channelState.Channel, collector.LastRecordId)
			}
		}

	}

}

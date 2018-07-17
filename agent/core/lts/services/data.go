package services

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/config"
	"github.com/huaweicloud/telescope/agent/core/lts/errs"
	"github.com/huaweicloud/telescope/agent/core/lts/logdumper"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

//收集文本日志
func StartExtractionTask() {
	logs.GetLtsLogger().Debug("Create extractors.")
	logdumper.CreateExtractors()
	logdumper.InitExtractorsFilesOffset()
	logs.GetLtsLogger().Info("Start to extract...")
	cronTime := lts_utils.COLLECT_LOG_CRON_JOB_TIME_SECOND //set default
	ticker := time.NewTicker(time.Duration(cronTime) * time.Second)
	go func() {
		for _ = range ticker.C {
			if config.GetConfig().Enable {
				logs.GetLtsLogger().Debugf("LTS enabled is %v", config.GetConfig().Enable)
				logdumper.Extract(GetchData())
			}
		}
	}()

	windowsTicker := time.NewTicker(60 * time.Second)
	go func() {
		for _ = range windowsTicker.C {
			if config.GetConfig().Enable {
				logs.GetLtsLogger().Debugf("Windows os log LTS enabled is %v", config.GetConfig().Enable)
				logdumper.ExtractOsLog(GetchData())
			}
		}
	}()

}

//发送日志数据到服务端
func StartDataService(data chan logdumper.FileEvent) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	for {
		fileEvent := <-data
		logs.GetLtsLogger().Infof("A group of logs from [%s] arrived and await to send to server...", fileEvent.FileState.FilePath)
		fileEvent.SendLogDataToServer(client)

		if !fileEvent.SuccessPreProcessLogEvent {
			logs.GetLtsLogger().Errorf("Failed to process log data from log file [%s], the log data is dismissed.", fileEvent.FileState.FilePath)
			if fileEvent.IsWindowsOsLog {
				logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.Offset)
			} else {
				logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.Offset, fileEvent.FileState.LineNumber+uint64(len(fileEvent.LogEvent.LogEvents)))
			}
			continue
		}

		switch fileEvent.ResStatusCode {
		case http.StatusCreated:
			logs.GetLtsLogger().Infof("Finished to send logs to server.")
			if fileEvent.IsWindowsOsLog {
				logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.Offset)
			} else {
				logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.Offset, fileEvent.FileState.LineNumber+uint64(len(fileEvent.LogEvent.LogEvents)))
			}

		case http.StatusBadRequest:
			if isNeedDrop(fileEvent.ErrRes.Message.Code) {
				logs.GetLtsLogger().Errorf("Failed to send log, drop it. Response is %s, and the log file is [%s].", fileEvent.ErrRes.Message.Details, fileEvent.FileState.FilePath)
				if fileEvent.IsWindowsOsLog {
					logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.Offset)
				} else {
					logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.Offset, fileEvent.FileState.LineNumber+uint64(len(fileEvent.LogEvent.LogEvents)))
				}
				errs.PutLtsDetail(fileEvent.ErrRes.Message.Code, fileEvent.ErrRes.Message.Details)
			} else if strings.Contains(fileEvent.ErrRes.Message.Code, "LTS") {
				logs.GetLtsLogger().Errorf("Failed to send log, not drop. Response is %s, and the log file is [%s].", fileEvent.ErrRes.Message.Details, fileEvent.FileState.FilePath)
				if fileEvent.IsWindowsOsLog {
					logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.WindowsOsLogChannelState.RecordId)
				} else {
					logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.FileState.OffSet, fileEvent.FileState.LineNumber)
				}
				errs.PutLtsDetail(fileEvent.ErrRes.Message.Code, fileEvent.ErrRes.Message.Details)
				if fileEvent.ErrRes.Message.Code == "LTS.0307" {
					time.Sleep(lts_utils.PUT_LOG_OVER_LIMIT_WAIT_MINITUES * time.Minute)
				}
			} else {
				logs.GetLtsLogger().Errorf("Failed to send log and response code is 400, response is {%s} ,the log file [%s] will be continued to collect.", fileEvent.ResponseStr, fileEvent.FileState.FilePath)
				if fileEvent.IsWindowsOsLog {
					logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.WindowsOsLogChannelState.RecordId)
				} else {
					logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.FileState.OffSet, fileEvent.FileState.LineNumber)
				}
				errs.PutLtsDetail(errs.BAD_REQUEST_UNKONOW_ERR.Code, fileEvent.ResponseStr)
			}
		case http.StatusUnauthorized:
			if strings.Contains(fileEvent.ErrRes.Message.Code, "LTS") {
				logs.GetLtsLogger().Errorf("Current user has no access to put log events,please check user status,response is {%s}", fileEvent.ErrRes.Message.Details)
				if fileEvent.IsWindowsOsLog {
					logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.WindowsOsLogChannelState.RecordId)
				} else {
					logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.FileState.OffSet, fileEvent.FileState.LineNumber)
				}
				errs.PutLtsDetail(fileEvent.ErrRes.Message.Code, fileEvent.ErrRes.Message.Details)
			} else {
				logs.GetLtsLogger().Errorf("Failed to send log due to authorization failed, response is {%s},please check configuration file conf.json. ", fileEvent.ResponseStr)
				if fileEvent.IsWindowsOsLog {
					logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.WindowsOsLogChannelState.RecordId)
				} else {
					logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.FileState.OffSet, fileEvent.FileState.LineNumber)
				}
				errs.PutLtsDetail(errs.AUTHORIZATION_FAILED_UNKOWN_ERR.Code, fileEvent.ResponseStr)
			}

		default:
			logs.GetLtsLogger().Errorf("Failed to send log, status code is %v, please check above logs.", fileEvent.ResStatusCode)
			if fileEvent.IsWindowsOsLog {
				logdumper.WindowsOsLogRecord.UpdateChannelStateRecordId(fileEvent.WindowsOsLogChannelState.Channel, fileEvent.WindowsOsLogChannelState.RecordId)
			} else {
				logdumper.Record.UpdateRecord(fileEvent.FileState, fileEvent.FileState.OffSet, fileEvent.FileState.LineNumber)
			}
			errs.PutLtsDetail(errs.SEND_LOG_DATA_FAILED.Code, errs.SEND_LOG_DATA_FAILED.Message)
		}
	}
}

//whether the logs need drop or not
//1. log format error, need drop
//2. log topic/group not existed and exceed quota, can't drop
func isNeedDrop(statusCode string) bool {
	needDropStatusCodes := strings.Split(lts_utils.LOGS_NEED_DROP_ERR_CODES, ",")
	if utils.StrArrContainsStr(needDropStatusCodes, statusCode) {
		return true
	}

	return false
}

package logdumper

import (
	"os"
	"regexp"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_error "github.com/huaweicloud/telescope/agent/core/lts/errs"
	"github.com/huaweicloud/telescope/agent/core/lts/file"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"

	"github.com/huaweicloud/telescope/agent/core/utils"
)

type extractor struct {
	groupId                     string
	topicId                     string
	path                        string
	timeExtractMode             string
	timeExtractPattern          string
	singleLineLog               bool
	multiLineLogMarkTimePattern string //"YYYY-MM-DD"
	multiLineLogMarkRegex       string
	curCollector                *Collector
	recentLogParseTime          time.Time
	files                       []string // files which can be extracted
	isOsLog                     bool
	windowsOsLogChannels        []string //windows操作系统需要收集的操作系统日志channel:system,security...,目前只收集system
}

// Read the target log file, returns the log text
func (e *extractor) read() error { //make sure the extractor's log is empty when read again
	files := e.files
	//collect the log text from files
	for fileNameIndex := range files {
		fileInfo, _ := os.Stat(files[fileNameIndex])
		finger := utils.GetFileFingerPrint(files[fileNameIndex])
		newState := file.FileState{FilePath: files[fileNameIndex], Info: fileInfo, FileStateOS: file.GetOSState(fileInfo), FingerPrint: finger, Timestamp: time.Now()}
		oldState := Record.FileStateExistInRecord(newState)
		collector := &Collector{fileState: newState}

		if oldState == nil {
			logs.GetLtsLogger().Debugf("Start to collect new file:%s.", files[fileNameIndex])
			collector.collectNewFile(e, newState, 0)
			e.curCollector = collector
			return nil
		} else if uint64(newState.Info.Size()) > oldState.OffSet && oldState.Finished { //if a file is collecting, we should not collect again,if file size not changed means no new logs write into the file,so no need to collect
			logs.GetLtsLogger().Debugf("Start to collect old file:%s, its offset is:%d.", files[fileNameIndex], oldState.OffSet)
			collector.collectExistingFile(e, *oldState, newState)
			e.curCollector = collector
			return nil
		} else {
			logs.GetLtsLogger().Debugf("Don't need to collect file :%s.", files[fileNameIndex])
			e.curCollector = nil
			continue
		}
	}
	return nil
}

// Assemble the logs to the Event type objects
func (e *extractor) assemble(data chan FileEvent) {
	collector := e.curCollector
	if collector == nil {
		logs.GetLtsLogger().Debug("Current collector is nil, so go to next loop.")
		return
	}
	localIp := utils.GetLocalIp()
	hostName := utils.GetHostName()

	logArr := collector.readLogs
	if logArr != nil && len(logArr) > 0 {
		timePattern := e.timeExtractPattern
		//Os log 的时间格式1 ：Jun 12 10:34:48 * The server is now ready to accept connections on port 6379
		//Os log 的时间格式2 ：Jun  2 10:34:48 * The server is now ready to accept connections on port 6379
		if e.isOsLog {
			logContentBytes := []byte(logArr[0])
			if len(logContentBytes) >= 5 && logContentBytes[4] == ' ' {
				timePattern = "MMM  D hh:mm:ss"
			} else {
				timePattern = "MMM DD hh:mm:ss"
			}
		}
		reg := utils.GetReg(timePattern)
		r := regexp.MustCompile(reg)
		var (
			logTime  time.Time
			timeLocs []int
			err      error
		)

		events := make([]Event, 0, lts_utils.PER_FILE_EVENT_LOGS_MAX_NUMBER/2)
		lineNumber := collector.fileState.LineNumber
		for logIndex := range logArr {
			if len(logArr[logIndex]) > 0 {
				lineNumber = lineNumber + 1
				if lts_utils.TIME_EXTRACT_MODE_SYSTEM == e.timeExtractMode {
					//利用系统时间作为日志时间
					events = append(events, Event{Message: logArr[logIndex], Time: uint64(utils.GetCurrTSInMs()), Path: collector.fileState.FilePath, Ip: localIp, HostName: hostName, LineNumber: lineNumber})
				} else {
					//从日志里提取日志时间
					if len(timeLocs) != 2 {
						timeLocs = r.FindIndex([]byte(logArr[logIndex]))
					}
					if len(timeLocs) == 2 && len(logArr[logIndex]) >= timeLocs[1] {
						logTime, err = utils.ParseLogTime(string([]byte(logArr[logIndex])[timeLocs[0]:timeLocs[1]]), timePattern, e.recentLogParseTime)
						if err != nil {
							logs.GetLtsLogger().Warnf("Failed to get log time and set it recent time, error is :%s, log content is %s, the log file is [%s]", err.Error(), logArr[logIndex], collector.fileState.FilePath)
							logTime = e.recentLogParseTime
							timeLocs = r.FindIndex([]byte(logArr[logIndex]))
						}
					} else if e.recentLogParseTime.UnixNano() > 0 {
						logTime = e.recentLogParseTime
					} else {
						logTime = time.Now()
					}

					if logTime.Year() == 0 {
						logTime = logTime.AddDate(time.Now().Year(), 0, 0)
					}
					e.recentLogParseTime = logTime
					if len(logArr[logIndex]) > lts_utils.CONTENT_LENGTH_LIMIT_PER_LOG_TEXT {
						logArr[logIndex] = utils.SubStr(logArr[logIndex], lts_utils.CONTENT_LENGTH_LIMIT_PER_LOG_TEXT) //cut out it if the log text is too long.
					}

					if utils.GetCurrTSInMs()-utils.GetMsFromTime(logTime) <= lts_utils.LOG_File_VALID_DURATION {
						events = append(events, Event{Message: logArr[logIndex], Time: uint64(utils.GetMsFromTime(logTime)), Path: collector.fileState.FilePath, Ip: localIp, HostName: hostName, LineNumber: lineNumber})
					}
				}
			}
		}

		if len(events) > 0 {
			logEventMsg := LogEventMessage{LogEvents: events, LogGroupId: e.groupId, LogTopicId: e.topicId}
			data <- FileEvent{FileState: collector.fileState, LogEvent: logEventMsg, Offset: collector.newOffSet}
			logs.GetLtsLogger().Debugf("File [%s], old state offset is %v, new offset is %v.", collector.fileState.FilePath, collector.fileState.OffSet, collector.newOffSet)
		} else {
			logs.GetLtsLogger().Warn("The logs are ignored due to time invalid.")
			logs.GetLtsLogger().Debugf("File [%s], old state offset is %v, new offset is %v.", collector.fileState.FilePath, collector.fileState.OffSet, collector.newOffSet)
			Record.UpdateRecord(collector.fileState, collector.newOffSet, lineNumber)
		}
	}
}

//get files which in the below cases:
//1.the file is new
//2.the file size is over the offset that represent the file has been read and the file is not reading in other goroutine
func (e *extractor) filterFiles() {
	needExtract := false
	files, err := utils.GetAllFilesFromDirectoryPath(e.path)
	if err != nil || len(files) < 1 {
		logs.GetLtsLogger().Debugf("There is no files under the extractor path: %s.", e.path)
		lts_error.PutLtsDetail(lts_error.NO_FILES_FOUNT.Code, lts_error.NO_FILES_FOUNT.Message)
	}

	files = utils.FileListSortTimeAsc(files) //sort in modification time ascending order
	//remove the file which produced in past 7days

	if lts_utils.TIME_EXTRACT_MODE_SYSTEM == e.timeExtractMode {
		//用系统时间作为日志时间的场景 只收集agent启动后新产生的日志
		files = file.FilterOldLogFile(files, 24*60*60*1000)
	} else {
		files = file.FilterOldLogFile(files, lts_utils.LOG_File_VALID_DURATION)
	}

	//filter the files which filestate Finished is false and which size is not increased after last read
	for fileIndex := range files {
		fileInfo, statErr := os.Stat(files[fileIndex])
		if statErr != nil {
			logs.GetLtsLogger().Warnf("Stat file [%s] error: %s", files[fileIndex], statErr.Error())
			continue
		}
		finger := utils.GetFileFingerPrint(files[fileIndex])
		currentState := file.FileState{FilePath: files[fileIndex], Info: fileInfo, FileStateOS: file.GetOSState(fileInfo), FingerPrint: finger}
		oldState := Record.fileStates.FindPrevious(currentState)

		if oldState == nil {
			logs.GetLtsLogger().Infof("Start to collect a new log file: [%s].", currentState.FilePath)
			needExtract = true
			break
		} else if oldState != nil && oldState.OffSet < uint64(fileInfo.Size()) && oldState.Finished {
			logs.GetLtsLogger().Debugf("The file [%s] old state is: %v", files[fileIndex], oldState.Finished)
			needExtract = true
			break
		}
	}

	if needExtract {
		e.files = files
	} else {
		e.files = make([]string, 0)
	}
}

//初始化文件的offset，对于使用系统时间作为日志时间的场景，应该只收集agent启动后的日志
func (e *extractor) initFileOffset() {
	if lts_utils.TIME_EXTRACT_MODE_SYSTEM == e.timeExtractMode {
		files, err := utils.GetAllFilesFromDirectoryPath(e.path)
		if err != nil || len(files) < 1 {
			logs.GetLtsLogger().Debugf("There is no files under the extractor path: %s.", e.path)
			lts_error.PutLtsDetail(lts_error.NO_FILES_FOUNT.Code, lts_error.NO_FILES_FOUNT.Message)
		}

		files = utils.FileListSortTimeAsc(files) //sort in modification time ascending order

		files = file.FilterOldLogFile(files, lts_utils.LOG_FILE_SYSTEM_TIME_VALID_DURATION)

		//把文件的offset更新到record
		for fileIndex := range files {
			fileInfo, statErr := os.Stat(files[fileIndex])
			if statErr != nil {
				logs.GetLtsLogger().Warnf("Stat file [%s] error: %s", files[fileIndex], statErr.Error())
				continue
			}
			finger := utils.GetFileFingerPrint(files[fileIndex])
			currentState := file.FileState{FilePath: files[fileIndex], Info: fileInfo, FileStateOS: file.GetOSState(fileInfo), FingerPrint: finger}
			oldState := Record.fileStates.FindPrevious(currentState)

			if oldState == nil {
				logs.GetLtsLogger().Infof("New file [%s] need to collect from offset [%v]", currentState.FilePath, fileInfo.Size())
				Record.UpdateRecord(currentState, uint64(fileInfo.Size()), 0)
			}
		}
	}
}

package logdumper

import (
	"github.com/huaweicloud/telescope/agent/core/lts/file"
	"github.com/huaweicloud/telescope/agent/core/lts/file/read"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type Collector struct {
	fileState file.FileState
	newOffSet uint64
	readLogs  []string
}

//collect logs from file which is not recorded in the record
func (c *Collector) collectNewFile(e *extractor, fileState file.FileState, offset uint64) {
	fileState.Finished = false //before reading file ,lock
	Record.UpdateState(fileState)
	filePath := fileState.FilePath

	logs, newOffset := collectFixedLengthFromOffset(e, filePath, offset, lts_utils.PER_FILE_EVENT_LOGS_MAX_TOTAL_SIZE)
	c.newOffSet = newOffset
	c.readLogs = logs
}

//collect logs from file which is stated in record file
func (c *Collector) collectExistingFile(e *extractor, oldState file.FileState, newState file.FileState) {
	oldState.Finished = false //before reading file ,lock
	Record.UpdateState(oldState)

	filepath := newState.FilePath //use newState beacause old state the file name may be renamed when a log file is written fully
	offSet := oldState.OffSet

	logs, newOffSet := collectFixedLengthFromOffset(e, filepath, offSet, lts_utils.PER_FILE_EVENT_LOGS_MAX_TOTAL_SIZE)
	c.newOffSet = newOffSet
	c.readLogs = logs
	c.fileState = oldState

	//update filepath when log file is renamed(eg. log file is full and rotate a new log file)
	if oldState.FilePath != "" && oldState.FilePath != newState.FilePath && oldState.Finished == true {
		oldState.FilePath = newState.FilePath
		Record.UpdateState(oldState)
	}
}

//collect the log data from offset
func collectFixedLengthFromOffset(e *extractor, filepath string, offset uint64, maxSize int64) (logs []string, newOffset uint64) {
	//collect the single line log
	if e.singleLineLog {
		logs, newOffset = read.ReadLogsTextByLine(filepath, offset, maxSize)
		return
	}

	//collect the multiple line log according to the time pattern
	if "" != e.multiLineLogMarkTimePattern {
		regex := utils.GetReg(e.multiLineLogMarkTimePattern)
		logs, newOffset = read.ReadLogsTextsByRegex(filepath, offset, maxSize, regex)
		return
	}

	//collect the multiple line log according to the regex
	if "" != e.multiLineLogMarkRegex {
		logs, newOffset = read.ReadLogsTextsByRegex(filepath, offset, maxSize, e.multiLineLogMarkRegex)
		return
	}

	//if no configuration in config file,collect single line as a log data
	logs, newOffset = read.ReadLogsTextByLine(filepath, offset, maxSize)
	return
}

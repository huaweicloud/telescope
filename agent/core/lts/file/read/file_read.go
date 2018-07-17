package read

import (
	"bufio"
	"io"
	"os"
	"regexp"

	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
)

const READ_LOG_METHOD = "{time_pattern}"

//validate
func validateFileEventLogsLength(lineLog string, logsLength uint64, logs []string, maxSize uint64) (valid bool) {
	lenLogLine := len(lineLog)
	logsLength = logsLength + uint64(lenLogLine)
	valid = logsLength > maxSize || len(logs) >= lts_utils.PER_FILE_EVENT_LOGS_MAX_NUMBER
	return !valid
}

// read logs according to regular express
func ReadLogsTextsByRegex(filePath string, offset uint64, maxSize int64, reg string) (logTexts []string, newOffset uint64) {
	newOffset = offset
	file, openFileErr := os.Open(filePath)
	defer file.Close()
	if openFileErr != nil {
		logs.GetLtsLogger().Errorf("Can't open the file:%s\n", filePath)
		return
	}
	_, err := file.Seek(int64(offset), os.SEEK_SET)
	if err != nil {
		logs.GetLtsLogger().Errorf("Seek file Error:%s\n", err.Error())
		return
	}
	buf := bufio.NewReader(file)
	r := regexp.MustCompile(reg)
	//read line in the file
	var validate bool
	logText := ""
	multi := false
	validLogLen := 0
	var isLineHeader bool
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err == io.EOF {
			logs.GetLtsLogger().Infof("The file is read fully:%s.", filePath)
			logTexts = append(logTexts, logText)
			return
		}
		if err != nil {
			logs.GetLtsLogger().Errorf("Read file error:%s", err.Error())
			return
		}

		lineHeader := r.FindIndex([]byte(line))

		if len(lineHeader) == 2 {
			isLineHeader = true
		} else {
			isLineHeader = false
		}
		if isLineHeader {
			multi = true
			if logText != "" {
				validate = validateFileEventLogsLength(logText, uint64(validLogLen), logTexts, uint64(maxSize)) //validate indicate that logs array size and logs content size both are limited
				if validate {
					logTexts = append(logTexts, logText)
					validLogLen = validLogLen + len(logText)
				} else {
					newOffset = newOffset - uint64(len(logText))
					return
				}
			}
			logText = line
			newOffset = newOffset + uint64(len(line))
		} else if !isLineHeader && multi == true {
			logText = logText + line
			newOffset = newOffset + uint64(len(line))
		} else {
			multi = false
			newOffset = newOffset + uint64(len(line))
			logs.GetLtsLogger().Errorf("The log %s cant't match the time pattern or configured regex.", line)
			continue
		}
	}
	return
}

//read logs, a line is a log
func ReadLogsTextByLine(filePath string, offset uint64, maxSize int64) (logTexts []string, newOffset uint64) {
	var validate bool
	newOffset = offset
	file, openFileErr := os.Open(filePath)
	defer file.Close()
	if openFileErr != nil {
		logs.GetLtsLogger().Errorf("Can't open the file:%s, and error is %s.", filePath, openFileErr.Error())
		return
	}
	_, err := file.Seek(int64(offset), os.SEEK_SET)
	if err != nil {
		logs.GetLtsLogger().Errorf("Seek file Error:%s.", err.Error())
		return
	}
	buf := bufio.NewReader(file)
	//read line in the file
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err == io.EOF {
			logs.GetLtsLogger().Infof("The file is read fully:%s", filePath)
			return
		}
		if err != nil {
			logs.GetLtsLogger().Errorf("Read file error:%s", err.Error())
			return
		}

		validate = validateFileEventLogsLength(line, newOffset-offset, logTexts, uint64(maxSize)) //validate indicate that logs array size and logs content size both are limited
		if validate {
			logTexts = append(logTexts, line)
			newOffset = newOffset + uint64(len(line))
		} else {
			break
		}

	}
	return
}

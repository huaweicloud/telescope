package logdumper

import (
	"os"
	"strings"

	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_config "github.com/huaweicloud/telescope/agent/core/lts/config"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	windowslog "github.com/huaweicloud/telescope/agent/core/lts/windowslog"
	utils "github.com/huaweicloud/telescope/agent/core/utils"
)

var Extractors []extractor
var Record Recorder
var WindowsOsLogRecord *windowslog.WindowsOsLogRecorder

func CreateExtractors() {
	Extractors = LoadExtractors()
	//Init Record
	dir, _ := os.Getwd()
	Record = Recorder{path: dir + lts_utils.RECORD_FILE_PATH}
	Record.InitRecord() //load the states from Record file into Record
	//加载收集windows操作系统日志目录的recorder
	WindowsOsLogRecord = &windowslog.WindowsOsLogRecorder{}
	WindowsOsLogRecord.InitRecord()
}

func Extract(data chan FileEvent) {
	for _, ext := range Extractors {
		if ext.isOsLog && utils.IsWindowsOs() {
			// windows系统日志的收集不在该方法内
		} else {
			// 根据目录读取文本日志
			ext.filterFiles()       //get the files can be extracted
			if len(ext.files) > 0 { //only within the scope of extractor, the files which size changed can be extracted again
				logs.GetLtsLogger().Debugf("Start to read files: %v", strings.Join(ext.files, ","))
				err := ext.read()
				if err != nil {
					logs.GetLtsLogger().Errorf("Error while reading the logs from %s, log topic id: %s, error is %s.", ext.path, ext.topicId, err.Error())
					continue
				}
			} else {
				logs.GetLtsLogger().Debugf("There is no available files to collect.")
			}
			//组装收集到的文本原始日志，放到data channel
			ext.assemble(data)
		}
	}
}

//收集windows操作系统日志
func ExtractOsLog(data chan FileEvent) {
	for _, ext := range Extractors {
		if ext.isOsLog && utils.IsWindowsOs() {
			// 读取windows系统日志
			ext.windowsOsLogChannels = []string{"System"}
			ext.readWindowsOsLog(WindowsOsLogRecord, data)
		}
	}
}

func ReloadExtractors() {
	Extractors = LoadExtractors()
}

func LoadExtractors() []extractor {
	config := lts_config.GetConfig()
	newExtractors := make([]extractor, 0)
	for _, group := range config.Groups {
		for _, e := range group.Topics {
			ext := extractor{groupId: group.GroupId, topicId: e.LogTopicId, path: e.Path, timeExtractMode: e.TimeExtractMode, singleLineLog: e.SingleLineLog, timeExtractPattern: e.TimeExtractPattern, multiLineLogMarkTimePattern: e.MultiLineLogMarkTimePattern, multiLineLogMarkRegex: e.MultiLineLogMarkRegex, isOsLog: e.IsOsLog}
			newExtractors = append(newExtractors, ext)
		}
	}

	return newExtractors
}

func InitExtractorsFilesOffset() {
	for _, ext := range Extractors {
		if ext.isOsLog && utils.IsWindowsOs() {
			// 读取windows系统日志
		} else {
			// 初始化file的offset
			ext.initFileOffset()
		}
	}
}

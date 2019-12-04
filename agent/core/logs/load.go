package logs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
)

var (
	logger          seelog.LoggerInterface
	cesLogger       seelog.LoggerInterface
	assistantLogger seelog.LoggerInterface
)

// disable: disables all library log output before load config
func disable() {
	logger = seelog.Disabled
	cesLogger = seelog.Disabled
	assistantLogger = seelog.Disabled
}

func loadAppConfig() {
	var err error
	logCommonConfig := getCommonLog()
	logger, err = seelog.LoggerFromConfigAsBytes([]byte(logCommonConfig))

	if err != nil {
		log.Println(err)
		return
	}

	cesLogConfig := getCesLog()
	cesLogger, err = seelog.LoggerFromConfigAsBytes([]byte(cesLogConfig))
	if err != nil {
		log.Println(err)
		return
	}

	assistantLogConfig := getAssistantLog()
	if assistantLogConfig == "" {
		assistantLogger = cesLogger
		return
	}
	assistantLogger, err = seelog.LoggerFromConfigAsBytes([]byte(assistantLogConfig))
	if err != nil {
		log.Println(err)
		return
	}

}

func init() {
	disable()
	loadAppConfig()
}

// GetAssistantLogger ...
func GetAssistantLogger() seelog.LoggerInterface {
	return assistantLogger
}

// GetCesLogger ...
func GetCesLogger() seelog.LoggerInterface {
	return cesLogger
}

// GetLogger ...
func GetLogger() seelog.LoggerInterface {
	return logger
}

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		GetLogger().Errorf(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

package logs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
)

var logger seelog.LoggerInterface
var cesLogger seelog.LoggerInterface
var ltsLogger seelog.LoggerInterface

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

	ltsLogConfig := getLtsLog()
	ltsLogger, err = seelog.LoggerFromConfigAsBytes([]byte(ltsLogConfig))
	if err != nil {
		log.Println(err)
		return
	}

}

func init() {
	DisableLog()
	loadAppConfig()
}

// DisableLog disables all library log output
func DisableLog() {
	logger = seelog.Disabled
	cesLogger = seelog.Disabled
	ltsLogger = seelog.Disabled

}

func GetCesLogger() seelog.LoggerInterface {
	return cesLogger
}

func GetLtsLogger() seelog.LoggerInterface {
	return ltsLogger
}

func GetLogger() seelog.LoggerInterface {
	return logger
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		GetLogger().Errorf(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

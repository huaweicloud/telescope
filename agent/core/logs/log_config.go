package logs

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	logsConfig LogsConfig
	isLoaded   bool
)

type LogsConfig struct {
	CommonConfig string `xml:"common"`
	LtsConfig    string `xml:"lts"`
	CesConfig    string `xml:"ces"`
}

var logfilenameRex = "filename=\"(.*).log\""

func findLogfileName(logConfig string) string {
	reg, err := regexp.Compile(logfilenameRex)
	if err != nil {
		return logConfig
	}
	data := reg.Find([]byte(logConfig))

	if data != nil {
		filenameConfig := string(data)
		reg, err := regexp.Compile("\"(.*)\"")
		if err != nil {
			return logConfig
		}
		filedata := reg.Find([]byte(filenameConfig))
		fileName := string(filedata)

		if fileName != "" {
			fileName = fileName[1 : len(fileName)-1]
			path := GetCurrentDirectory()
			fileName = path + "/" + fileName
			return strings.Replace(logConfig, filenameConfig, "filename=\""+fileName+"\"", -1)
		}
	}

	return logConfig
}

func getCommonLog() (config string) {
	if !isLoaded {
		LoadConfig()
	}
	return findLogfileName(logsConfig.CommonConfig)
}

func getCesLog() (config string) {
	if !isLoaded {
		LoadConfig()
	}
	return findLogfileName(logsConfig.CesConfig)
}

func getLtsLog() (config string) {
	if !isLoaded {
		LoadConfig()
	}
	return findLogfileName(logsConfig.LtsConfig)
}

func LoadConfig() {
	pwd := GetCurrentDirectory()
	content, err := ioutil.ReadFile(pwd + "/logs_config.xml")
	if err != nil {
		fmt.Printf("Load logs_config.xml failed, error is: %s.\n", err.Error())
		return
	}

	logsConfig = LogsConfig{}
	err = xml.Unmarshal(content, &logsConfig)
	if err != nil {
		fmt.Printf("Parse logs_config.xml failed, error is: %s.\n", err.Error())
		return
	} else {
		isLoaded = true
	}
}

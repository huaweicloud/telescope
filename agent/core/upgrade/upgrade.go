package upgrade

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var AgentHome string
var AgentTmpHome string

const (
	PackageInfoFile = "info"
)

func init() {
	file, _ := exec.LookPath(os.Args[0])
	AgentHome = filepath.Dir(file)
	AgentTmpHome = AgentHome + "/.tmp"

	utils.CreateDir(AgentTmpHome)
}

// check local file
// true: already download the new package
// false: should download the new package
func isDownloaded(version string) bool {
	filePath := AgentTmpHome + "/" + PackageInfoFile
	info, err := os.Stat(filePath)
	if err != nil || info.Size() == 0 {
		return false
	}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logs.GetLogger().Errorf("Read info file failed, err:%s", err.Error())
		return false
	}

	packageInfo := Info{}
	err = json.Unmarshal(bytes, &packageInfo)
	if err != nil {
		logs.GetLogger().Errorf("Parse info file failed, err:%s", err.Error())
		return false
	}

	if packageInfo.Version != version {
		logs.GetLogger().Infof("Should download the new package, version:%s", version)
		return false
	} else {
		logs.GetLogger().Info("Already download the new package.")
	}

	return true
}

// download file from remote server
func Download(url string, version string, md5str string) error {
	if isDownloaded(version) {
		return nil
	}

	//create tmp dir if not exists
	err := utils.CreateDir(AgentTmpHome)
	if err != nil {
		logs.GetLogger().Errorf("Create agent tmp dir failed, err:%s", err)
		return err
	}

	//build file info
	tmpFileName := extractNameFromUrl(url, "agent_"+version+".tar.gz")
	tmpFilePath := AgentTmpHome + "/" + tmpFileName
	info := Info{Version: version, File: tmpFilePath}

	//get new package
	bytes, err := utils.HTTPGet(url)
	if err != nil {
		logs.GetLogger().Errorf("Fetch new package failed, err: %s.", err.Error())
		return err
	}

	//wirte new package
	err = ioutil.WriteFile(tmpFilePath, bytes, 0700)
	if err != nil {
		logs.GetLogger().Errorf("Write new package to local failed, err:%s.", err.Error())
		return err
	}

	//verify file
	newMd5str := fmt.Sprintf("%x", md5.Sum(bytes))
	if strings.ToLower(md5str) != strings.ToLower(newMd5str) {
		logs.GetLogger().Errorf("New package's md5[%s] is not equals to the md5[%s] of heartbeat.", newMd5str, md5str)
		return errors.New("new package is invalid")
	}
	info.Md5 = md5str
	info.Size = len(bytes)

	//write info file
	bytes, err = json.Marshal(info)
	if err != nil {
		logs.GetLogger().Errorf("Parse info to json failed, err:%s", err.Error())
		return err
	}
	err = ioutil.WriteFile(AgentTmpHome+"/"+PackageInfoFile, bytes, 0600)
	if err != nil {
		logs.GetLogger().Errorf("Write info to local failed, err:%s.", err.Error())
		return err
	}
	logs.GetLogger().Info("Download new package success.")
	return nil
}

//url format: http://ip:port/bucketName/file.name.tar.gz
func extractNameFromUrl(url string, defaultName string) string {
	array := strings.Split(url, "/")
	if len(array) > 0 {
		return array[len(array)-1]
	}
	return defaultName
}

package upgrade

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// AgentHome ...
var AgentHome string

// AgentTmpHome ...
var AgentTmpHome string

const (
	// PackageInfoFile ...
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
	filePath := filepath.Join(AgentTmpHome, PackageInfoFile)
	info, err := os.Stat(filePath)
	if err != nil || info.Size() == 0 {
		return false
	}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logs.GetCesLogger().Errorf("Read info file failed, err: %v", err)
		return false
	}

	packageInfo := Info{}
	err = json.Unmarshal(bytes, &packageInfo)
	if err != nil {
		logs.GetCesLogger().Errorf("Parse info file failed, err: %s", err)
		return false
	}

	if packageInfo.Version != version {
		logs.GetCesLogger().Infof("Should download the new package, version is: %s", version)
		return false
	} else {
		logs.GetCesLogger().Info("Already download the new package.")
	}

	return true
}

// Download file from remote server
func Download(url string, version string, md5str string) error {
	if isDownloaded(version) {
		return nil
	}

	//create tmp dir if not exists
	err := utils.CreateDir(AgentTmpHome)
	if err != nil {
		logs.GetCesLogger().Errorf("Create agent tmp dir failed, err: %v", err)
		return err
	}

	//build file info
	tmpFileName := extractNameFromUrl(url, "agent_"+version+".tar.gz")
	// before daemon can be upgrade, tmpFilePath must use / to contract
	tmpFilePath := AgentTmpHome + "/" + tmpFileName
	info := Info{Version: version, File: tmpFilePath}

	//get new package
	bytes, err := utils.HTTPGet(url)
	if err != nil {
		logs.GetCesLogger().Errorf("Fetch new package failed, err: %v.", err)
		return err
	}

	//wirte new package
	err = ioutil.WriteFile(tmpFilePath, bytes, 0700)
	if err != nil {
		logs.GetCesLogger().Errorf("Write new package to local failed, err: %v.", err)
		return err
	}

	//verify file
	newMd5str := fmt.Sprintf("%x", md5.Sum(bytes))
	if strings.ToLower(md5str) != strings.ToLower(newMd5str) {
		logs.GetCesLogger().Errorf("New package's md5[%s] is not equals to the md5[%s] of heartbeat.", newMd5str, md5str)
		return errors.New("new package is invalid")
	}
	info.Md5 = md5str
	info.Size = len(bytes)

	//write info file
	bytes, err = json.Marshal(info)
	if err != nil {
		logs.GetCesLogger().Errorf("Parse info to json failed, err: %v", err)
		return err
	}
	filePath := filepath.Join(AgentTmpHome, PackageInfoFile)
	err = ioutil.WriteFile(filePath, bytes, 0600)
	if err != nil {
		logs.GetCesLogger().Errorf("Write info to local failed, err: %v.", err)
		return err
	}
	logs.GetCesLogger().Info("Download new package success.")
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

package upgrade

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/daemon/process"

	agent "github.com/huaweicloud/telescope/agent/core/upgrade"
	utils "github.com/huaweicloud/telescope/agent/core/utils"
	"path/filepath"
)

const (
	PackageInfoFile = "info"
)

// ScanAgentTmpDir ...
// scan tmp dir, if find new version, send signal to upgrade goroutine
func ScanAgentTmpDir(dir string, agentName string, upgradeSignal chan *agent.Info) {
	for {
		files, _ := ioutil.ReadDir(dir)
		if len(files) > 0 {
			for _, f := range files {
				if f.Name() == PackageInfoFile {
					info, err := parsePackageInfo(filepath.Join(dir, f.Name()))
					if err == nil {
						logs.GetLogger().Infof("Find new version:%s, begin to upgrade.", info.Version)
						upgradeSignal <- info
					}
					break
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}

// parse the info file
func parsePackageInfo(infoFile string) (*agent.Info, error) {
	bytes, err := ioutil.ReadFile(infoFile)
	if err != nil {
		logs.GetLogger().Errorf("Read info file failed, err:%s", err.Error())
		return nil, err
	}
	info := agent.Info{}
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		logs.GetLogger().Errorf("Parse info file failed, err:%s", err.Error())
		return nil, err
	}
	return &info, nil
}

// DoUpgrade ...
// upgrade the agent
func DoUpgrade(agentHome, agentTmpHome, agentName, daemonName string, info *agent.Info, oldProc *os.Process) (*os.Process, error) {
	defer os.RemoveAll(agentTmpHome)
	agentBinPath := filepath.Join(agentHome, agentName)
	//Uncompress new package
	destDir, err := utils.UncompressFile(info.File, agentTmpHome)
	if err != nil {
		return oldProc, err
	}

	logs.GetLogger().Info("Begin to backup agent.")

	//backup current agent
	err = backup(agentHome, agentTmpHome, agentName, daemonName)
	if err != nil {
		return oldProc, err
	}
	//kill current process
	err = process.SigAndKillProcess(filepath.Join(agentHome, agentName), agent.SIG_UPGRADE, oldProc)
	if err != nil {
		logs.GetLogger().Errorf("Stop agent process failed, err:%s", err.Error())
		rollbackFile(agentHome, agentTmpHome, agentName, daemonName)
		return oldProc, err
	}

	logs.GetLogger().Info("Kill old process finished.")

	//upgrade
	err = upgrade(agentHome, agentTmpHome, destDir, agentName, daemonName)
	if err != nil {
		return rollback(agentHome, agentTmpHome, agentName, daemonName, agentBinPath)
	}

	//start process
	proc, err := process.StartProcess(agentBinPath)
	if err == nil {
		logs.GetLogger().Info("Start new process success.")
		return proc, nil
	}

	logs.GetLogger().Errorf("Start new process failed, err:%s", err.Error())
	return rollback(agentHome, agentTmpHome, agentName, daemonName, agentBinPath)
}

// back files
func backup(agentHome, agentTmpHome, agentName, daemonName string) error {
	bakDir := filepath.Join(agentTmpHome, "bak")
	err := os.RemoveAll(bakDir)
	if err != nil {
		logs.GetLogger().Errorf("Delete backup dir failed, err:%s", err.Error())
		return err
	}

	err = utils.CreateDir(bakDir)
	if err != nil {
		logs.GetLogger().Errorf("Make backup dir failed, err:%s", err.Error())
		return err
	}

	//backup daemon, only support linux
	osName := runtime.GOOS
	if osName == "linux" {
		err := os.Rename(filepath.Join(agentHome, daemonName), filepath.Join(bakDir, daemonName))
		if err != nil {
			logs.GetLogger().Errorf("Backup daemon failed, err:%s", err.Error())
			return err
		}
	}

	//backup agent
	err = os.Rename(filepath.Join(agentHome, agentName), filepath.Join(bakDir, agentName))
	if err != nil {
		logs.GetLogger().Errorf("Backup agent failed, err:%s", err.Error())
		return err
	}

	logs.GetLogger().Info("Backup agent success.")
	return nil
}

// replace file
func upgrade(agentHome, agentTmpHome, destDir, agentName, daemonName string) error {
	//upgrade agent
	err := utils.CopyFile(filepath.Join(agentTmpHome, destDir, "bin", agentName), filepath.Join(agentHome, agentName))
	if err != nil {
		logs.GetLogger().Errorf("Upgrade agent failed, err:%s", err.Error())
		return err
	}
	err = os.Chmod(filepath.Join(agentHome, agentName), 0700)
	if err != nil {
		logs.GetLogger().Errorf("Chmod agent failed, err:%s", err.Error())
		return err
	}

	//upgrade daemon, next start daemon process will be effective
	osName := runtime.GOOS
	if osName == "linux" {
		err := utils.CopyFile(filepath.Join(agentTmpHome, destDir, "bin", daemonName), filepath.Join(agentHome, daemonName))
		if err != nil {
			logs.GetLogger().Errorf("Upgrade daemon failed, err:%s", err.Error())
			return err
		}
		err = os.Chmod(filepath.Join(agentHome, daemonName), 0700)
		if err != nil {
			logs.GetLogger().Errorf("Chmod daemonName failed, err:%s", err.Error())
			return err
		}
	}

	logs.GetLogger().Info("Upgrade agent success.")
	return nil
}

func rollback(agentHome, agentTmpHome, agentName, daemonName, agentBinPath string) (*os.Process, error) {
	err := rollbackFile(agentHome, agentTmpHome, agentName, daemonName)
	if err != nil {
		logs.GetLogger().Errorf("Rollback file error when rollback, err:%s", err.Error())
		return nil, err
	} else {
		return startProc(agentBinPath)
	}
}

func startProc(agentBinPath string) (*os.Process, error) {
	logs.GetLogger().Infof("Begin to start old agent.")
	proc, err := process.StartProcess(agentBinPath)
	if err != nil {
		logs.GetLogger().Errorf("Start old agent failed, err:%s.", err.Error())
		return nil, err
	} else {
		logs.GetLogger().Infof("Start old agent success.")
		return proc, nil
	}
}

// rollback file
func rollbackFile(agentHome, agentTmpHome, agentName, daemonName string) error {
	logs.GetLogger().Info("Begin to rollback agent.")
	bakDir := agentTmpHome + "/bak"
	//rollback agent
	err := os.Rename(filepath.Join(bakDir, agentName), filepath.Join(agentHome, agentName))
	if err != nil {
		logs.GetLogger().Errorf("Rollback agent failed, err:%s", err.Error())
		return err
	}

	//rollback daemon, only support linux
	osName := runtime.GOOS
	if osName == "linux" {
		err = os.Rename(filepath.Join(bakDir, daemonName), filepath.Join(agentHome, daemonName))
		if err != nil {
			logs.GetLogger().Errorf("Rollback daemon failed, err:%s", err.Error())
			return err
		}
	}
	logs.GetLogger().Info("Rollback agent success.")
	return nil
}

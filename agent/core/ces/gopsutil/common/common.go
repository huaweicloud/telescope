package common

import (
	"bufio"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
)

func HostProc(combineWith ...string)string{
	return GetEnv("HOST_PROC", "/proc", combineWith...)
}

//GetEnv retrieves the environment variable key. If it does not exist it returns the default.
func GetEnv(key string, defaultValue string, combineWith ...string)string{
	defer func(){
		if p := recover(); p != nil{
			logs.GetCesLogger().Errorf("panic when try to get environment variable key, error is:%v",p)
		}
	}()

	value := os.Getenv(key)
	if value == ""{
		value = defaultValue
	}

	switch len(combineWith){
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith) + 1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
	panic("invalid switch case")
}

// ReadLines reads contents from a file and splits them by new lines.
// A convenience wrapper to ReadLinesOffsetN(filename, 0, -1).
func ReadLines(filename string)([]string, error){
	return ReadLinesOffsetN(filename, 0, -1)
}

// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
// n >= 0: at most n lines
// n <0: whole file
func ReadLinesOffsetN(filename string, offset uint, n int)([]string, error){
	f, err := os.Open(filename)
	if err != nil{
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n + int(offset) || n < 0; i++{
		line, err := r.ReadString('\n')
		if err != nil{
			break
		}
		if i < int(offset){
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret, nil
}


func GetDeviceNum(filesystem string, mountpoint string) string {

	lines := getLsblkResult()

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if strings.Contains(fields[0], filesystem) && len(fields) > 6 && mountpoint == fields[6] {
			return fields[1]
		}

	}
	logs.GetCesLogger().Errorf("Can't get device number by mountpoint.")

	return ""
}

func getLsblkResult() []string {
	cmd := exec.Command("lsblk")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.GetCesLogger().Errorf("Execute cmd StdoutPipe error, error is %s", err)
		return nil
	}
	defer stdout.Close()
	defer cmd.Wait()
	if err = cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Execute cmd Start error, error is %s", err)
		return nil
	}
	reader := bufio.NewReader(stdout)
	var lines []string

	for i := 0; i >= 0; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, strings.Trim(line, "\n"))
	}

	if len(lines) < 1 {
		return nil
	}

	return lines[1:len(lines)]
}

// 根据lsblk第六列结果获取磁盘/分区类型
func GetDeviceTypeMap() map[string][]string {
	lines := getLsblkResult()
	deviceTypeMap := make(map[string][]string, 5)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		switch fields[5] {
		case "disk":
			deviceTypeMap["disk"] = append(deviceTypeMap["disk"], fields[0])
		case "part":
			deviceTypeMap["part"] = append(deviceTypeMap["part"], fields[0])
		case "lvm":
			deviceTypeMap["lvm"] = append(deviceTypeMap["lvm"], fields[0])
		default:
			logs.GetCesLogger().Infof("Other value of device type from lsblk:%v", fields[0])
		}

	}

	return deviceTypeMap
}

// deviceNum format (majNum:minNum)
func GetDeviceNameByDeviceNum(deviceNum string) string {
	deviceNumArr := strings.Split(deviceNum, ":")
	if len(deviceNum) > 1 {
		majNum := deviceNumArr[0]
		minNum := deviceNumArr[1]

		filename := HostProc("diskstats")
		lines, err := ReadLines(filename)

		if err != nil {
			return ""
		}
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			majNumFile := fields[0]
			minNumFile := fields[1]
			name := fields[2]

			if majNumFile == majNum && minNumFile == minNum {
				return name
			}
		}
	}
	return ""
}

package common

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/huaweicloud/telescope/agent/core/logs"
)

func getLsblkResult() []string {
	cmd := exec.Command("lsblk")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.GetCesLogger().Errorf("Execute cmd StdoutPipe error: %s", err.Error())
		return nil
	}
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Execute cmd Start error: %s", err.Error())
		return nil
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	var lines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, strings.Trim(line, "\n"))
	}

	if len(lines) < 1 {
		return nil
	}

	return lines[1:]
}

// GetDeviceTypeMap ...
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
			logs.GetCesLogger().Warnf("Unknown value of device type from lsblk for device:%+v", fields[0])
		}

	}

	return deviceTypeMap
}

func execLsblkWithOpts(opts []string) []string {
	cmd := exec.Command("lsblk", opts...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.GetCesLogger().Errorf("Execute cmd StdoutPipe error: %s", err.Error())
		return nil
	}
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Execute cmd Start error: %s", err.Error())
		return nil
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	var lines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, strings.Trim(line, "\n"))
	}

	if len(lines) < 1 {
		return nil
	}

	return lines[1:]
}

// GetDeviceMap returns the map
func GetDeviceMap() map[string]string {
	opts := []string{"-o", "MOUNTPOINT,MAJ:MIN"}
	lines := execLsblkWithOpts(opts)
	m := make(map[string]string)

	// trim title
	for _, line := range lines[1:] {
		// trim the line which mount point is empty
		fields := strings.Fields(line)
		fieldsCount := len(fields)
		if fieldsCount < 2 {
			continue
		}

		mountPoint := strings.Join(fields[:fieldsCount-1], "") // mount point with space
		devSN := fields[fieldsCount-1]
		m[mountPoint] = devSN
	}

	logs.GetCesLogger().Debugf("Device map is %v", m)
	return m
}

package disk

import (
	"github.com/huaweicloud/telescope/agent/core/ces/gopsutil/common"
	"strings"
)

// GetFileSystemStatus returns filesystem state from /proc/mounts,eg. "rw":0, "ro":1
func GetFileSystemStatus(names ...string)(map[string]ces_linux.FSMountStat, error){
	filename := common.HostProc("mounts")
	lines, err := common.ReadLines(filename)

	if err != nil{
		return nil, err
	}
	ret := make(map[string]ces_linux.FsMountstat, 0)
	var rwstate int64 = -1
	for _, line := range lines{
		fields := strings.Fields(line)
		if len(fields) < 4{
			// malformed line in /proc/mounts, avoid panic by ignoring
			continue
		}
		partition := fields[0]
		mountpoint := fields[1]
		if strings.Contaions(fields[3], "rw"){
			rwstate = 0
		}else if strings.Contains(fields[3], "ro"){
			rwstate =1
		}

		d := ces_linux.FSMountstat{
			Partition: partition,
			MountPoint: mountpoint,
			State: rwstate,
		}
		if partition == "" || mountpoint == ""{
			continue
		}
		ret[mountpoint] = d
	}
	return ret, nil
}

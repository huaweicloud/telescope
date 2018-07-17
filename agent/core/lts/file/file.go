package file

import (
	"os"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

//reomove files which produced in past n days
func FilterOldLogFile(paths []string, validTime int64) (newPath []string) {
	for index := range paths {
		fileInfo, statErr := os.Stat(paths[index])
		if statErr != nil {
			logs.GetLtsLogger().Warnf("Stat file [%s] error: %s ", paths[index], statErr)
			continue
		}
		modTime := utils.GetMsFromTime(fileInfo.ModTime())
		currentTime := utils.GetCurrTSInMs()
		if currentTime-modTime <= validTime {
			newPath = append(newPath, paths[index])
		}
	}
	return
}

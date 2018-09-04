// +build !windows

package file

import (
	"os"
	"strings"
	"syscall"
	"time"
)

type StateOS struct {
	Inode  uint64 `json:"inode,"`
	Device uint64 `json:"device,"`
}

type FileState struct {
	FilePath    string      `json:"filePath"` //file path
	Info        os.FileInfo `json:"-"`        // the file info
	OffSet      uint64      `json:"offset"`
	FileStateOS StateOS     `json:"fileStateOs"`  //StateOS indentify the unique of file
	FingerPrint string      `json:"finger_print"` //file first line hash
	Finished    bool        `json:"finished"`     //false indicate the file is being collected
	Timestamp   time.Time   `json:"timestamp"`
	LineNumber  uint64      `json:"line_number"`
}

type FileStates struct {
	States []FileState
}

func (s *FileStates) FindPrevious(newState FileState) *FileState {
	for index := range s.States {
		if s.States[index].FileStateOS.IsSame(newState.FileStateOS) && strings.Compare(s.States[index].FingerPrint, newState.FingerPrint) == 0 {
			return &s.States[index]
		}
	}
	return nil
}

func (fs StateOS) IsSame(state StateOS) bool {
	return fs.Inode == state.Inode && fs.Device == state.Device
}

func GetOSState(info os.FileInfo) StateOS {
	stat := info.Sys().(*syscall.Stat_t)
	fileState := StateOS{
		Inode:  uint64(stat.Ino),
		Device: uint64(stat.Dev),
	}
	return fileState
}

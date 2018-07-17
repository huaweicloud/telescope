package file

import (
	"os"
	"reflect"
	"strings"
	"time"
)

type StateOS struct {
	IdxHi uint64 `json:"idxhi,"`
	IdxLo uint64 `json:"idxlo,"`
	Vol   uint64 `json:"vol,"`
}

type FileState struct {
	FilePath    string      `json:"filePath"` //file path
	Info        os.FileInfo `json:"-"`        // the file info
	OffSet      uint64      `json:"offset"`
	FileStateOS StateOS     `json:"fileStateOs"`  //StateOS identify the unique of file
	FingerPrint string      `json:"finger_print"` //file first line hash
	Finished    bool        `json:"finished"`     //false indicate the file is being collected
	Timestamp   time.Time   `json:"timestamp"`
	LineNumber  uint64      `json:"line_number"`
}

type FileStates struct {
	States []FileState
}

// findPreviousState returns the previous state fo the file
// In case no previous state exists, return empty
func (s *FileStates) FindPrevious(newState FileState) *FileState {
	for index := range s.States {
		// This is using the FileStateOS for comparison as FileInfo identifiers can only be fetched for existing files
		if s.States[index].FileStateOS.IsSame(newState.FileStateOS) && strings.Compare(s.States[index].FingerPrint, newState.FingerPrint) == 0 {
			return &s.States[index]
		}
	}
	return nil
}

// IsSame file checks if the files are identical
func (fs StateOS) IsSame(state StateOS) bool {
	return fs.IdxHi == state.IdxHi && fs.IdxLo == state.IdxLo && fs.Vol == state.Vol
}

// GetOSState returns the platform specific StateOS
func GetOSState(info os.FileInfo) StateOS {

	os.SameFile(info, info)

	fileStat := reflect.ValueOf(info).Elem()

	fileState := StateOS{
		IdxHi: uint64(fileStat.FieldByName("idxhi").Uint()),
		IdxLo: uint64(fileStat.FieldByName("idxlo").Uint()),
		Vol:   uint64(fileStat.FieldByName("vol").Uint()),
	}

	return fileState
}

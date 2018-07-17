package logdumper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/file"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type Recorder struct {
	path       string
	fileStates *file.FileStates
}

//initialize
func (p *Recorder) InitRecord() {
	CleanInvalidStatesInLocal()
	fileStates := GetStatesFromRecordFile()
	if len(fileStates.States) > 0 {
		for index := range fileStates.States {
			fileStates.States[index].Finished = true
		}
	}
	p.fileStates = fileStates
	logs.GetLtsLogger().Debugf("Init [%v] file states from persisted record.", len(fileStates.States))
}

//decode the states from local record file
func GetStatesFromRecordFile() *file.FileStates {
	dir, _ := os.Getwd()
	recordFilePath := dir + lts_utils.RECORD_FILE_PATH
	_, err := os.Stat(recordFilePath)
	if os.IsNotExist(err) {
		recordPath := filepath.Dir(recordFilePath)
		err := os.MkdirAll(recordPath, 0750)
		if err != nil {
			logs.GetLtsLogger().Errorf("Failed to created record file dir %s: %v", recordPath, err)
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}
		f, createError := os.Create(recordFilePath) //create file
		defer f.Close()
		if createError != nil {
			logs.GetLtsLogger().Errorf("Failed to create record file:%s, error is %s ", recordFilePath, createError.Error())
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}
		return &file.FileStates{States: []file.FileState{}}
	} else {
		openfile, openError := os.Open(recordFilePath) //open file
		defer openfile.Close()
		if openError != nil {
			logs.GetLtsLogger().Errorf("Failed to open file: %s, error is %s", recordFilePath, openError.Error())
			logs.GetLtsLogger().Flush()
			os.Exit(1)
		}

		decoder := json.NewDecoder(openfile)
		states := []file.FileState{}
		err := decoder.Decode(&states)
		//make sure the file is finished to read when agent works again
		for stateIndex := range states {
			states[stateIndex].Finished = true
		}
		if err != nil {
			return &file.FileStates{States: []file.FileState{}}
		} else {
			return &file.FileStates{States: states}
		}
	}

}

//if oldstate exist in memory,then return ,whether return {}
func (p *Recorder) FileStateExistInRecord(newFileState file.FileState) *file.FileState {
	oldStates := p.fileStates
	oldState := oldStates.FindPrevious(newFileState)
	return oldState
}

//update State in Mem
func (p *Recorder) UpdateState(fileState file.FileState) {
	if p.FileStateExistInRecord(fileState) == nil {
		p.fileStates.States = append(p.fileStates.States, fileState)
	} else {
		for index := range p.fileStates.States {
			if p.fileStates.States[index].FileStateOS.IsSame(fileState.FileStateOS) && strings.Compare(p.fileStates.States[index].FingerPrint, fileState.FingerPrint) == 0 {
				p.fileStates.States[index] = fileState
			}
		}
	}

}

//write the new states into the local Record file
func ReWriteRecordFile(fileStates *file.FileStates) error {
	dir, _ := os.Getwd()
	path := dir + lts_utils.RECORD_FILE_PATH
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to open record File : %s, error is %s.", path, err.Error())
		logs.GetLtsLogger().Flush()
		os.Exit(1)
	}
	states := fileStates.States
	encoder := json.NewEncoder(f)
	err = encoder.Encode(states)
	if err != nil {
		f.Close()
		logs.GetLtsLogger().Errorf("Error when encoding the states: %s", err)
		return err
	}
	return nil
}

//backup the record file into local
//update record file when success to send a log request to server
func (p *Recorder) UpdateRecord(fileState file.FileState, offset uint64, lineNumber uint64) {
	//update Record in memory
	fileState.Finished = true
	fileState.OffSet = offset
	fileState.LineNumber = lineNumber
	p.UpdateState(fileState)
	//update Record in local file
	PersistFileStateInLocal(fileState)
}

//backup the record file into local
func PersistFileStateInLocal(state file.FileState) {
	fileStates := GetStatesFromRecordFile()
	fileState := fileStates.FindPrevious(state)
	if fileState != nil {
		for oldStateIndex := range fileStates.States {
			if fileStates.States[oldStateIndex].FileStateOS.IsSame(state.FileStateOS) && strings.Compare(fileStates.States[oldStateIndex].FingerPrint, state.FingerPrint) == 0 {
				fileStates.States[oldStateIndex] = state
			}
		}
	} else {
		fileStates.States = append(fileStates.States, state)
	}
	//update file name for the case: log file is full and is renamed
	fileStates = UpdateFileStateFileName(fileStates, Extractors)
	ReWriteRecordFile(fileStates)
}

//Clean the invalid states in Record to upgrade the performance
func CleanInvalidStatesInLocal() {
	//get the valid state in regisrar file
	fileStates := GetStatesFromRecordFile()
	validStatesArr := make([]file.FileState, 0, len(fileStates.States))
	for fileStateIndex := range fileStates.States {
		fileInfo, err := os.Stat(fileStates.States[fileStateIndex].FilePath)
		if err != nil {
			logs.GetLtsLogger().Warnf("State file [%s] error: %s", fileStates.States[fileStateIndex].FilePath, err.Error())
			continue
		}
		if utils.GetCurrTSInMs()-utils.GetMsFromTime(fileInfo.ModTime()) <= lts_utils.LOG_File_VALID_DURATION {
			validStatesArr = append(validStatesArr, fileStates.States[fileStateIndex])
		}
	}
	validStates := file.FileStates{States: validStatesArr}
	ReWriteRecordFile(&validStates)
}

//update attribute filePath in FileState
func UpdateFileStateFileName(fileStates *file.FileStates, extractors []extractor) (newFileStates *file.FileStates) {
	//get all filestate in all extractor
	fileNames := make([]string, 0, len(fileStates.States))
	for extIndex := range extractors {
		extFiles, err := utils.GetAllFilesFromDirectoryPath(extractors[extIndex].path)
		if err == nil {
			fileNames = utils.MergeStringArr(extFiles, fileNames)
		}
	}
	allExistFileStates := make([]file.FileState, 0, len(fileNames))
	for fileIndex := range fileNames {
		fileInfo, _ := os.Stat(fileNames[fileIndex])
		finger := utils.GetFileFingerPrint(fileNames[fileIndex])
		fileState := file.FileState{FileStateOS: file.GetOSState(fileInfo), FingerPrint: finger, FilePath: fileNames[fileIndex]}
		allExistFileStates = append(allExistFileStates, fileState)
	}
	//update the file path according to the newest file name
	existFileStates := file.FileStates{States: allExistFileStates}
	for fIndex := range fileStates.States {
		existState := existFileStates.FindPrevious(fileStates.States[fIndex])
		if existState != nil && existState.FilePath != "" {
			fileStates.States[fIndex].FilePath = existState.FilePath
		}
	}
	newFileStates = fileStates
	return
}

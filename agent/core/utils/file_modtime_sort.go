package utils

import (
	"os"
	"sort"
)

type fileInfoStr []string

// Less ...
// if sort ,must implement Less,Len,Swap function
func (self fileInfoStr) Less(i, j int) bool {
	fileInfo1, _ := os.Stat(self[i])
	fileInfo2, _ := os.Stat(self[j])
	return fileInfo1.ModTime().UnixNano() <= fileInfo2.ModTime().UnixNano()
}

// Len ..
func (self fileInfoStr) Len() int {
	return len(self)
}

// Swap ...
func (self fileInfoStr) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// FileListSortTimeAsc sort the file in ascending order
func FileListSortTimeAsc(fileStrArr []string) (fileStrArrSortResult []string) {
	if len(fileStrArr) > 1 {
		sort.Sort(fileInfoStr(fileStrArr))
		for fileIndex := range fileStrArr {
			fileStrArrSortResult = append(fileStrArrSortResult, fileStrArr[fileIndex])
		}
	} else {
		fileStrArrSortResult = fileStrArr
	}

	return
}

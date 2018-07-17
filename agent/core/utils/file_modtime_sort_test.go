package utils

import (
	"path/filepath"
	"testing"
)

func TestFileListSortTimeAsc(t *testing.T) {
	var files []string
	pattern := "D:/test/*.txt"
	formatDir := filepath.FromSlash(pattern) //format “/” to path separator according to system
	matches, err := filepath.Glob(formatDir)
	if err != nil {
		t.Errorf("there is no matches under the path: %s", pattern)
	}
	if len(matches) == 0 {
		t.Logf("there is no matches under the path: %s", pattern)
	}
	for _, matchStr := range matches {
		if IsFileOrDir(matchStr, false) {
			files = append(files, matchStr)
		}
	}
	files = FileListSortTimeAsc(files)
	for _, file := range files {
		t.Log(file)
	}
}

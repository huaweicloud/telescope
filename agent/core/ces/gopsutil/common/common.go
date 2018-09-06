package common

import (
	"bufio"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"os"
	"path/filepath"
	"strings"
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
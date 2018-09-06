package utils

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"sync"

	"encoding/json"

	"bytes"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	reuse "github.com/libp2p/go-reuseport"
	"os/user"
)

var HTTPClient *http.Client
var clientMutex = new(sync.Mutex)
var CLIENT_POINT = 0

// limit float number to two decimal points
func Limit2Decimal(number float64) float64 {
	limitedString := fmt.Sprintf("%.2f", number)
	limitedNumber, _ := strconv.ParseFloat(limitedString, 64)
	return limitedNumber
}

// return os type, eg: windows_amd64
func GetOsType() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	osType := os + "_" + arch
	return osType
}

// compare two files with md5 of file
func CompareFiles(srcFile, tarFile string) bool {
	srcMd5 := GetMd5OfFile(srcFile)
	tarMd5 := GetMd5OfFile(tarFile)
	return srcMd5 == tarMd5
}

// get md5 of file
func GetMd5OfFile(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	md5 := md5.Sum(bytes)
	return string(md5[:])
}

//Persiste the struct to file
func WriteToFile(obj interface{}, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetLogger().Errorf("Failed to open file [%s], error is %s.", filepath, err.Error())
		return err
	}
	encoder := json.NewEncoder(f)
	err = encoder.Encode(obj)
	if err != nil {
		logs.GetLogger().Errorf("Failed to encoding the object: %v, error is %s", obj, err.Error())
		return err
	}
	return nil
}

func WriteStrToFile(str string, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetLogger().Errorf("Failed to open file [%s], error is %s.", filepath, err.Error())
		return err
	}
	_, err = f.WriteString(str)
	if err != nil {
		logs.GetLogger().Errorf("Failed to write string: %s, error is %s", str, err.Error())
		return err
	}
	return nil
}

// Concat the two str, str1 first, strs follows
func ConcatStr(str1 string, str2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(str1)
	buffer.WriteString(str2)

	return buffer.String()
}

// We don't use the time.Millisecond
// Because 1000000 is more accurate and efficient
func GetCurrTSInMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetMsFromTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

//according the file name and decideDirBool to judge file or directory
func IsFileOrDir(fileName string, decideDirBool bool) bool {
	fileInfo, error := os.Stat(fileName)
	if error != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDirBool {
		return isDir
	}
	return !isDir
}

//get all files according to the pattern
func getAllFileWithPattern(pattern string) []string {
	var files []string
	formatDir := filepath.FromSlash(pattern) //format “/” to path separator according to system
	matches, err := filepath.Glob(formatDir) // get all files and directory according to the pattern
	if err != nil || len(matches) == 0 {
		logs.GetLogger().Infof("there is no matches with the pattern: %s", pattern)
	}
	for _, matchStr := range matches {
		if IsFileOrDir(matchStr, false) {
			files = append(files, matchStr)
		}
	}
	return files
}

//cut out string according to the size
func SubStr(str string, size int) string {

	strRune := []rune(str)
	contentSize := 0
	strIndex := 0

	for strIndex = range strRune {
		if contentSize <= size {
			contentSize = contentSize + len(string(strRune[strIndex]))
		} else {
			break
		}
	}
	if strIndex > 0 {
		strRune = strRune[:strIndex-1]
	}
	return string(strRune)

}

//转换时间字符串为time.Time
func ParseLogTime(logTimeContent string, timePattern string, recentTime time.Time) (time.Time, error) {
	if len(timePattern) == 0 {
		return time.Now(), Errors.NoTimeFormat
	} else {
		if len(logTimeContent) >= 0 {
			logTime, err := Parse(timePattern, logTimeContent)
			if err != nil {
				if !recentTime.IsZero() {
					return recentTime, err
				}
				return time.Now(), err
			} else {
				return logTime, nil
			}
		} else {
			if !recentTime.IsZero() {
				return recentTime, Errors.NoTimeInTheLog
			}
			return time.Now(), Errors.NoTimeInTheLog
		}
	}

}

//retrieve all files in the path, eg. /var/err*,/var/log/*
func GetAllFilesFromDirectoryPath(path string) (files []string, err error) {
	pathArr := make([]string, 0)
	if strings.Contains(path, ",") {
		pathArr = strings.Split(path, ",")
	} else {
		pathArr = append(pathArr, path)
	}
	logs.GetLtsLogger().Debugf("Paths: %v", pathArr)
	var matches []string
	for i, _ := range pathArr {
		formatDir := filepath.FromSlash(pathArr[i]) //format “/” to path separator according to system
		matches, err = filepath.Glob(formatDir)
		if err != nil {
			return
		}
		if len(matches) == 0 {
			err = Errors.NoMatchedFileFound
			return
		}
		for _, matchStr := range matches {
			if IsFileOrDir(matchStr, false) && isNotCompressedFile(matchStr) && isNotTempFile(matchStr) {
				files = append(files, matchStr)
			}
		}
	}

	return
}

//is compressed file
func isNotCompressedFile(file string) bool {
	if !strings.Contains(file, ".zip") && !strings.Contains(file, ".tar") && !strings.Contains(file, ".gz") && !strings.Contains(file, ".rar") {
		return true
	}
	return false

}

func isNotTempFile(file string) bool {
	if !strings.HasPrefix(filepath.Base(file), ".") {
		return true
	}
	return false
}

//merge the oriArr to destArr
func MergeStringArr(oriArr []string, destArr []string) []string {
	for curIndex := range oriArr {
		if !StrArrContainsStr(destArr, oriArr[curIndex]) {
			destArr = append(destArr, oriArr[curIndex])
		}
	}
	return destArr

}

//An array which type is string if contain a string
func StrArrContainsStr(strArr []string, str string) bool {
	if len(strArr) > 0 {
		for strIndex := range strArr {
			if strArr[strIndex] == str {
				return true
			} else {
				continue
			}
		}
	}
	return false
}

//get local ip
func GetLocalIp() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// check whether ip address is loop or not
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return
}

//get system host name
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return name
	}
}

//check current user,if user is root ,do not start this agent
func CheckCurrentUser(args []string) bool {
	if runtime.GOOS == "linux" {
		for _, arg := range args {
			if arg == "-R" {
				return true
			}
		}
		usr, err := user.Current()
		if err != nil {
			logs.GetLogger().Errorf("Get current user error: %s", err)
			logs.GetLogger().Flush()
			return false
		}
		if usr.Username == "root" {
			logs.GetLogger().Errorf("Running as root user or group is not recommended! If you really want to force running as root, use -R command line option!")
			logs.GetLogger().Flush()
			return false
		}
	}
	return true
}

func UncompressFile(zipFilePath string, destPath string) (string, error) {
	if strings.HasSuffix(strings.ToLower(zipFilePath), ".zip") {
		return UncompressZipFile(zipFilePath, destPath)
	}

	return UncompressTgzFile(zipFilePath, destPath)
}

// uncompress .zip
func UncompressZipFile(zipFilePath string, destPath string) (string, error) {
	os.Mkdir(destPath, 0777)
	// zip 为windows升级包，需要特殊处理
	if runtime.GOOS == "windows" {
		filePath := strings.Split(zipFilePath, "/")
		zipFilePath = destPath + "/" + filePath[len(filePath)-1]
	}

	cf, err := zip.OpenReader(zipFilePath) //读取zip文件
	if err != nil {
		return "", err
	}
	defer cf.Close()
	var destDir string
	for _, file := range cf.File {
		if destDir == "" {
			destDir = file.Name
		}
		if !file.FileInfo().IsDir() {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			os.MkdirAll(destPath+"/"+path.Dir(file.Name), os.ModePerm)

			f, err := os.Create(destPath + "/" + file.Name)
			if err != nil {
				return "", err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return "", err
			}
		}
	}
	return destDir, nil
}

// uncompress .tar.gz or .tgz file
func UncompressTgzFile(tgzFilePath string, destPath string) (string, error) {

	os.Mkdir(destPath, os.ModePerm)
	fr, err := os.Open(tgzFilePath)
	if err != nil {
		logs.GetLogger().Errorf("Open tgz file failed, err:%s", err.Error())
		return "", err
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		logs.GetLogger().Errorf("Create gzip reader failed, err:%s", err.Error())
		return "", err
	}
	tr := tar.NewReader(gr)
	var destDir string
	for {
		hdr, err := tr.Next()
		if destDir == "" {
			destDir = hdr.Name
		}
		if err == io.EOF {
			break
		}
		if hdr.Typeflag != tar.TypeDir {
			os.MkdirAll(destPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			fw, err := os.Create(destPath + "/" + hdr.Name)
			if err != nil {
				logs.GetLogger().Errorf("Create new file failed, err:%s", err.Error())
				return "", err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				logs.GetLogger().Errorf("Copy reader failed, err:%s", err.Error())
				return "", err
			}
		}
	}
	return destDir, nil
}

// copy file
func CopyFile(srcFilePath, destFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return destFile.Sync()
}

// create dir
func CreateDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logs.GetLogger().Warnf("Create agent tmp dir failed, err:%s", err)
			return err
		}
	}
	return nil
}

//calculate finger print of the file
func GetFileFingerPrint(filePath string) (fingerPrint string) {
	var firstLineContent string
	h := md5.New()
	file, openFileErr := os.Open(filePath)
	if openFileErr != nil {
		logs.GetLtsLogger().Errorf("Calculate file fingerprint, can't open the file:%s", filePath)
		return
	}

	buf := bufio.NewReader(file)
	for {
		firstLineContent, _ = buf.ReadString('\n')
		break
	}
	defer file.Close()
	firstLineContent = strings.TrimSpace(firstLineContent)
	if firstLineContent != "" {
		io.WriteString(h, firstLineContent)
		fingerPrint = strings.ToLower(fmt.Sprintf("%x", h.Sum(nil)))
	}
	return

}

//get md5 hash from bytes
func GetMd5FromBytes(contentBytes []byte) string {
	h := md5.New()
	h.Write(contentBytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//IsFileExist check if file exist
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// HTTPGet get result by http
func HTTPGet(url string) ([]byte, error) {
	client := GetHttpClient()
	request, rErr := http.NewRequest("GET", url, nil)
	if rErr != nil{
		logs.GetLogger().Errorf("Create request Error:",rErr.Error())
		return []byte{}, rErr
	}
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func HTTPSend(req *http.Request, service string) (*http.Response, error) {
	client := GetHttpClient()
	authReq, _ := generateAuthHeader(req, service)
	res, err := client.Do(authReq)

	if err == nil && (res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden) {
		ChooseConfOrApiAksk(true)
	} else {
		ChooseConfOrApiAksk(false)
	}

	return res, err
}

func GetHttpClient() *http.Client{
	clientMutex.Lock()
	defer clientMutex.Unlock()

	currentPoint := GetClientPort()
	netAddr := &net.TCPAddr{Port:currentPoint}
	transport := &http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},
		DisableKeepAlives:true,
		DialContext:(&net.Dialer{
			Timeout:	10 * time.Second,
			LocalAddr:  netAddr,
		}).DialContext,
	}
	HTTPClient = &http.Client{Transport:transport, Timeout: 10 * time.Second}
	logs.GetLogger().Infof("New client had create, CLIENT_POINT: %d, GetClientPort: %d", CLIENT_POINT, currentPoint)
	CLIENT_POINT = currentPoint
	return HTTPClient
}

func ReuseDial(network, addr string)(net.Conn, error){
	dialPort := strconv.Itoa(CLIENT_POINT)
	if !reuse.Available(){
		dialPort = "0"
	}
	return reuse.Dial(network, "0.0.0.0:" + dialPort, addr)
}

func generateAuthHeader(req *http.Request, service string) (*http.Request, error) {
	nowTime := time.Unix(time.Now().Unix(), 0)
	x_sdk_date := nowTime.UTC().Format(BasicDateFormat)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set(HeaderXDate, x_sdk_date)
	req.Header.Set(HeaderProjectId, GetConfig().ProjectId)

	s := NewSigner(service)

	authkey, err := s.GetAuthorization(req)
	if err != nil {
		logs.GetLtsLogger().Errorf("Get authorization header error, error is %v", err)
	}
	req.Header.Set(HeaderAuthorization, authkey)
	if s.AKSKToken != "" {
		req.Header.Set("X-Security-Token", s.AKSKToken)
	}
	return req, err
}

func IsWindowsOs() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

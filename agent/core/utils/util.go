package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	reuse "github.com/libp2p/go-reuseport"
)

var (
	HTTPClient   *http.Client
	clientMutex  = new(sync.Mutex)
	CLIENT_POINT = 0
	workingPath  string
)

// Limit2Decimal limit float number to two decimal points
func Limit2Decimal(number float64) float64 {
	limitedString := fmt.Sprintf("%.2f", number)
	limitedNumber, _ := strconv.ParseFloat(limitedString, 64)
	return limitedNumber
}

// GetOsType return os type, eg: windows_amd64
func GetOsType() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	osType := os + "_" + arch
	return osType
}

// CompareFiles compare two files with md5 of file
func CompareFiles(srcFile, tarFile string) bool {
	srcMd5 := GetMd5OfFile(srcFile)
	tarMd5 := GetMd5OfFile(tarFile)
	return srcMd5 == tarMd5
}

// GetMd5OfFile get md5 of file
func GetMd5OfFile(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	md5 := md5.Sum(bytes)
	return string(md5[:])
}

// WriteToFile Persiste the struct to file
func WriteToFile(obj interface{}, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to open file [%s], error is %s.", filepath, err.Error())
		return err
	}
	encoder := json.NewEncoder(f)
	err = encoder.Encode(obj)
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to encoding the object: %v, error is %s", obj, err.Error())
		return err
	}
	return nil
}

// WriteStrToFile ...
func WriteStrToFile(str string, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to open file [%s], error is %s.", filepath, err.Error())
		return err
	}
	_, err = f.WriteString(str)
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to write string: %s, error is %s", str, err.Error())
		return err
	}
	return nil
}

// ConcatStr Concat the two str, str1 first, strs follows
func ConcatStr(str1 string, str2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(str1)
	buffer.WriteString(str2)

	return buffer.String()
}

// GetCurrTSInMs ...
// We don't use the time.Millisecond
// Because 1000000 is more accurate and efficient
func GetCurrTSInMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetMsFromTime ...
func GetMsFromTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// IsFileOrDir according the file name and decideDirBool to judge file or directory
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
		logs.GetCesLogger().Infof("there is no matches with the pattern: %s", pattern)
	}
	for _, matchStr := range matches {
		if IsFileOrDir(matchStr, false) {
			files = append(files, matchStr)
		}
	}
	return files
}

// SubStr cut out string according to the size
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

// MergeStringArr merge the oriArr to destArr
func MergeStringArr(oriArr []string, destArr []string) []string {
	for curIndex := range oriArr {
		if !StrArrContainsStr(destArr, oriArr[curIndex]) {
			destArr = append(destArr, oriArr[curIndex])
		}
	}
	return destArr

}

// StrArrContainsStr An array which type is string if contain a string
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

// GetLocalIp get local ip
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

// GetHostName get system host name
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return name
	}
}

// CheckCurrentUser check current user,if user is root ,do not start this agent
func CheckCurrentUser(args []string) bool {
	if runtime.GOOS == "linux" {
		for _, arg := range args {
			if arg == "-R" {
				return true
			}
		}
		usr, err := user.Current()
		if err != nil {
			logs.GetCesLogger().Errorf("Get current user error: %s", err)
			logs.GetCesLogger().Flush()
			return false
		}
		if usr.Username == "root" {
			logs.GetCesLogger().Errorf("Running as root user or group is not recommended! If you really want to force running as root, use -R command line option!")
			logs.GetCesLogger().Flush()
			return false
		}
	}
	return true
}

// UncompressFile ...
func UncompressFile(zipFilePath string, destPath string) (string, error) {
	if strings.HasSuffix(strings.ToLower(zipFilePath), ".zip") {
		return UncompressZipFile(zipFilePath, destPath)
	}

	return UncompressTgzFile(zipFilePath, destPath)
}

// UncompressZipFile uncompress .zip
func UncompressZipFile(zipFilePath string, destPath string) (string, error) {
	os.Mkdir(destPath, 0777)
	// zip 为windows升级包，需要特殊处理
	if runtime.GOOS == "windows" {
		_, pathfile := filepath.Split(zipFilePath)
		zipFilePath = filepath.Join(destPath, pathfile)
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
			os.MkdirAll(filepath.Join(destPath, path.Dir(file.Name)), os.ModePerm)

			f, err := os.Create(filepath.Join(destPath, file.Name))
			if err != nil {
				return "", err
			}
			_, err = io.Copy(f, rc)
			if err != nil {
				f.Close()
				return "", err
			}
			f.Close()
		}
	}
	return destDir, nil
}

// UncompressTgzFile uncompress .tar.gz or .tgz file
func UncompressTgzFile(tgzFilePath string, destPath string) (string, error) {

	os.Mkdir(destPath, os.ModePerm)
	fr, err := os.Open(tgzFilePath)
	if err != nil {
		logs.GetCesLogger().Errorf("Open tgz file failed, err:%s", err.Error())
		return "", err
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		logs.GetCesLogger().Errorf("Create gzip reader failed, err:%s", err.Error())
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
			os.MkdirAll(filepath.Join(destPath, path.Dir(hdr.Name)), os.ModePerm)
			fw, err := os.Create(filepath.Join(destPath, hdr.Name))
			if err != nil {
				logs.GetCesLogger().Errorf("Create new file failed, err:%s", err.Error())
				return "", err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				logs.GetCesLogger().Errorf("Copy reader failed, err:%s", err.Error())
				fw.Close()
				return "", err
			}
			fw.Close()
		}
	}
	return destDir, nil
}

// CopyFile copy file
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

// CreateDir create dir
func CreateDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logs.GetCesLogger().Warnf("Create agent tmp dir failed, err:%s", err)
			return err
		}
	}
	return nil
}

// GetMd5FromBytes get md5 hash from bytes
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
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error:", rErr.Error())
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

// HTTPSend ...
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

func GetHttpClient() *http.Client {

	clientMutex.Lock()
	defer clientMutex.Unlock()

	currentPoint := GetClientPort()
	netAddr := &net.TCPAddr{Port: currentPoint}
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			LocalAddr: netAddr,
		}).DialContext,
	}
	HTTPClient = &http.Client{Transport: transport, Timeout: 10 * time.Second}

	logs.GetCesLogger().Debugf("New client had create, CLIENT_POINT: %d , GetClientPort: %d", CLIENT_POINT, currentPoint)
	CLIENT_POINT = currentPoint
	return HTTPClient
}

// ReuseDial ...
func ReuseDial(network, addr string) (net.Conn, error) {
	dialPort := strconv.Itoa(CLIENT_POINT)
	if !reuse.Available() {
		dialPort = "0"
	}
	return reuse.Dial(network, "0.0.0.0:"+dialPort, addr)
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
		logs.GetCesLogger().Errorf("Get authorization header error, error is %v", err)
	}
	req.Header.Set(HeaderAuthorization, authkey)
	if s.AKSKToken != "" {
		req.Header.Set("X-Security-Token", s.AKSKToken)
	}
	return req, err
}

// IsWindowsOs ...
func IsWindowsOs() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func GetWorkingPath() string {
	var err error

	if workingPath == "" {
		workingPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logs.GetCesLogger().Errorf("Get working path path failed and error is: %v", err)
			return ""
		}
	}

	return workingPath
}

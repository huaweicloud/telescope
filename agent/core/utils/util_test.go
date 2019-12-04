package utils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

func TestLimit2Decimal(t *testing.T) {
	number1 := 1.3456789
	number2 := 1.3446789
	number3 := float64(12345)
	if Limit2Decimal(number1) != 1.35 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number1))
	}

	if Limit2Decimal(number2) != 1.34 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number2))
	}

	if Limit2Decimal(number3) != 12345 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number3))
	}
}

//GetOsType
func TestGetOsType(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetOsType", t, func() {
		Convey("test case 1", func() {
			osType := GetOsType()
			So(osType, ShouldNotBeBlank)
		})
	})
}

//CompareFiles
func TestCompareFiles(t *testing.T) {
	//go InitConfig()
	Convey("Test_CompareFiles", t, func() {
		Convey("test case 1", func() {
			files := CompareFiles("", "")
			So(files, ShouldBeTrue)
		})
	})
}

//GetMd5OfFile
func TestGetMd5OfFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetMd5OfFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, errors.New("")
			})
			file := GetMd5OfFile("")
			So(file, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, nil
			})
			file := GetMd5OfFile("")
			So(file, ShouldNotBeBlank)
		})
	})
}

//WriteToFile
func TestWriteToFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_WriteToFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, errors.New("123")
			})
			file := WriteToFile("", "")
			So(file.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*json.Encoder).Encode, func(j *json.Encoder, v interface{}) error {
				return errors.New("123")
			})
			file := WriteToFile("", "")
			So(file.Error(), ShouldEqual, "123")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*json.Encoder).Encode, func(j *json.Encoder, v interface{}) error {
				return nil
			})
			file := WriteToFile("", "")
			So(file, ShouldBeNil)
		})
	})
}

//WriteStrToFile
func TestWriteStrToFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_WriteStrToFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, errors.New("123")
			})
			file := WriteStrToFile("", "")
			So(file.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*os.File).WriteString, func(o *os.File, s string) (n int, err error) {
				return 0, errors.New("123")
			})
			file := WriteStrToFile("", "")
			So(file.Error(), ShouldEqual, "123")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace((*os.File).WriteString, func(o *os.File, s string) (n int, err error) {
				return 0, nil
			})
			file := WriteStrToFile("", "")
			So(file, ShouldBeNil)
		})
	})
}

//ConcatStr
func TestConcatStr(t *testing.T) {
	//go InitConfig()
	Convey("Test_ConcatStr", t, func() {
		Convey("test case 1", func() {
			str := ConcatStr("1", "2")
			So(str, ShouldEqual, "12")
		})
	})
}

//GetCurrTSInMs
func TestGetCurrTSInMs(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetCurrTSInMs", t, func() {
		Convey("test case 1", func() {
			ms := GetCurrTSInMs()
			So(ms, ShouldNotBeNil)
		})
	})
}

//GetMsFromTime
func TestGetMsFromTime(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetMsFromTime", t, func() {
		Convey("test case 1", func() {
			ms := GetMsFromTime(time.Now())
			So(ms, ShouldNotBeNil)
		})
	})
}

//IsFileOrDir
func TestIsFileOrDir(t *testing.T) {
	//go InitConfig()
	Convey("Test_IsFileOrDir", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				return nil, errors.New("123")
			})
			dir := IsFileOrDir("a", true)
			So(dir, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				dir, _ := os.Getwd()
				path := dir + "/util_test.go"
				info, _ := os.Lstat(path)
				return info, nil
			})
			dir := IsFileOrDir("a", true)
			So(dir, ShouldBeFalse)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				dir, _ := os.Getwd()
				path := dir + "/util_test.go"
				info, _ := os.Lstat(path)
				return info, nil
			})
			dir := IsFileOrDir("a", false)
			So(dir, ShouldBeTrue)
		})
	})
}

//getAllFileWithPattern
func TestGetAllFileWithPattern(t *testing.T) {
	//go InitConfig()
	Convey("Test_getAllFileWithPattern", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(filepath.Glob, func(pattern string) (matches []string, err error) {
				return []string{"1", "2"}, errors.New("123")
			})
			//IsFileOrDir
			mockfn.Replace(IsFileOrDir, func(fileName string, decideDirBool bool) bool {
				return true
			})
			pattern := getAllFileWithPattern("a")
			So(pattern, ShouldNotBeNil)
		})
	})
}

//SubStr
func TestSubStr(t *testing.T) {
	//go InitConfig()
	Convey("Test_SubStr", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
			str := SubStr("123", 0)
			So(str, ShouldBeBlank)
		})
	})
}

//isNotCompressedFile
func TestIsNotCompressedFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_isNotCompressedFile", t, func() {
		Convey("test case 1", func() {
			file := isNotCompressedFile("a.zip")
			So(file, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			file := isNotCompressedFile("a")
			So(file, ShouldBeTrue)
		})
	})
}

//isNotTempFile
func TestIsNotTempFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_isNotTempFile", t, func() {
		Convey("test case 1", func() {
			file := isNotTempFile(".a")
			So(file, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			file := isNotTempFile("a")
			So(file, ShouldBeTrue)
		})
	})
}

//MergeStringArr
func TestMergeStringArr(t *testing.T) {
	//go InitConfig()
	Convey("Test_MergeStringArr", t, func() {
		Convey("test case 1", func() {
			arr := MergeStringArr([]string{"a"}, []string{"b"})
			So(arr, ShouldNotBeEmpty)
		})
	})
}

//StrArrContainsStr
func TestStrArrContainsStr(t *testing.T) {
	//go InitConfig()
	Convey("Test_StrArrContainsStr", t, func() {
		Convey("test case 1", func() {
			str := StrArrContainsStr([]string{"a"}, "a")
			So(str, ShouldBeTrue)
		})
		Convey("test case 2", func() {
			str := StrArrContainsStr([]string{"a"}, "b")
			So(str, ShouldBeFalse)
		})
	})
}

//GetLocalIp
func TestGetLocalIp(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetLocalIp", t, func() {
		Convey("test case 1", func() {
			ip := GetLocalIp()
			So(ip, ShouldNotBeBlank)
		})
	})
}

//GetHostName
func TestGetHostName(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetHostName", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Hostname, func() (name string, err error) {
				return "", errors.New("123")
			})
			name := GetHostName()
			So(name, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Hostname, func() (name string, err error) {
				return "", nil
			})
			name := GetHostName()
			So(name, ShouldBeBlank)
		})
	})
}

//CheckCurrentUser
func TestCheckCurrentUser(t *testing.T) {
	//go InitConfig()
	Convey("Test_CheckCurrentUser", t, func() {
		Convey("test case 1", func() {
		})
	})
}

//UncompressFile
func TestUncompressFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_UncompressFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(UncompressZipFile, func(zipFilePath string, destPath string) (string, error) {
				return "", nil
			})
			s, e := UncompressFile("a.zip", "b")
			So(s, ShouldBeBlank)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(UncompressTgzFile, func(zipFilePath string, destPath string) (string, error) {
				return "", nil
			})
			s, e := UncompressFile("a", "b")
			So(s, ShouldBeBlank)
			So(e, ShouldBeNil)
		})
	})
}

//UncompressZipFile
func TestUncompressZipFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_UncompressZipFile", t, func() {
		/*Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.MkdirAll, func(path string, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(zip.OpenReader, func(name string) (*zip.ReadCloser, error) {
				return nil, errors.New("123")
			})
			s, e := UncompressZipFile("z", "d1")
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.MkdirAll, func(path string, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(zip.OpenReader, func(name string) (*zip.ReadCloser, error) {
				closer := &zip.ReadCloser{}
				closer.File = []*zip.File{{}}
				return closer, nil
			})
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return &os.File{}, errors.New("123")
			})
			s, e := UncompressZipFile("z", "d2")
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(zip.OpenReader, func(name string) (*zip.ReadCloser, error) {
				closer := &zip.ReadCloser{}
				return closer, nil
			})
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.MkdirAll, func(path string, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, errors.New("12")
			})
			s, e := UncompressZipFile("z", "d2")
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(zip.OpenReader, func(name string) (*zip.ReadCloser, error) {
				closer := &zip.ReadCloser{}
				return closer, nil
			})
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.MkdirAll, func(path string, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, nil
			})
			//io.Copy
			mockfn.Replace(io.Copy, func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, errors.New("12")
			})
			s, e := UncompressZipFile("z", "d2")
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(zip.OpenReader, func(name string) (*zip.ReadCloser, error) {
				closer := &zip.ReadCloser{}
				return closer, nil
			})
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.MkdirAll, func(path string, perm os.FileMode) error {
				return nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, nil
			})
			//io.Copy
			mockfn.Replace(io.Copy, func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			})
			s, e := UncompressZipFile("z", "d2")
			So(s, ShouldBeBlank)
			So(e, ShouldBeNil)
		})*/
	})
}

//UncompressTgzFile
func TestUncompressTgzFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_UncompressTgzFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})

		})
	})
}

//CopyFile
func TestCopyFile(t *testing.T) {
	//go InitConfig()
	Convey("Test_CopyFile", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, errors.New("")
			})
			file := CopyFile("a", "b")
			So(file, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, errors.New("12")
			})
			file := CopyFile("a", "b")
			So(file, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(io.Copy, func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, errors.New("12")
			})
			file := CopyFile("a", "b")
			So(file, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(os.Create, func(name string) (*os.File, error) {
				return nil, nil
			})
			mockfn.Replace(io.Copy, func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			})
			file := CopyFile("a", "b")
			So(file, ShouldNotBeNil)
		})
	})
}

//CreateDir
func TestCreateDir(t *testing.T) {
	//go InitConfig()
	Convey("Test_CreateDir", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				return nil, nil
			})
			dir := CreateDir("")
			So(dir, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(os.Stat, func(name string) (os.FileInfo, error) {
				return nil, errors.New("")
			})
			//os.Mkdir
			mockfn.Replace(os.Mkdir, func(name string, perm os.FileMode) error {
				return errors.New("")
			})
			dir := CreateDir("")
			So(dir, ShouldNotBeNil)
		})
	})
}

//GetMd5FromBytes
func TestGetMd5FromBytes(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetMd5FromBytes", t, func() {
		Convey("test case 1", func() {
			bytes := GetMd5FromBytes([]byte(""))
			So(bytes, ShouldNotBeBlank)
		})
	})
}

//IsFileExist
func TestIsFileExist(t *testing.T) {
	//go InitConfig()
	Convey("Test_IsFileExist", t, func() {
		Convey("test case 1", func() {
			exist := IsFileExist("")
			So(exist, ShouldBeFalse)
		})
	})
}

//HTTPGet
func TestHTTPGet(t *testing.T) {
	//go InitConfig()
	Convey("Test_HTTPGet", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			bytes, e := HTTPGet("")
			So(bytes, ShouldBeEmpty)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return &http.Request{}, nil
			})
			bytes, e := HTTPGet("")
			So(bytes, ShouldBeEmpty)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			mockfn.Replace((*http.Client).Do, func(h *http.Client, req *http.Request) (*http.Response, error) {
				return nil, errors.New("123")
			})
			bytes, e := HTTPGet("")
			So(bytes, ShouldBeEmpty)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return &http.Request{}, nil
			})
			mockfn.Replace((*http.Client).Do, func(h *http.Client, req *http.Request) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, errors.New("123")
			})
			bytes, e := HTTPGet("")
			So(bytes, ShouldBeEmpty)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			mockfn.Replace((*http.Client).Do, func(h *http.Client, req *http.Request) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, nil
			})
			bytes, e := HTTPGet("")
			So(bytes, ShouldBeEmpty)
			So(e, ShouldBeNil)
		})
	})
}

//HTTPSend
func TestHTTPSend(t *testing.T) {
	//go InitConfig()
	Convey("Test_HTTPSend", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			//generateAuthHeader
			mockfn.Replace(generateAuthHeader, func(req *http.Request, service string) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace((*http.Client).Do, func(h *http.Client, req *http.Request) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			response, e := HTTPSend(nil, "")
			So(response, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetHttpClient, func() *http.Client {
				return &http.Client{}
			})
			//generateAuthHeader
			mockfn.Replace(generateAuthHeader, func(req *http.Request, service string) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace((*http.Client).Do, func(h *http.Client, req *http.Request) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, errors.New("123")
			})
			response, e := HTTPSend(nil, "")
			So(response, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
	})
}

//GetHttpClient
func TestGetHttpClient(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetHttpClient", t, func() {
		Convey("test case 1", func() {
			client := GetHttpClient()
			So(client, ShouldNotBeNil)
		})
	})
}

//ReuseDial
func TestReuseDial(t *testing.T) {
	//go InitConfig()
	Convey("TestReuseDial", t, func() {
		Convey("test case 1", func() {
			conn, e := ReuseDial("", "")
			So(conn, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
	})
}

//generateAuthHeader
func TestGenerateAuthHeader(t *testing.T) {
	//go InitConfig()
	Convey("Test_generateAuthHeader", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).GetAuthorization, func(s *Signer, req *http.Request) (string, error) {
				return "", errors.New("123")
			})
			//GetConfig
			mockfn.Replace(GetConfig, func() *GeneralConfig {
				return &GeneralConfig{}
			})
			request := &http.Request{
				Header: map[string][]string{},
			}
			header, e := generateAuthHeader(request, "")
			So(header, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
	})
}

//IsWindowsOs
func TestIsWindowsOs(t *testing.T) {
	//go InitConfig()
	Convey("Test_IsWindowsOs", t, func() {
		Convey("test case 1", func() {
			windowsOs := IsWindowsOs()
			So(windowsOs, ShouldBeTrue)
		})
	})
}

package heartbeat

import (
	"errors"
	"github.com/huaweicloud/telescope/agent/core/assistant/utils"
	utl "github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"testing"
	"time"
)

//InitConfig
func TestInitConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitConfig", t, func() {
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
		})
	})
}

//SendHBTicker
func TestSendHBTicker(t *testing.T) {
	Convey("Test_SendHBTicker", t, func() {
		Convey("test case 1", func() {
			//time.NewTicker
			defer mockfn.RevertAll()
			mockfn.Replace(time.NewTicker, func(d time.Duration) *time.Ticker {
				t.SkipNow()
				return nil
			})
			bools := make(chan bool, 1)
			SendHBTicker(bools)
		})
	})
}

//SendHBExec
func TestSendHBExec(t *testing.T) {
	Convey("Test_SendHBExec", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getHB, func() *Heartbeat {
				return nil
			})
			mockfn.Replace(sendHB, func(requestBody *Heartbeat) (error, *HBResp) {
				resp := &HBResp{
					Config: &HBConfig{},
				}
				return nil, resp
			})
			bools := make(chan bool, 1)
			SendHBExec(bools)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getHB, func() *Heartbeat {
				return nil
			})
			mockfn.Replace(sendHB, func(requestBody *Heartbeat) (error, *HBResp) {
				resp := &HBResp{
					Config: &HBConfig{AssistSwitch: true},
				}
				return nil, resp
			})
			bools := make(chan bool, 1)
			SendHBExec(bools)
		})
	})
}

//updateHBState
func TestUpdateHBState(t *testing.T) {
	Convey("Test_updateHBState", t, func() {
		Convey("test case 1", func() {
			updateHBState(nil)
		})
	})
}

//sendHB
func TestSendHB(t *testing.T) {
	Convey("Test_sendHB", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			/*	mockfn.Replace(logs.GetAssistantLogger, func() seelog.LoggerInterface {
					return seelog.Disabled
				})
				mockfn.Replace(logs.GetAssistantLogger().Debugf, func(format string, params ...interface{}) {
					return
				})*/
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, errors.New("123")
			})
			heartbeat := &Heartbeat{}
			e, resp := sendHB(heartbeat)
			So(e, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			//http.NewRequest
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			heartbeat := &Heartbeat{}
			e, resp := sendHB(heartbeat)
			So(e, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			//http.NewRequest
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utl.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				return nil, errors.New("123")
			})
			heartbeat := &Heartbeat{}
			e, resp := sendHB(heartbeat)
			So(e, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			//http.NewRequest
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utl.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			heartbeat := &Heartbeat{}
			e, resp := sendHB(heartbeat)
			So(e, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			//http.NewRequest
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utl.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			//ioutil.ReadAll
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return nil, errors.New("123")
			})
			heartbeat := &Heartbeat{}
			e, resp := sendHB(heartbeat)
			So(e, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
		Convey("test case 6", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.BuildURL, func(destURI string) string {
				return ""
			})
			//GetMarshalledRequestBody
			mockfn.Replace(utils.GetMarshalledRequestBody, func(v interface{}, url string) ([]byte, error) {
				return nil, nil
			})
			//http.NewRequest
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utl.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 204,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			//ioutil.ReadAll
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), nil
			})
		})
	})
}

//getHB
func TestGetHB(t *testing.T) {
	Convey("Test_getHB", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace(utl.GetConfig, func() *utl.GeneralConfig {
				config := &utl.GeneralConfig{}
				return config
			})
			//getCurrentUser
			mockfn.Replace(getCurrentUser, func() string {
				return ""
			})
			mockfn.Replace(getEnv, func() *Env {
				return nil
			})
			hb := getHB()
			So(hb, ShouldNotBeNil)
		})
	})
}

//getEnv
func TestGetEnv(t *testing.T) {
	Convey("Test_getEnv", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(enableEnvReport, func() bool {
				return true
			})
			hb := getEnv()
			So(hb, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace(enableEnvReport, func() bool {
				return false
			})
			hb := getEnv()
			So(hb, ShouldBeNil)
		})
	})
}

//getCurrentUser
func TestGetCurrentUser(t *testing.T) {
	Convey("Test_getEnv", t, func() {
		Convey("test case 1", func() {
			user := getCurrentUser()
			So(user, ShouldNotBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace(user.Current, func() (*user.User, error) {
				return nil, errors.New("123")
			})
			user := getCurrentUser()
			So(user, ShouldBeBlank)
		})
	})
}

//getHostname
func TestGetHostname(t *testing.T) {
	Convey("Test_getHostname", t, func() {
		Convey("test case 1", func() {
			hostname := getHostname()
			So(hostname, ShouldNotBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace(os.Hostname, func() (name string, err error) {
				return "", errors.New("123")
			})
			hostname := getHostname()
			So(hostname, ShouldBeBlank)
		})
	})
}

//getFirstIP
func TestGetFirstIP(t *testing.T) {
	/*ip := getFirstIP()
	t.Log("ip is:", ip)*/
	Convey("Test_getFirstIP", t, func() {
		Convey("test case 1", func() {
			ip := getFirstIP()
			So(ip, ShouldNotBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace(net.Interfaces, func() ([]net.Interface, error) {
				return nil, errors.New("123")
			})
			ip := getFirstIP()
			So(ip, ShouldEqual, "<nil>")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace((*net.Interface).Addrs, func(*net.Interface) ([]net.Addr, error) {
				return nil, errors.New("123")
			})
			ip := getFirstIP()
			So(ip, ShouldEqual, "<nil>")
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			//utils.GetConfig
			mockfn.Replace((*net.Interface).Addrs, func(*net.Interface) ([]net.Addr, error) {
				addrs := []net.Addr{}
				addr := &net.IPAddr{}
				addrs = append(addrs, addr)
				return addrs, nil
			})
			ip := getFirstIP()
			So(ip, ShouldEqual, "<nil>")
		})
	})
}

//enableEnvReport
func TestEnableEnvReport(t *testing.T) {
	Convey("Test_enableEnvReport", t, func() {
		Convey("test case 1", func() {
			HBS.IsEnvReported = true
			report := enableEnvReport()
			So(report, ShouldBeFalse)
		})
	})
}

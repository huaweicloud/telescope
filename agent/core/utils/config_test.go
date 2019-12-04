package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/huaweicloud/telescope/agent/core/logs"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

var confData = `{
    "InstanceId":"",
    "ProjectId": "",
    "AccessKey": "",
    "SecretKey": "",
    "RegionId": "cn-north-1",
    "ClientPort": 0,
    "PortNum": 200
}`

//ReadConfig
func TestReadConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_ReadConfig", t, func() {
		Convey("test case 1", func() {
			//os.Open
			defer mockfn.RevertAll()
			mockfn.Replace(os.Open, func(name string) (*os.File, error) {
				return nil, errors.New("123")
			})
			var i, j int = 0, 0
			config, e := ReadConfig(&i, &j)
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*json.Decoder).Decode, func(j *json.Decoder, v interface{}) error {
				return errors.New("123")
			})
			var i, j int = 0, 0
			config, e := ReadConfig(&i, &j)
			So(config, ShouldBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			//logs.GetCurrentDirectory()
			defer mockfn.RevertAll()
			mockfn.Replace(logs.GetCurrentDirectory, func() string {
				pwd, _ := os.Getwd()
				return pwd
			})
			//getConfFromOpenstack
			mockfn.Replace(getConfFromOpenstack, func() (*MetaData, error) {
				data := MetaData{ProjectId: "123", ExternalService: "AGT.ECS123"}
				return &data, errors.New("")
			})
			//getAKSKFromOpenStack
			mockfn.Replace(getAKSKFromOpenStack, func() (bool, error) {
				return false, errors.New("")
			})
			pwd, _ := os.Getwd()
			if err := ioutil.WriteFile(pwd+"/conf.json", []byte(confData), 0666); err != nil {
				t.Fatal(err)
			}
			var i, j int = 0, 0
			config, e := ReadConfig(&i, &j)
			So(config, ShouldNotBeEmpty)
			So(e, ShouldBeNil)
			os.Remove(pwd + "/conf.json")
		})

	})
}

//InitConfig
func TestInitConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitConfig", t, func() {
		Convey("test case 2", func() {
			//os.Open
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*GeneralConfig, error) {
				return nil, errors.New("123")
			})
			//time.Sleep
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				t.SkipNow()
			})
			InitConfig()
		})
		Convey("test case 1", func() {
			InitConfig()
		})
	})
}

//GetConfig
func TestGetConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetConfig", t, func() {
		Convey("test case 2", func() {
			//InitConfig
			defer mockfn.RevertAll()
			mockfn.Replace(InitConfig, func() {

			})
			//setConfAKSK
			mockfn.Replace(setConfAKSK, func(config *GeneralConfig) {

			})
			config := GetConfig()
			So(config, ShouldBeNil)
		})
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(InitConfig, func() {

			})
			mockfn.Replace(setConfAKSK, func(config *GeneralConfig) {

			})
			config = nil
			config := GetConfig()
			So(config, ShouldBeNil)
		})
	})
}

//setConfAKSK
func TestSetConfAksk(t *testing.T) {
	//go InitConfig()
	Convey("Test_setConfAksk", t, func() {
		Convey("test case 1", func() {
			buse_api_aksk = true
			generalConfig := GeneralConfig{}
			setConfAKSK(&generalConfig)
		})
		Convey("test case 2", func() {
			buse_api_aksk = false
			generalConfig := GeneralConfig{}
			setConfAKSK(&generalConfig)
		})
	})
}

//ChooseConfOrApiAksk
func TestChooseConfOrApiAksk(t *testing.T) {
	//go InitConfig()
	Convey("Test_ChooseConfOrApiAksk", t, func() {
		Convey("test case 1", func() {
			securityDataFromConf.access = ""
			ChooseConfOrApiAksk(false)
		})
		Convey("test case 2", func() {
			securityDataFromConf.access = "123"
			now_method_count = 5
			retry_count = 1
			ChooseConfOrApiAksk(true)
		})
		Convey("test case 3", func() {
			securityDataFromConf.access = "123"
			ChooseConfOrApiAksk(false)
		})
	})
}

//ReloadConfig
func TestReloadConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_ReloadConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*GeneralConfig, error) {
				config := &GeneralConfig{}
				return config, nil
			})
			config = &GeneralConfig{}
			var i, j int = 0, 0
			reloadConfig := ReloadConfig(&i, &j)
			So(reloadConfig, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ReadConfig, func() (*GeneralConfig, error) {
				return nil, errors.New("123")
			})
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				t.SkipNow()
				return
			})
			config = &GeneralConfig{}
			var i, j int = 0, 0
			reloadConfig := ReloadConfig(&i, &j)
			So(reloadConfig, ShouldNotBeNil)
		})
	})
}

//getConfFromOpenstack
func TestGetConfFromOpenstack(t *testing.T) {
	//go InitConfig()
	Convey("Test_getConfFromOpenstack", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return nil, errors.New("123")
			})
			data, e := getConfFromOpenstack()
			So(data, ShouldBeNil)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`"uuid":"123"`), nil
			})
			data, e := getConfFromOpenstack()
			So(data, ShouldBeNil)
			So(e.Error(), ShouldEqual, "Key path not found")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123"}`), nil
			})
			data, e := getConfFromOpenstack()
			So(data, ShouldBeNil)
			So(e.Error(), ShouldEqual, "Key path not found")
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), nil
			})
			//jsonparser.Get
			mockfn.Replace(jsonparser.Get, func(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), 1, 0, errors.New("123")
			})
			data, e := getConfFromOpenstack()
			So(data, ShouldBeNil)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), nil
			})
			//jsonparser.Get
			mockfn.Replace(jsonparser.Get, func(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
				return []byte(`physical`), 1, 0, nil
			})
			data, e := getConfFromOpenstack()
			So(data, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
	})
}

//getAKSKFromOpenStack
func TestGetAKSKFromOpenStack(t *testing.T) {
	//go InitConfig()
	Convey("Test_getAKSKFromOpenStack", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isNeedFreshAKSK, func(aksk_data SecurityData) bool {
				return false
			})
			b, e := getAKSKFromOpenStack()
			So(b, ShouldBeFalse)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isNeedFreshAKSK, func(aksk_data SecurityData) bool {
				return true
			})
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), nil
			})
			//isOpenStackAKSKJsonValid
			mockfn.Replace(isOpenStackAKSKJsonValid, func(strAkskData string) bool {
				return true
			})
			retry_count = 1
			b, e := getAKSKFromOpenStack()
			So(b, ShouldBeFalse)
			So(e, ShouldBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isNeedFreshAKSK, func(aksk_data SecurityData) bool {
				return true
			})
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), errors.New("123")
			})
			//isOpenStackAKSKJsonValid
			mockfn.Replace(isOpenStackAKSKJsonValid, func(strAkskData string) bool {
				return true
			})
			retry_count = 1
			b, e := getAKSKFromOpenStack()
			So(b, ShouldBeFalse)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isNeedFreshAKSK, func(aksk_data SecurityData) bool {
				return true
			})
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), nil
			})
			//isOpenStackAKSKJsonValid
			mockfn.Replace(isOpenStackAKSKJsonValid, func(strAkskData string) bool {
				return true
			})
			mockfn.Replace(parseSecurityToken, func(strAkskData string) (SecurityData, error) {
				return SecurityData{}, nil
			})
			retry_count = 1
			b, e := getAKSKFromOpenStack()
			So(b, ShouldBeTrue)
			So(e, ShouldBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(isNeedFreshAKSK, func(aksk_data SecurityData) bool {
				return true
			})
			mockfn.Replace(HTTPGet, func(url string) ([]byte, error) {
				return []byte(`{"uuid":"123","project_id":"123"}`), errors.New("123")
			})
			//isOpenStackAKSKJsonValid
			mockfn.Replace(isOpenStackAKSKJsonValid, func(strAkskData string) bool {
				return true
			})
			retry_count = 1
			b, e := getAKSKFromOpenStack()
			So(b, ShouldBeFalse)
			So(e, ShouldNotBeNil)
		})
	})
}

//parseSecurityToken
func TestParseSecurityToken(t *testing.T) {
	//go InitConfig()
	Convey("Test_parseSecurityToken", t, func() {
		Convey("test case 1", func() {
			strAkskData := ""
			data, e := parseSecurityToken(strAkskData)
			So(data, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			strAkskData := `{"credential":{"access":"a"}}`
			data, e := parseSecurityToken(strAkskData)
			So(data, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			strAkskData := `{"credential":{"access":"a","secret":"s"}}`
			data, e := parseSecurityToken(strAkskData)
			So(data, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			strAkskData := `{"credential":{"access":"a","secret":"s","securitytoken":"se"}}`
			data, e := parseSecurityToken(strAkskData)
			So(data, ShouldNotBeNil)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 5", func() {
			strAkskData := `{"credential":{"access":"a","secret":"s","securitytoken":"se","expires_at":"ex"}}`
			data, e := parseSecurityToken(strAkskData)
			So(data, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
	})
}

//isNeedFreshAKSK
func TestIsNeedFreshAKSK(t *testing.T) {
	//go InitConfig()
	Convey("Test_isNeedFreshAKSK", t, func() {
		Convey("test case 1", func() {
			data := SecurityData{}
			buse_api_aksk = true
			now_method_count = 5
			retry_count = 1
			aksk := isNeedFreshAKSK(data)
			So(aksk, ShouldBeTrue)
		})
		Convey("test case 2", func() {
			data := SecurityData{}
			buse_api_aksk = false
			now_method_count = 5
			retry_count = 1
			aksk := isNeedFreshAKSK(data)
			So(aksk, ShouldBeFalse)
		})
		Convey("test case 3", func() {
			data := SecurityData{}
			buse_api_aksk = true
			now_method_count = 5
			retry_count = 8
			aksk := isNeedFreshAKSK(data)
			So(aksk, ShouldBeTrue)
		})
		Convey("test case 4", func() {
			//time.Parse
			defer mockfn.RevertAll()
			mockfn.Replace(time.Parse, func(layout, value string) (time.Time, error) {
				return time.Now(), errors.New("123")
			})
			data := SecurityData{access: "123"}
			buse_api_aksk = true
			now_method_count = 5
			retry_count = 8
			aksk := isNeedFreshAKSK(data)
			So(aksk, ShouldBeFalse)
		})
		Convey("test case 5", func() {
			//time.Parse
			defer mockfn.RevertAll()
			mockfn.Replace(time.Parse, func(layout, value string) (time.Time, error) {
				return time.Time{}, nil
			})
			data := SecurityData{access: "123"}
			buse_api_aksk = true
			now_method_count = 5
			retry_count = 8
			aksk := isNeedFreshAKSK(data)
			So(aksk, ShouldBeTrue)
		})
	})
}

//isOpenStackAKSKJsonValid
func TestIsOpenStackAKSKJsonValid(t *testing.T) {
	//go InitConfig()
	Convey("Test_isOpenStackAKSKJsonValid", t, func() {
		Convey("test case 1", func() {
			valid := isOpenStackAKSKJsonValid("")
			So(valid, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			valid := isOpenStackAKSKJsonValid("expires_ataccesssecretsecuritytoken")
			So(valid, ShouldBeTrue)
		})
	})
}

//GetClientPort
func TestGetClientPort(t *testing.T) {
	//go InitConfig()
	Convey("TestGetClientPort", t, func() {
		Convey("test case 1", func() {
			configClientPort = 1
			port := GetClientPort()
			So(port, ShouldEqual, 0)
		})
		Convey("test case 2", func() {
			configClientPort = 1025
			configClientPort = 65534
			configPortNum = -1
			GetClientPort()
		})
	})
}

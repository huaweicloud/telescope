package logs

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//findLogfileName
func TestFindLogfileName(t *testing.T) {

	Convey("Test_findLogfileName", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(regexp.Compile, func(expr string) (*regexp.Regexp, error) {
				return nil, errors.New("")
			})
			name := findLogfileName("123")
			So(name, ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			name := findLogfileName("filename=\"aa.log\"")
			So(name, ShouldNotBeBlank)
		})
		Convey("test case 3", func() {
			name := findLogfileName("filename=\"///.log\"")
			So(name, ShouldNotBeBlank)
		})
	})
}

//getCommonLog
func TestGetCommonLog(t *testing.T) {
	//go InitConfig()
	Convey("Test_getCommonLog", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = false
			config := getCommonLog()
			So(config, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = true
			config := getCommonLog()
			So(config, ShouldBeBlank)
		})
	})
}

//getCesLog
func TestGetCesLog(t *testing.T) {
	Convey("Test_getCesLog", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = false
			config := getCesLog()
			So(config, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = true
			config := getCesLog()
			So(config, ShouldBeBlank)
		})
	})
}

//getAssistantLog
func TestGetAssistantLog(t *testing.T) {
	Convey("Test_getAssistantLog", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = false
			config := getAssistantLog()
			So(config, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(LoadConfig, func() {
				return
			})
			isLoaded = true
			logsConfig.AssistantConfig = ""
			config := getAssistantLog()
			So(config, ShouldBeBlank)
		})
	})
}

var logsData = `<?xml version="1.0" encoding="UTF-8"?>
<root>
    <common>
        <![CDATA[
			<seelog minlevel="info">
		 		<outputs formatid="common">
        			<rollingfile type="size" filename="../log/common.log" maxsize="20000000" maxrolls="5"/>
    			</outputs>
    			<formats>
        			<format id="common" format="%Date/%Time [%LEV] [%File:%Line] %Msg%r%n" />
    			</formats>
			</seelog>
		]]>
    </common>
    <ces>
        <![CDATA[
			<seelog minlevel="info">
				<outputs formatid="ces">
					<rollingfile type="size" filename="../log/ces.log" maxsize="20000000" maxrolls="5"/>
				</outputs>
				<formats>
					<format id="ces" format="%Date/%Time [%LEV] [%File:%Line] %Msg%r%n" />
				</formats>
			</seelog>
		]]>
    </ces>
    <assistant>
        <![CDATA[
			<seelog minlevel="info">
				<outputs formatid="assistant">
					<rollingfile type="size" filename="../log/assistant.log" maxsize="20000000" maxrolls="5"/>
				</outputs>
				<formats>
					<format id="assistant" format="%Date/%Time [%LEV] [%File:%Line] %Msg%r%n" />
				</formats>
			</seelog>
		]]>
    </assistant>
</root>`

//LoadConfig
func TestLoadConfig(t *testing.T) {
	LoadConfig()
	pwd := GetCurrentDirectory()
	if err := ioutil.WriteFile(pwd+"/logs_config.xml", []byte(logsData), 0666); err != nil {
		t.Fatal(err)
	}
	LoadConfig()
	defer os.Remove(pwd + "/logs_config.xml")

	Convey("Test_LoadConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetCurrentDirectory, func() string {
				return ""
			})
			//ioutil.ReadFile
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, errors.New("")
			})
			LoadConfig()
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetCurrentDirectory, func() string {
				return ""
			})
			//ioutil.ReadFile
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, nil
			})
			//xml.Unmarshal
			mockfn.Replace(xml.Unmarshal, func(data []byte, v interface{}) error {
				return errors.New("123")
			})
			LoadConfig()
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(GetCurrentDirectory, func() string {
				return ""
			})
			//ioutil.ReadFile
			mockfn.Replace(ioutil.ReadFile, func(filename string) ([]byte, error) {
				return nil, nil
			})
			//xml.Unmarshal
			mockfn.Replace(xml.Unmarshal, func(data []byte, v interface{}) error {
				return nil
			})
			LoadConfig()
		})
	})
}

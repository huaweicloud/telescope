package heartbeat

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/cihub/seelog"
	cesConfig "github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
)

//LoadHbServicesDetails
func TestLoadHbServicesDetails(t *testing.T) {
	//go InitConfig()
	Convey("Test_LoadHbServicesDetails", t, func() {
		Convey("test case 3", func() {
			data := make(chan channel.HBServiceData, 1)
			data <- channel.HBServiceData{Service: "CES"}
			beat := &HeartBeat{}
			go beat.LoadHbServicesDetails(data)
		})
	})
}

//ProduceHeartBeat
func TestProduceHeartBeat(t *testing.T) {
	//go InitConfig()
	Convey("Test_ProduceHeartBeat", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.NewTicker, func(d time.Duration) *time.Ticker {
				c := make(chan time.Time, 10)
				c <- time.Time{}
				t := &time.Ticker{
					C: c,
				}
				return t
			})
			//channel.NewHBEntity
			mockfn.Replace(channel.NewHBEntity, func(status channel.StatusEnum, time int64, cesDetails string) *channel.HBEntity {
				return nil
			})
			mockfn.Replace(utils.ReloadConfig, func() *utils.GeneralConfig {
				return nil
			})
			mockfn.Replace(cesConfig.ReloadConfig, func() *cesConfig.CESConfig {
				t.SkipNow()
				return nil
			})
			beat := &HeartBeat{}
			entities := make(chan *channel.HBEntity, 1)
			beat.ProduceHeartBeat(entities)
		})
	})
}

//ConsumeHeartBeat
func TestConsumeHeartBeat(t *testing.T) {
	//go InitConfig()
	Convey("Test_ConsumeHeartBeat", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(buildHeartBeatUrl, func(uri string) string {
				return ""
			})
			mockfn.Replace(sendHeartBeat, func(url string, hb *channel.HBEntity) *channel.HBResponse {
				return nil
			})
			mockfn.Replace(logs.GetCesLogger, func() seelog.LoggerInterface {
				t.SkipNow()
				return nil
			})
			mockfn.Replace(seelog.LoggerInterface.Errorf, func(s seelog.LoggerInterface, format string, params ...interface{}) error {
				t.SkipNow()
				return nil
			})
			beat := &HeartBeat{}
			entities := make(chan *channel.HBEntity, 1)
			entities <- &channel.HBEntity{}
			go beat.ConsumeHeartBeat(entities)
		})
	})
}

//updateConfig
func TestUpdateConfig(t *testing.T) {
	//go InitConfig()
	Convey("Test_updateConfig", t, func() {
		Convey("test case 1", func() {
			entities := &channel.HBResponse{}
			updateConfig(entities)
		})
	})
}

//updateAgent
func TestUpdateAgent(t *testing.T) {
	//go InitConfig()
	Convey("Test_updateAgent", t, func() {
		/*Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(logs.GetLogger().Debug, func(v ...interface{}) {

			})
			entities := &channel.HBResponse{
				Version: utils.AgentVersion,
			}
			updateAgent(entities)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(logs.GetLogger().Error, func(v ...interface{}) error {
				return nil
			})
			entities := &channel.HBResponse{
				DownloadUrl: "",
			}
			updateAgent(entities)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(logs.GetLogger().Error, func(v ...interface{}) error {
				return nil
			})
			entities := &channel.HBResponse{
				Md5: "",
			}
			updateAgent(entities)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(logs.GetLogger().Errorf, func(format string, params ...interface{}) error {
				//syscall.Exit(0)
				//configFromReader
				return nil
			})
			mockfn.Replace(upgrade.Download, func(url string, version string, md5str string) error {
				return errors.New("")
			})
			entities := &channel.HBResponse{
				Md5: "",
			}
			updateAgent(entities)
		})*/
	})
}

//buildHeartBeatUrl
func TestBuildHeartBeatUrl(t *testing.T) {
	//go InitConfig()
	Convey("Test_buildHeartBeatUrl", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(cesConfig.GetConfig, func() *cesConfig.CESConfig {
				return &cesConfig.CESConfig{}
			})
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			url := buildHeartBeatUrl("url")
			So(url, ShouldEqual, "/V1.0/url")
		})
	})
}

//SendSignalHeartBeat
func TestSendSignalHeartBeat(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendSignalHeartBeat", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(buildHeartBeatUrl, func(uri string) string {
				return ""
			})
			mockfn.Replace(sendHeartBeat, func(url string, hb *channel.HBEntity) *channel.HBResponse {
				return nil
			})
			SendSignalHeartBeat(nil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(buildHeartBeatUrl, func(uri string) string {
				return ""
			})
			mockfn.Replace(sendHeartBeat, func(url string, hb *channel.HBEntity) *channel.HBResponse {
				return &channel.HBResponse{}
			})
			SendSignalHeartBeat(nil)
		})
	})
}

//sendHeartBeat
func TestSendHeartBeat(t *testing.T) {
	//go InitConfig()
	Convey("Test_sendHeartBeat", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(json.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("123")
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("123")
			})
			entity := &channel.HBEntity{}
			beat := sendHeartBeat("", entity)
			So(beat, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(nil),
				}
				return response, errors.New("")
			})
			entity := &channel.HBEntity{}
			beat := sendHeartBeat("", entity)
			So(beat, ShouldBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), errors.New("123")
			})
			entity := &channel.HBEntity{}
			beat := sendHeartBeat("", entity)
			So(beat, ShouldBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			//utils.HTTPSend
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			})
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), errors.New("123")
			})
			entity := &channel.HBEntity{}
			beat := sendHeartBeat("", entity)
			So(beat, ShouldBeNil)
		})
	})
}

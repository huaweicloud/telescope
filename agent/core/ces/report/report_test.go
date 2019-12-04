package report

import (
	. "github.com/smartystreets/goconvey/convey"

	"bytes"
	"errors"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
	"github.com/yougg/mockfn"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//SendMetricData
func TestSendMetricData(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendMetricData", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildCesMetricData, func(inputMetric *model.InputMetric, isAggregated bool) model.CesMetricDataArr {
				return nil
			})
			mockfn.Replace(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(json.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("12123")),
				}
				return response, errors.New("")
			})
			SendMetricData("url", nil, true)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildCesMetricData, func(inputMetric *model.InputMetric, isAggregated bool) model.CesMetricDataArr {
				return nil
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("12123")),
				}
				return response, errors.New("")
			})
			SendMetricData("url", nil, true)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildCesMetricData, func(inputMetric *model.InputMetric, isAggregated bool) model.CesMetricDataArr {
				return nil
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("12123")),
				}
				return response, nil
			})
			SendMetricData("url", nil, true)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildCesMetricData, func(inputMetric *model.InputMetric, isAggregated bool) model.CesMetricDataArr {
				return nil
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(strings.NewReader("12123")),
				}
				return response, nil
			})
			SendMetricData("url", nil, true)
		})
	})
}

//SendProcessInfo
func TestSendProcessInfo(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendProcessInfo", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildProcessInfoByList, func(processList model.ChProcessList) model.ProcessInfoDB {
				return model.ProcessInfoDB{}
			})
			//bytes.NewBuffer
			mockfn.Replace(bytes.NewBuffer, func(buf []byte) *bytes.Buffer {
				return nil
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {

				return nil, errors.New("")
			})
			SendProcessInfo("url", model.ChProcessList{})
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildProcessInfoByList, func(processList model.ChProcessList) model.ProcessInfoDB {
				return model.ProcessInfoDB{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("121213")),
				}
				return response, nil
			})
			SendProcessInfo("url", nil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(model.BuildProcessInfoByList, func(processList model.ChProcessList) model.ProcessInfoDB {
				return model.ProcessInfoDB{}
			})
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(strings.NewReader("112123")),
				}
				return response, nil
			})
			SendProcessInfo("url", nil)

		})
	})
}

//SendData
func TestSendData(t *testing.T) {
	//go InitConfig()
	Convey("Test_SendData", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("")
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {

				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(strings.NewReader("112123")),
				}
				return response, nil
			})
			SendData("", nil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {

				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("112123")),
				}
				return response, nil
			})
			SendData("", nil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, nil
			})
			mockfn.Replace(utils.HTTPSend, func(req *http.Request, service string) (*http.Response, error) {

				response := &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(strings.NewReader("112123")),
				}
				return response, errors.New("")
			})
			SendData("", nil)
		})
	})
}

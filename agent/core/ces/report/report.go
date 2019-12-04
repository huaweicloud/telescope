package report

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// SendMetricData used for ces post metric-data api
func SendMetricData(url string, data *model.InputMetric, isAggregate bool) {

	metricData, err := json.Marshal(model.BuildCesMetricData(data, isAggregate))

	if err != nil {
		logs.GetCesLogger().Errorf("Failed marshall ces metric data. Error: %s", err.Error())
		return
	}
	logs.GetCesLogger().Debugf("Result metricData to send: %s", string(metricData))
	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(metricData))
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error:", rErr.Error())
	}

	res, err := utils.HTTPSend(request, cesUtils.Service)

	if err != nil {
		logs.GetCesLogger().Errorf("request error %s", err.Error())
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated { //TODO the codes need be optimized
		logs.GetCesLogger().Infof("Send metric successfully, url is: %s", url)
	} else {
		resBody, _ := ioutil.ReadAll(res.Body)
		logs.GetCesLogger().Errorf("Failed to send metric, url is: %s, response code: %d, response body is: %s", url, res.StatusCode, string(resBody))
	}
}

// SendProcessInfo used for ces post process-info api
func SendProcessInfo(url string, plist model.ChProcessList) {

	processData, err := json.Marshal(model.BuildProcessInfoByList(plist))

	if err != nil {
		logs.GetCesLogger().Infof("Failed marshall ces process info.\n")
		return
	}

	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(processData))
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error:", rErr.Error())
	}

	res, err := utils.HTTPSend(request, cesUtils.Service)

	if err != nil {
		logs.GetCesLogger().Errorf("request error %s", err.Error())
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated { //TODO the codes need be optimized
		logs.GetCesLogger().Info("Send process info success")
	} else {
		resBody, _ := ioutil.ReadAll(res.Body)
		logs.GetCesLogger().Infof("Failed to send ces process info, the response code: %d, url is: %s, response body is: %s", res.StatusCode, url, string(resBody))
	}
}

// SendData ...
func SendData(url string, data []byte) {
	logs.GetCesLogger().Debugf("Data to send: %s, url is: %s", string(data), url)
	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error: %s, url is: %s", rErr.Error(), url)
	}

	res, err := utils.HTTPSend(request, cesUtils.Service)

	if err != nil {
		logs.GetCesLogger().Errorf("request error %s, url is: %s", err.Error(), url)
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated {
		logs.GetCesLogger().Infof("Send data successfully, url is: %s", url)
	} else {
		resBody, _ := ioutil.ReadAll(res.Body)
		logs.GetCesLogger().Errorf("Failed to send data, the response code: %d, url is: %s, response body is: %s", res.StatusCode, url, string(resBody))
	}
}

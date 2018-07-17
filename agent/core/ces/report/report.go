package report

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// SendMetricData used for ces post metric-data api
func SendMetricData(client *http.Client, url string, data *model.InputMetric, isAggregate bool) {

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

	res, err := utils.HTTPSend(client, request, ces_utils.Service)

	if err != nil {
		logs.GetCesLogger().Errorf("request error %s", err.Error())
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated { //TODO the codes need be optimized
		logs.GetCesLogger().Info("Send metric success")
	} else {
		logs.GetCesLogger().Infof("Failed to send metric and the response code:%d", res.StatusCode)
	}
}

// SendProcessInfo used for ces post process-info api
func SendProcessInfo(client *http.Client, url string, plist model.ChProcessList) {

	processData, err := json.Marshal(model.BuildProcessInfoByList(plist))

	if err != nil {
		logs.GetCesLogger().Infof("Failed marshall ces process info.\n")
		return
	}

	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(processData))
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error:", rErr.Error())
	}

	res, err := utils.HTTPSend(client, request, ces_utils.Service)

	if err != nil {
		logs.GetCesLogger().Errorf("request error %s", err.Error())
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated { //TODO the codes need be optimized
		logs.GetCesLogger().Info("Send process info success")
	} else {
		logs.GetCesLogger().Infof("Failed to send ces process info and the response code:%d", res.StatusCode)
	}
}

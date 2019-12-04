package heartbeat

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	cesConfig "github.com/huaweicloud/telescope/agent/core/ces/config"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// HeartBeat ...
type HeartBeat struct {
	CesDetails string
}

// LoadHbServicesDetails ...
// load services(ces) detail content from service data channel
func (hb *HeartBeat) LoadHbServicesDetails(ch chan channel.HBServiceData) {
	for {
		hbServiceData := <-ch
		if hbServiceData.Service == cesUtils.Service {
			hb.CesDetails = hbServiceData.Detail
		} else {
			logs.GetCesLogger().Warnf("Return service [%s] not matched. ", hbServiceData.Service)
		}
	}
}

// ProduceHeartBeat ...
// Start the heartbeat timer, it will send heartbeat message to heartbeat channel periodically
func (hb *HeartBeat) ProduceHeartBeat(heartbeat chan *channel.HBEntity) {
	cronTime := utils.HB_CRON_JOB_TIME_SECOND //set default
	ticker := time.NewTicker(time.Duration(cronTime) * time.Second)

	i, j := 0, 0
	for range ticker.C {
		heartbeat <- channel.NewHBEntity(channel.Running, time.Now().Unix()*1000, hb.CesDetails)
		logs.GetCesLogger().Infof("Start to produce heartbeat and current goroutine number is: %d", runtime.NumGoroutine())
		//clear last hb details to avoid duplicate metadata send to hb server
		hb.CesDetails = ""
		//support hot load services and common config file
		utils.ReloadConfig(&i, &j)
		cesConfig.ReloadConfig()
	}
}

// ConsumeHeartBeat ...
// Start the control service, it will keep receiving the heartbeat and re-send it to server
func (hb *HeartBeat) ConsumeHeartBeat(heartbeat chan *channel.HBEntity) {
	for {
		HBData := <-heartbeat
		hbResponse := sendHeartBeat(buildHeartBeatUrl(utils.POST_HEART_BEAT_URI), HBData)
		if hbResponse != nil {
			updateConfig(hbResponse)
			go updateAgent(hbResponse)
		} else {
			logs.GetCesLogger().Errorf("Failed to send heart beat, so current heartbeat entity is dismissed.")
		}
	}
}

func updateConfig(hbResponse *channel.HBResponse) {
	//put services(ces) config to config channel
	channel.GetCesConfigChan() <- hbResponse.CesConfig
}

func updateAgent(hbResponse *channel.HBResponse) {
	if hbResponse.Version == utils.AgentVersion {
		logs.GetCesLogger().Debug("Agent version matches the version supported by server and do not need to update.")
		return
	}
	if hbResponse.DownloadUrl == "" {
		logs.GetCesLogger().Error("Prepare UpdateAgent failed, hbResponse DownloadUrl is empty.")
		return
	}
	if hbResponse.Md5 == "" {
		logs.GetCesLogger().Error("Prepare UpdateAgent failed, hbResponse Md5 is empty.")
		return
	}

	err := upgrade.Download(hbResponse.DownloadUrl, hbResponse.Version, hbResponse.Md5)
	if err != nil {
		logs.GetCesLogger().Errorf("Download new package failed, err:%s", err.Error())
	}
}

func buildHeartBeatUrl(uri string) string {
	return cesConfig.GetConfig().Endpoint + utils.SLASH + utils.APICESVersion + utils.SLASH + utils.GetConfig().ProjectId + uri
}

// SendSignalHeartBeat ...
func SendSignalHeartBeat(hb *channel.HBEntity) {
	sigHBResponse := sendHeartBeat(buildHeartBeatUrl(utils.POST_HEART_BEAT_URI), hb)

	if sigHBResponse != nil {
		logs.GetCesLogger().Info("Success to send agent signal heartbeat.")
	} else {
		logs.GetCesLogger().Error("Failed to send agent signal heartbeat.")
	}
}

//send heartbeat to server
func sendHeartBeat(url string, hb *channel.HBEntity) *channel.HBResponse {
	hbEntityBytes, err := json.Marshal(*hb)
	if err != nil {
		logs.GetCesLogger().Infof("Failed marshall ces heartbeat, error is %s", err.Error())
		return nil
	}
	logs.GetCesLogger().Debugf("Heartbeat request url is: %s", url)

	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(hbEntityBytes))
	if rErr != nil {
		logs.GetCesLogger().Errorf("Create request Error:", rErr.Error())
		return nil
	}

	res, err := utils.HTTPSend(request, "HB")
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to request for server, error is %s", err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		logs.GetCesLogger().Debug("Success to send heartbeat.")
		hbResponse := channel.HBResponse{}
		resBodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.GetCesLogger().Debugf("HeartBeat response: %s", string(resBodyBytes))
		err = json.Unmarshal(resBodyBytes, &hbResponse)
		if err != nil {
			logs.GetCesLogger().Errorf("Failed to unmarshal response [%s].", string(resBodyBytes))
			return nil
		}
		return &hbResponse
	} else {
		resBodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.GetCesLogger().Errorf("Failed to send heartbeat and the response code [%d], response content is %s", res.StatusCode, string(resBodyBytes))
		return nil
	}
}

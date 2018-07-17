package heartbeat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	ces_config "github.com/huaweicloud/telescope/agent/core/ces/config"
	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/channel"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts"
	lts_config "github.com/huaweicloud/telescope/agent/core/lts/config"
	lts_errs "github.com/huaweicloud/telescope/agent/core/lts/errs"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	"github.com/huaweicloud/telescope/agent/core/upgrade"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

type HeartBeat struct {
	LtsDetails string
	CesDetails string
}

//load services(lts,ces) detail content from service data channel
func (hb *HeartBeat) LoadHbServicesDetails(ch chan channel.HBServiceData) {
	for {
		hbServiceData := <-ch
		if hbServiceData.Service == lts_utils.SERVICE {
			hb.LtsDetails = hbServiceData.Detail
		} else if hbServiceData.Service == ces_utils.Service {
			hb.CesDetails = hbServiceData.Detail
		} else {
			logs.GetLogger().Warnf("Return service [%s] not matched. ", hbServiceData.Service)
		}
	}
}

// Start the heartbeat timer, it will send heartbeat message to heartbeat channel periodically
func (hb *HeartBeat) ProduceHeartBeat(heartbeat chan *channel.HBEntity) {
	cronTime := utils.HB_CRON_JOB_TIME_SECOND //set default
	ticker := time.NewTicker(time.Duration(cronTime) * time.Second)

	for _ = range ticker.C {
		heartbeat <- channel.NewHBEntity(channel.Running, time.Now().Unix()*1000, lts_config.GetConfig().Enable, lts_errs.GetLtsDetail(), hb.CesDetails)
		//clear last hb details to avoid duplicate metadata send to hb server
		hb.CesDetails = ""
		hb.LtsDetails = ""
		//support hot load services and common config file
		utils.ReloadConfig()
		lts.CronUpdateConfig()
		ces_config.ReloadConfig()
	}
}

// Start the control service, it will keep receiving the heartbeat and re-send it to server
func (hb *HeartBeat) ConsumeHeartBeat(heartbeat chan *channel.HBEntity) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	for {
		HBData := <-heartbeat
		hbResponse := sendHeartBeat(client, buildHeartBeatUrl(utils.POST_HEART_BEAT_URI), HBData)
		if hbResponse != nil {
			updateAgent(hbResponse)
		} else {
			logs.GetLogger().Errorf("Failed to send heart beat, so current heartbeat entity is dismissed.")
		}
	}
}

func updateAgent(hbReponse *channel.HBResponse) {
	if hbReponse.Version != utils.AGENT_VERSION {
		err := upgrade.Download(hbReponse.DownloadUrl, hbReponse.Version, hbReponse.Md5)
		if err != nil {
			logs.GetLogger().Errorf("Download new package failed, err:%s", err.Error())
		}
	}

	//put services(lts,ces) config to config channel
	channel.GetLtsConfigChan() <- hbReponse.LtsConfig
	channel.GetCesConfigChan() <- hbReponse.CesConfig
}

func buildHeartBeatUrl(uri string) string {
	return ces_config.GetConfig().Endpoint + "/" + utils.API_CES_VERSION + "/" + utils.GetConfig().ProjectId + uri
}

func SendSignalHeartBeat(hb *channel.HBEntity) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	sigHBResponse := sendHeartBeat(client, buildHeartBeatUrl(utils.POST_HEART_BEAT_URI), hb)

	if sigHBResponse != nil {
		logs.GetLogger().Info("Success to send agent signal hearbeat.")
	} else {
		logs.GetLogger().Error("Failed to send agent signal hearbeat.")
	}
}

//send heartbeat to server
func sendHeartBeat(client *http.Client, url string, hb *channel.HBEntity) *channel.HBResponse {
	hbEntityBytes, err := json.Marshal(*hb)
	if err != nil {
		logs.GetLogger().Infof("Failed marshall ces heartbeat, error is %s", err.Error())
		return nil
	}
	logs.GetLogger().Debugf("Heartbeat request url is: %s", url)
	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(hbEntityBytes))
	if rErr != nil {
		logs.GetLogger().Errorf("Create request Error:", rErr.Error())
		return nil
	}

	res, err := utils.HTTPSend(client, request, "HB")

	if err != nil {
		logs.GetLogger().Errorf("Failed to request for server, error is %s", err.Error())
		return nil
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		logs.GetLogger().Debug("Success to send heartbeat.")
		hbResponse := channel.HBResponse{}
		resBodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.GetLogger().Debugf("HeartBeat response: %s", string(resBodyBytes))
		err = json.Unmarshal(resBodyBytes, &hbResponse)
		if err != nil {
			logs.GetLogger().Errorf("Failed to unmarshal response [%s].", string(resBodyBytes))
			return nil
		}
		return &hbResponse
	} else {
		resBodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.GetLogger().Errorf("Failed to send heartbeat and the response code [%d], response content is %s", res.StatusCode, string(resBodyBytes))
		return nil
	}
}

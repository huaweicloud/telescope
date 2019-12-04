package heartbeat

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"

	assistantUtils "github.com/huaweicloud/telescope/agent/core/assistant/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func SendHBTicker(switchChan chan bool) {
	// TODO 从配置文件或server端获取间隔时间
	ticker := time.NewTicker(time.Second * assistantUtils.SEND_HEARTBEAT_INTERVAL_IN_SECOND)
	for range ticker.C {
		logs.GetAssistantLogger().Info("Send heart beat starts")
		go SendHBExec(switchChan)
	}
}

// SendHBExec is the heartbeat timer entry
func SendHBExec(switchChan chan bool) {
	heartbeatPtr := getHB()
	err, hbResp := sendHB(heartbeatPtr)
	updateHBState(err)
	if hbResp != nil && hbResp.Config != nil && hbResp.Config.AssistSwitch {
		logs.GetAssistantLogger().Debugf("Assist from server switch on, now.")
		switchChan <- true
	} else {
		logs.GetAssistantLogger().Debugf("Assist from server switch off, now.")
		switchChan <- false
	}
}

// updateHBState updates the heartbeat sent state
func updateHBState(err error) {
	now := time.Now()
	HBS.LastUpdateTime = now
	if err == nil {
		HBS.LastUpdateSucceededTime = now
		HBS.IsEnvReported = true
	}
}

// sendHB sends heartbeat to server periodically
func sendHB(requestBody *Heartbeat) (error, *HBResp) {
	logs.GetAssistantLogger().Debugf("Send heartbeat request body is: %v", *requestBody)
	url := assistantUtils.BuildURL(assistantUtils.POST_HEARTBEAT_URI)
	requestBodyInBytes, err := assistantUtils.GetMarshalledRequestBody(*requestBody, url)
	if err != nil {
		return err, nil
	}
	logs.GetAssistantLogger().Debugf("MarshalledRequestBody is: %v", string(requestBodyInBytes))

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyInBytes))
	if err != nil {
		logs.GetAssistantLogger().Errorf("[Heartbeat]Create request Error: %v", err)
		return err, nil
	}

	res, err := utils.HTTPSend(request, assistantUtils.SERVICE)
	var resBody []byte
	if err != nil {
		if res != nil {
			resBody, _ = ioutil.ReadAll(res.Body)
		}
		logs.GetAssistantLogger().Errorf("[Heartbeat]Failed to request, url is: %s, error is: %v, response body is: %s", assistantUtils.POST_HEARTBEAT_URI, err, string(resBody))
		return err, nil
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		if res != nil {
			resBody, _ = ioutil.ReadAll(res.Body)
		}
		logs.GetAssistantLogger().Errorf("[Heartbeat]Request failed, status code is %d(should be %d or %d), url is: %s, response body is: %s", res.StatusCode, http.StatusOK, http.StatusNoContent, url, string(resBody))
		return errors.New("status code does not match"), nil
	}
	logs.GetAssistantLogger().Debugf("Send heartbeat successfully.")

	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		resBodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logs.GetAssistantLogger().Errorf("[Heartbeat] ioutil.ReadAll failed, error is: %s, res.Body is: %v", err.Error(), res.Body)
			return err, nil
		}
		logs.GetAssistantLogger().Debugf("[Heartbeat]Response: %s", string(resBodyBytes))

		hbRespBody := HBResp{}
		err = json.Unmarshal(resBodyBytes, &hbRespBody)
		if err != nil {
			logs.GetAssistantLogger().Errorf("[Heartbeat]Failed to unmarshal response [%s], error is: %v.", string(resBodyBytes), err)
			return err, nil
		}
		return nil, &hbRespBody
	} else {
		logs.GetAssistantLogger().Debugf("[Heartbeat]Response no content.")
		return nil, nil
	}
}

// getHB
func getHB() *Heartbeat {
	return &Heartbeat{
		InstanceID: utils.GetConfig().InstanceId,
		Version:    assistantUtils.ASSISTANT_VERSION,
		User:       getCurrentUser(),
		PID:        os.Getpid(),
		// TODO 目前依赖telescope，独立后处理
		Status: AGENT_STATUS_RUNNING,
		// TODO
		Metric: nil,
		// TODO
		Extension: []Extension{},
		EnvInfo:   getEnv(),
	}
}

// getEnv executes only when the server never received it
func getEnv() *Env {
	if enableEnvReport() {
		return &Env{
			Hostname: getHostname(),
			IP:       getFirstIP(),
			Arch:     runtime.GOARCH,
			OS:       runtime.GOOS,
		}
	} else {
		return nil
	}
}

// getCurrentUser returns the user name who runs the process
func getCurrentUser() string {
	u, err := user.Current()
	if err != nil {
		logs.GetAssistantLogger().Errorf("Get current user failed and error is: %v", err)
		return ""
	}

	return u.Username
}

// getHostname returns the hostname of os
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		logs.GetAssistantLogger().Errorf("Get hostname failed and error is: %v", err)
		return ""
	}

	return hostname
}

// getFirstIP returns the first IP(non loopback ip) of first interface
func getFirstIP() string {
	var ip net.IP

	interfaces, err := net.Interfaces()
	if err != nil {
		logs.GetAssistantLogger().Errorf("Get interfaces failed and error is: %v", err)
		return ip.String()
	}

L:
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			logs.GetAssistantLogger().Errorf("Get address from interface failed and error is: %v", err)
			return ip.String()
		}

		for _, address := range addresses {
			switch v := address.(type) {
			// 测试时发现类型都是IPNet，遗留IPNet和IPAddr的区别
			case *net.IPNet:
				ip = v.IP
				if !ip.IsLoopback() {
					break L
				}
			case *net.IPAddr:
				ip = v.IP
				if !ip.IsLoopback() {
					break L
				}
			}
		}
	}

	return ip.String()
}

func enableEnvReport() bool {
	return !HBS.IsEnvReported
}

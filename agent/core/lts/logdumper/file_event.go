package logdumper

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	lts_config "github.com/huaweicloud/telescope/agent/core/lts/config"
	"github.com/huaweicloud/telescope/agent/core/lts/file"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	windowslog "github.com/huaweicloud/telescope/agent/core/lts/windowslog"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)

var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary

type MetaData_ struct {
	HasEventsTooOld bool `json:"has_events_too_old"`
	IsLogTooLarge   bool `json:"is_log_too_large"`
	TopicNotExist   bool `json:"topic_or_group_not_exist"`
}
type Response struct {
	StatusCode      int       `json:_`
	MetaData        MetaData_ `json:"meta_data"`
	SuccessPutCount int       `json:"success_put_count"`
}

type ErrorResponseMessage struct {
	Details string `json:"details"`
	Code    string `json:"code"`
}

type ErrorResponse struct {
	Message ErrorResponseMessage `json:"message"`
}

type FileEvent struct {
	IsWindowsOsLog            bool
	WindowsOsLogChannelState  windowslog.WindowsOsLogChannelState
	FileState                 file.FileState
	LogEvent                  LogEventMessage
	Offset                    uint64 //读取后的offset
	ErrRes                    ErrorResponse
	ResponseStr               string
	ResStatusCode             int
	SuccessPreProcessLogEvent bool //before send to server,log need to be marshaled and compressed,it's flaged the status
}

//send the log event from file event to server
func (fileEvent *FileEvent) SendLogDataToServer() {
	logEventMessage := fileEvent.LogEvent
	logBytes, err := jsonx.Marshal(logEventMessage)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed marshall log event message.")
		fileEvent.SuccessPreProcessLogEvent = false
		return
	}

	var compressErr error

	//compress the http content
	logBytes, compressErr = gzipCompress(logBytes)

	if compressErr != nil {
		logs.GetLtsLogger().Errorf("The log content is failed to compress, compress error: %s", compressErr.Error())
		fileEvent.SuccessPreProcessLogEvent = false
		return
	}
	fileEvent.SuccessPreProcessLogEvent = true
	var uri string = "/groups/" + fileEvent.LogEvent.LogGroupId + "/logs"
	//send logs to server,if due to internet issue, should start retry,there are 2 scenario about internet issue
	//1.internet issue from agent to APIGW
	//2.internet issue from APIGW to Logtank Server, but this case Agent can get 500 status code
	res := sendLogBytes(uri, logBytes, true)
	defer res.Body.Close()
	statusCode, errorResponse, responseStr := handleResponse(uri, logBytes, true, res)
	fileEvent.ResStatusCode = statusCode
	fileEvent.ErrRes = errorResponse
	fileEvent.ResponseStr = responseStr
}

//build http request
func buildRequest(httpContentEncoding bool, uri string, logBytes []byte) *http.Request {
	var url string = lts_config.GetConfig().Endpoint + "/" + utils.API_LTS_VERSION + "/" + utils.GetConfig().ProjectId + uri
	request, rErr := http.NewRequest("POST", url, bytes.NewBuffer(logBytes))
	if rErr != nil {
		logs.GetLtsLogger().Errorf("Create request error: %s", rErr.Error())
		return nil
	}
	if httpContentEncoding {
		request.Header.Set("Content-Encoding", "gzip")
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("x-log-md5", utils.GetMd5FromBytes(logBytes))
	return request
}

//send log data to server, if internet error, start to retry
func sendLogBytes(uri string, logBytes []byte, httpContentEncoding bool) (res *http.Response) {
	resendCount := 0
	retry := lts_utils.PUT_LOG_MAX_RETRY

	request := buildRequest(httpContentEncoding, uri, logBytes)
	if request == nil {
		logs.GetLtsLogger().Errorf("Failed to create request, please check.")
		return
	}
	var err error
	res, err = utils.HTTPSend(request, lts_utils.SERVICE)
	for {
		if err == nil {
			break
		}

		if resendCount >= retry-1 {
			time.Sleep(lts_utils.PUT_LOG_RETRY_INTERVAL_SEC * time.Second)
			logs.GetLtsLogger().Errorf("Data send %d times,server is not available beacuse:%s.", lts_utils.PUT_LOG_MAX_RETRY, err.Error())
			request = buildRequest(httpContentEncoding, uri, logBytes)
			res, err = utils.HTTPSend(request, lts_utils.SERVICE)
			resendCount++
		} else {
			time.Sleep(lts_utils.PUT_LOG_RETRY_INTERVAL_MS * time.Millisecond)
			logs.GetLtsLogger().Errorf("Data send %d times,server is not available beacuse:%s.", lts_utils.PUT_LOG_MAX_RETRY, err.Error())
			request = buildRequest(httpContentEncoding, uri, logBytes)
			res, err = utils.HTTPSend(request, lts_utils.SERVICE)
			resendCount++
		}
	}
	return
}

//according to the response,go on the flow.
//There is one special scenario APIGW return a 500 response which is not from logtank server.
func handleResponse(uri string, logBytes []byte, httpContentEncoding bool, res *http.Response) (int, ErrorResponse, string) {
	resendCount := 0
	retry := lts_utils.PUT_LOG_MAX_RETRY
	response := ErrorResponse{}
	for {
		if res.StatusCode == http.StatusCreated {
			return res.StatusCode, response, ""
		}
		if resendCount >= retry-1 {
			time.Sleep(lts_utils.PUT_LOG_RETRY_INTERVAL_SEC * time.Second)
			resBodyBytes, _ := ioutil.ReadAll(res.Body)
			logs.GetLtsLogger().Errorf("Data send %d times, the request id is [%s] and response is %s", resendCount+1, res.Header.Get("x-request-id"), string(resBodyBytes))
			_ = jsonx.Unmarshal(resBodyBytes, &response)
			if res.StatusCode >= http.StatusInternalServerError {
				logs.GetLtsLogger().Warnf("Server error [%v], retry... ", http.StatusInternalServerError)
				res = sendLogBytes(uri, logBytes, httpContentEncoding)
			} else {
				logs.GetLtsLogger().Errorf("Status code is %v, response is %s", res.StatusCode, string(resBodyBytes))
				return res.StatusCode, response, string(resBodyBytes)
			}
			resendCount++
		} else {
			time.Sleep(lts_utils.PUT_LOG_RETRY_INTERVAL_MS * time.Millisecond)
			resBodyBytes, _ := ioutil.ReadAll(res.Body)
			_ = jsonx.Unmarshal(resBodyBytes, &response)
			if res.StatusCode >= http.StatusInternalServerError {
				res = sendLogBytes(uri, logBytes, httpContentEncoding)
			} else {
				logs.GetLtsLogger().Errorf("Status code is %v, response is %s", res.StatusCode, string(resBodyBytes))
				return res.StatusCode, response, string(resBodyBytes)
			}
			resendCount++
		}

	}

}

//gzip compress
func gzipCompress(logBytes []byte) (result []byte, err error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(logBytes)
	if err != nil {
		return
	}
	if err = w.Flush(); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}
	resultBytes := b.Bytes()
	result = resultBytes
	return
}

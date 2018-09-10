package utils

import (
	"errors"
	"github.com/huaweicloud/telescope/agent/core/assistant/config"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/json-iterator/go"
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary
// BuildURL build URL string by URI
func BuildURL(destURI string) string {
	var url = config.GetConfig().Endpoint + utils.SLASH + API_ASSISTANT_VERSION + utils.SLASH + utils.GetConfig().ProjectId + destURI
	return url
}

// GetMarshalledRequestBody ...
func GetMarshalledRequestBody(v interface{}, url string) ([]byte, error) {
	marshalledData, err := json.Marshal(v)

	if err != nil {
		logs.GetAssistantLogger().Errorf("Failed marshall request body for URL[%s]", url)
		return nil, errors.New("Failed marshall request body for URL[" + url + "]")
	}

	return marshalledData, err
}
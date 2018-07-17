package utils

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/buger/jsonparser"
)

type GeneralConfig struct {
	InstanceId      string
	ProjectId       string
	AccessKey       string
	SecretKey       string
	AKSKToken       string
	RegionId        string
	ExternalService string
}

type MetaData struct {
	InstanceId      string
	ProjectId       string
	RegionId        string
	ExternalService string
}

type SecurityData struct {
	access        string
	secret        string
	securitytoken string
	expires_at    string
}

var (
	config             *GeneralConfig
	metaData           *MetaData
	mutex              sync.Mutex
	security_data      = SecurityData{}
	security_data_conf = SecurityData{}
	retry_count        = 3
	buse_api_aksk      = true
	//当前方式连续失败调用次数
	now_method_count = 0
)

// Read the config from conf.json
func ReadConfig() (*GeneralConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf.json")
	defer file.Close()

	if err != nil {
		logs.GetLogger().Errorf("Loading general configuration file error: %s", err)
		return nil, Errors.NoConfigFileFound
	}

	decoder := json.NewDecoder(file)
	conf := GeneralConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetLogger().Errorf("Parsing general configuration file error: %s", err)
		return nil, Errors.ConfigFileValidationError
	}
	logs.GetLogger().Infof("Successfully loaded general configuration file")

	// update conf by openstack api
	if metaData == nil {
		logs.GetLogger().Infof("Try to get metadata from openstack")
		metaData, err = getConfFromOpenstack()
		if err != nil {
			logs.GetLogger().Errorf("Can not load openstack meta data, err is %v", err)
		}
	}

	if metaData != nil && (conf.ProjectId != metaData.ProjectId || conf.InstanceId != metaData.InstanceId) {
		logs.GetLogger().Warnf("The projectId or instanceId of config.json is not consistent with metadata, use metadata.\n In conf.json, projectId is [%s], instanceId is [%s]. MetaData is %v", conf.ProjectId, conf.InstanceId, *metaData)
		conf.ProjectId = metaData.ProjectId
		conf.InstanceId = metaData.InstanceId
	}

	if metaData != nil && ces_utils.NameSpace != metaData.ExternalService {
		conf.ExternalService = metaData.ExternalService
	}

	security_data_conf.access = conf.AccessKey
	security_data_conf.secret = conf.SecretKey
	_, err = getAKSKFromOpenStack() //get based on the aksk if valid
	if err != nil {
		logs.GetLogger().Warnf("Failed to get aksk data from openstack, %s", err.Error())
	}

	return &conf, nil
}

// Initialize configuration
func InitConfig() {
	var err error
	if config == nil || err != nil {
		mutex.Lock()
		defer mutex.Unlock()
		config, err = ReadConfig()
		for {
			if config == nil || err != nil {
				logs.GetLogger().Errorf("Init common config file error is %s", err.Error())
				time.Sleep(time.Second * 5)
				config, err = ReadConfig()
			} else {
				return
			}
		}
	}

}

// Get wrapper
func GetConfig() *GeneralConfig {
	if config == nil {
		InitConfig()
	}

	setConfAksk(config)

	return config
}

func setConfAksk(config *GeneralConfig) {
	if buse_api_aksk {
		logs.GetLogger().Debugf("Using ak sk from api.")
		config.AccessKey = security_data.access
		config.SecretKey = security_data.secret
		config.AKSKToken = security_data.securitytoken
	} else {
		logs.GetLogger().Debugf("Using ak sk from config.")
		config.AccessKey = security_data_conf.access
		config.SecretKey = security_data_conf.secret
		config.AKSKToken = ""
	}
}

// Use it after a api reponse, if the response is 401/403, increase the now_method_count
// if now_method_count greater than retry_count, change the get aksk method to the other
func ChooseConfOrApiAksk(needExchange bool) {

	if security_data_conf.access == "" {
		buse_api_aksk = true
		return
	}

	if needExchange {
		now_method_count = now_method_count + 1
		if now_method_count >= retry_count {
			buse_api_aksk = !buse_api_aksk
			now_method_count = 0
		}
	} else {
		now_method_count = 0
	}
}

// Reload config to support hot load config file
func ReloadConfig() *GeneralConfig {
	newCommonConfig, err := ReadConfig()
	for {
		if newCommonConfig == nil || err != nil {
			logs.GetLtsLogger().Errorf("Reload common config error is %s.", err.Error())
			time.Sleep(time.Second * 10)
			newCommonConfig, err = ReadConfig()
		} else {
			config = newCommonConfig
			return config
		}
	}
}

func getConfFromOpenstack() (*MetaData, error) {
	conf := MetaData{}

	metaData, err := HTTPGet(GetOpenstackMetaDataUrl)
	if err != nil {
		return nil, err
	}

	conf.InstanceId, err = jsonparser.GetString([]byte(metaData), "uuid")
	if err != nil {
		return nil, err
	}
	conf.ProjectId, err = jsonparser.GetString([]byte(metaData), "project_id")
	if err != nil {
		return nil, err
	}

	// get if it is BMS from meta_data
	serviceName, _, _, err := jsonparser.Get([]byte(metaData), "meta", ces_utils.TagServiceBMS)
	serviceNameStr := string(serviceName[:])
	if err != nil {
		return nil, err
	}

	if strings.Contains(serviceNameStr, ces_utils.TagBMS) {
		conf.ExternalService = ces_utils.ExternalServiceBMS
		logs.GetLogger().Infof("Get External Service of BMS from metadata.")
	}

	return &conf, nil
}

func getAKSKFromOpenStack() (bool, error) {
	var err error
	var strAkskData string
	if !isNeedFreshAKSK(security_data) {
		return false, nil
	}

	logs.GetLogger().Infof("Need to refresh aksk")
	var bAkSKStrValid bool
	for count := 0; count < retry_count; count++ {
		err = nil
		bAkSKStrValid = false
		strAkskData, err = HTTPGet(OpenStackURL4AKSK)
		if err == nil {
			bAkSKStrValid = isOpenStackAKSKJsonValid(strAkskData)
			if bAkSKStrValid {
				break
			}
		}
	}

	if err != nil {
		logs.GetLogger().Errorf("Failed to httpGet:%s", err.Error())
		return false, nil
	}
	if !bAkSKStrValid {
		logs.GetLogger().Errorf("AKSK data is invalid")
		return false, Errors.AkskStrInvalid
	}

	temp_sec_data, err := parseSecurityToken(strAkskData)
	if err != nil {
		logs.GetLogger().Errorf("Failed to parse AKSK:%s", err.Error())
		return false, nil
	}
	security_data = temp_sec_data
	return true, nil

}

// parse SECRECT_DATA from string aksk
func parseSecurityToken(strAkskData string) (SecurityData, error) {
	temp_security_data := SecurityData{}
	var err error
	temp_security_data.access, err = jsonparser.GetString([]byte(strAkskData), "credential", "access")
	if err != nil {
		return temp_security_data, err
	}
	temp_security_data.secret, err = jsonparser.GetString([]byte(strAkskData), "credential", "secret")
	if err != nil {
		return temp_security_data, err
	}
	temp_security_data.securitytoken, err = jsonparser.GetString([]byte(strAkskData), "credential", "securitytoken")
	if err != nil {
		return temp_security_data, err
	}
	temp_security_data.expires_at, err = jsonparser.GetString([]byte(strAkskData), "credential", "expires_at")
	if err != nil {
		return temp_security_data, err
	}
	return temp_security_data, nil
}

func isNeedFreshAKSK(aksk_data SecurityData) bool {

	if buse_api_aksk && now_method_count >= retry_count-1 {
		return true
	}

	if !buse_api_aksk {
		return false
	}

	var bNeedFresh = false
	//"expires_at": "2017-12-20T14:46:04.488000Z"

	if aksk_data.access == "" {
		bNeedFresh = true
	} else {
		layout := "2006-01-02T15:04:05.000000Z"

		expiretime := aksk_data.expires_at
		expiretimeFmt, err := time.Parse(layout, expiretime)

		if err != nil {
			logs.GetLogger().Errorf("Failed to get format expired time:%s", err.Error())
			return bNeedFresh
		}

		current_time := time.Now().Unix()

		deltaBefore5min, _ := time.ParseDuration("-5m")
		needUpdateTime := expiretimeFmt.Add(deltaBefore5min) // before expire time 5 min

		if current_time >= needUpdateTime.Unix() {
			logs.GetLogger().Debugf("time.Now(): %v, current_time:%v, needUpdateTime: %v", time.Now(), current_time, needUpdateTime.Unix())
			bNeedFresh = true
			logs.GetLogger().Infof("Need to update aksk")
		}
	}

	return bNeedFresh
}

func isOpenStackAKSKJsonValid(strAkskData string) bool {

	legalAKSKJsonKey := []string{"expires_at", "access", "secret", "securitytoken"}

	for _, value := range legalAKSKJsonKey {
		if !strings.Contains(strAkskData, value) {
			return false
		}
	}

	return true
}

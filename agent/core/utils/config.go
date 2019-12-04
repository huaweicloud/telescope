package utils

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// GeneralConfig ...
type GeneralConfig struct {
	InstanceId      string
	ProjectId       string
	AppTask         string
	AccessKey       string
	SecretKey       string
	ClientPort      int
	PortNum         int
	AKSKToken       string
	RegionId        string
	ExternalService string
	BmsFlag         bool

	// default false: 10s collect switch off
	DetailMonitorEnable bool
	CPU1stPctThreshold  float64 `json:"cpu_first_pct_threshold"`
	Memory1stThreshold  uint64  `json:"memory_first_threshold"`
	CPU2ndPctThreshold  float64 `json:"cpu_second_pct_threshold"`
	Memory2ndThreshold  uint64  `json:"memory_second_threshold"`
}

// MetaData ...
type MetaData struct {
	InstanceId      string
	ProjectId       string
	RegionId        string
	ExternalService string
	AppTask         string
}

// SecurityData ...
type SecurityData struct {
	access        string
	secret        string
	securitytoken string
	expires_at    string
}

var (
	config               *GeneralConfig
	metaData             *MetaData
	mutex                sync.Mutex
	securityDataFromAPI  = SecurityData{}
	securityDataFromConf = SecurityData{}
	retry_count          = 3
	buse_api_aksk        = true
	//当前方式连续失败调用次数
	now_method_count = 0
	configClientPort = 0
	configPortNum    = 200
	waitTimeArr      = []int{0, 1, 1, 2, 3, 5, 8, 10}

	// DefaultMetricDeltaDataTimeInSecond is the interval to get metrics
	DefaultMetricDeltaDataTimeInSecond = DisableDetailDataCronJobTimeSecond
)

// ReadConfig Read the config from conf.json
func ReadConfig(i, j *int) (*GeneralConfig, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf.json")
	if err != nil {
		logs.GetCesLogger().Errorf("Loading general configuration file error: %s", err)
		return nil, Errors.NoConfigFileFound
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := GeneralConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetCesLogger().Errorf("Parsing general configuration file error: %s", err)
		return nil, Errors.ConfigFileValidationError
	}
	configClientPort = conf.ClientPort
	configPortNum = conf.PortNum
	logs.GetCesLogger().Infof("Successfully loaded general configuration file")

	updateConfigUsingMeta(i, j, &conf)
	if conf.DetailMonitorEnable {
		logs.GetCesLogger().Infof("DetailMonitorEnable is true, set delta time to %d s.", DetailDataCronJobTimeSecond)
		DefaultMetricDeltaDataTimeInSecond = DetailDataCronJobTimeSecond
	}
	return &conf, nil
}

// InitConfig Initialize configuration
func InitConfig() {
	var err error

	if config == nil {
		mutex.Lock()
		defer mutex.Unlock()

		var i, j int = 0, 0
		config, err = ReadConfig(&i, &j)
		for {
			if err == nil && config != nil {
				return
			}

			if err != nil {
				logs.GetCesLogger().Errorf("Init common config file failed and error is:%v", err)
			}
			time.Sleep(time.Second * 5)
			var i, j int = 0, 0
			config, err = ReadConfig(&i, &j)
		}
	}
}

// GetConfig Get wrapper
func GetConfig() *GeneralConfig {
	if config == nil {
		InitConfig()
	}

	setConfAKSK(config)

	return config
}

func setConfAKSK(config *GeneralConfig) {
	if buse_api_aksk {
		logs.GetCesLogger().Debugf("Using ak sk from api.")
		config.AccessKey = securityDataFromAPI.access
		config.SecretKey = securityDataFromAPI.secret
		config.AKSKToken = securityDataFromAPI.securitytoken
	} else {
		logs.GetCesLogger().Debugf("Using ak sk from config.")
		config.AccessKey = securityDataFromConf.access
		config.SecretKey = securityDataFromConf.secret
		config.AKSKToken = ""
	}
}

// ChooseConfOrApiAksk ...
// Use it after a api response, if the response is 401/403, increase the now_method_count
// if now_method_count greater than retry_count, change the get aksk method to the other
func ChooseConfOrApiAksk(needExchange bool) {

	if securityDataFromConf.access == "" {
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

// ReloadConfig to support hot load config file
func ReloadConfig(i, j *int) *GeneralConfig {
	newCommonConfig, err := ReadConfig(i, j)
	for {
		if err == nil && newCommonConfig != nil {
			config = newCommonConfig
			return config
		}

		if err != nil {
			logs.GetCesLogger().Errorf("Reload common config failed and error is:%v.", err)
		}

		time.Sleep(time.Second * 10)
		newCommonConfig, err = ReadConfig(i, j)
	}
}

func getConfFromOpenstack() (*MetaData, error) {
	conf := MetaData{}

	metaData, err := HTTPGet(GetOpenstackMetaDataUrl)
	if err != nil {
		return nil, err
	}

	conf.InstanceId, err = jsonparser.GetString(metaData, "uuid")
	if err != nil {
		return nil, err
	}
	conf.ProjectId, err = jsonparser.GetString(metaData, "project_id")
	if err != nil {
		return nil, err
	}

	conf.AppTask, err = jsonparser.GetString(metaData, "meta", "__app_task")
	logs.GetCesLogger().Debugf("app_task value is :v%", conf.AppTask)
	if err != nil {
		logs.GetCesLogger().Debugf("Get AppTask error is:%v.", err)
	}

	// get if it is BMS from meta_data
	serviceName, _, _, err := jsonparser.Get(metaData, "meta", cesUtils.TagServiceBMS)
	serviceNameStr := string(serviceName[:])
	if err != nil {
		return nil, err
	}

	if strings.Contains(serviceNameStr, cesUtils.TagBMS) {
		conf.ExternalService = cesUtils.ExternalServiceBMS
		logs.GetCesLogger().Infof("Get External Service of BMS from metadata.")
	}

	return &conf, nil
}

func getAKSKFromOpenStack() (bool, error) {
	var err error
	var strAkskData string
	if !isNeedFreshAKSK(securityDataFromAPI) {
		return false, nil
	}

	logs.GetCesLogger().Infof("Need to refresh aksk")
	var bAkSKStrValid bool
	for count := 0; count < retry_count; count++ {
		err = nil
		bAkSKStrValid = false
		bytes, err := HTTPGet(OpenStackURL4AKSK)
		strAkskData = string(bytes)
		if err == nil {
			bAkSKStrValid = isOpenStackAKSKJsonValid(strAkskData)
			if bAkSKStrValid {
				break
			}
		}
	}

	if err != nil {
		logs.GetCesLogger().Errorf("Failed to httpGet:%s", err.Error())
		return false, nil
	}
	if !bAkSKStrValid {
		logs.GetCesLogger().Errorf("AKSK data is invalid")
		return false, Errors.AkskStrInvalid
	}

	temp_sec_data, err := parseSecurityToken(strAkskData)
	if err != nil {
		logs.GetCesLogger().Errorf("Failed to parse AKSK:%s", err.Error())
		return false, nil
	}
	securityDataFromAPI = temp_sec_data
	return true, nil

}

// parse SECRECT_DATA from string aksk
func parseSecurityToken(strAkSk string) (SecurityData, error) {
	var (
		s SecurityData
	)

	for _, v := range []string{"access", "secret", "securitytoken", "expires_at"} {
		str, err := jsonparser.GetString([]byte(strAkSk), "credential", v)
		if err != nil {
			logs.GetCesLogger().Errorf("Exec jsonparser.GetString failed and error is:%v, when parse %s", err, v)
			return s, err
		}

		switch v {
		case "access":
			s.access = str
		case "secret":
			s.secret = str
		case "securitytoken":
			s.securitytoken = str
		case "expires_at":
			s.expires_at = str
		}
	}

	return s, nil
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
			logs.GetCesLogger().Errorf("Failed to get format expired time:%s", err.Error())
			return bNeedFresh
		}

		current_time := time.Now().Unix()

		deltaBefore5min, _ := time.ParseDuration("-5m")
		needUpdateTime := expiretimeFmt.Add(deltaBefore5min) // before expire time 5 min

		if current_time >= needUpdateTime.Unix() {
			logs.GetCesLogger().Debugf("time.Now(): %v, current_time:%v, needUpdateTime: %v", time.Now(), current_time, needUpdateTime.Unix())
			bNeedFresh = true
			logs.GetCesLogger().Infof("Need to update aksk")
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

// GetClientPort ...
func GetClientPort() int {

	// reserved port cannot used, 0 for random port
	if configClientPort < 1024 || configClientPort > 65535 {
		return 0
	}

	if configPortNum <= 0 {
		configPortNum = 200
	}

	if configClientPort+configPortNum > 65536 {
		configPortNum = 65536 - configClientPort
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return configClientPort + r.Intn(configPortNum)
}

// updateBasicInfo updates projectID/instanceID
func updateBasicInfo(conf *GeneralConfig) {
	var err error

	if metaData == nil {
		logs.GetCesLogger().Infof("Try to get metadata from openstack")
		metaData, err = getConfFromOpenstack()
		if err != nil {
			logs.GetCesLogger().Errorf("Can not load openstack meta data, err is %v", err)
			return
		}
		if metaData == nil {
			logs.GetCesLogger().Errorf("Meta data from openstack is nil")
			return
		}
	}

	if conf.ProjectId != metaData.ProjectId || conf.InstanceId != metaData.InstanceId {
		logs.GetCesLogger().Warnf("The projectId or instanceId of conf.json is not consistent with metadata, use metadata. In conf.json, projectId is [%s], instanceId is [%s]. MetaData is %v", conf.ProjectId, conf.InstanceId, *metaData)
		conf.ProjectId = metaData.ProjectId
		conf.InstanceId = metaData.InstanceId
	}

	if cesUtils.NameSpace != metaData.ExternalService {
		conf.ExternalService = metaData.ExternalService
	}
	if metaData.AppTask != "" {
		conf.AppTask = metaData.AppTask
	}
}

// updateAKSK updates ak/sk using openstack API
func updateAKSK(conf *GeneralConfig) {
	securityDataFromConf.access = conf.AccessKey
	securityDataFromConf.secret = conf.SecretKey
	_, err := getAKSKFromOpenStack() //get based on the aksk if valid
	if err != nil {
		logs.GetCesLogger().Warnf("Failed to get aksk data from openstack, %s", err.Error())
		return
	}
	logs.GetCesLogger().Debugf("Update ak/sk successfully")
}

func updateConfigUsingMeta(i, j *int, conf *GeneralConfig) {
	// 如果ak/sk已经获取，则需要按照之前的频率更新
	// nova 应该是在临时ak/sk快过期的时候（可能是过期的5min前）更新ak/sk
	if securityDataFromAPI.access != "" && securityDataFromAPI.secret != "" {
		logs.GetCesLogger().Info("Start to get ak/sk to update config")
		updateAKSK(conf)
	}

	// 获取ak/sk失败时，按照waitTimeArr退避，到达最后一个值时一直按照最后一个值退避
	if waitTimeArr[*i] == *j {
		logs.GetCesLogger().Info("Start to get meta data and ak/sk to update config")
		updateBasicInfo(conf)
		updateAKSK(conf)
		*j = 0
		if *i == len(waitTimeArr)-1 {
			*i = len(waitTimeArr) - 1
		} else {
			*i++
		}
	}

	// projectID/instanceID/externalService use the value in last meta data
	// when the conf.json does not configured
	if config != nil && conf.ProjectId == "" {
		conf.ProjectId = config.ProjectId
	}
	if config != nil && conf.InstanceId == "" {
		conf.InstanceId = config.InstanceId
	}
	if config != nil && conf.ExternalService == "" {
		conf.ExternalService = config.ExternalService
	}

	*j++
}

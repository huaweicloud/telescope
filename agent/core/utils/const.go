package utils

const (
	// HB_CRON_JOB_TIME_SECOND ...
	HB_CRON_JOB_TIME_SECOND            = 60
	DetailDataCronJobTimeSecond        = 10
	DisableDetailDataCronJobTimeSecond = 60
	HTTP_CLIENT_TIME_OUT               = 30
	POST_HEART_BEAT_URI                = "/agent-status"

	AgentNameWin   = "agent.exe"
	AgentNameLinux = "agent"

	DaemonNameWin   = "telescope.exe"
	DaemonNameLinux = "telescope"

	// GetOpenstackMetaDataUrl ...
	GetOpenstackMetaDataUrl = "http://169.254.169.254/openstack/latest/meta_data.json"
	// OpenStackURL4AKSK details
	// https://support.huaweicloud.com/usermanual-ecs/zh-cn_topic_0042400609.html
	OpenStackURL4AKSK = "http://169.254.169.254/openstack/latest/securitykey"

	CPU1stPctThreshold float64 = 10                // 10%
	Memory1stThreshold uint64  = 200 * 1024 * 1024 // 200MB
	CPU2ndPctThreshold float64 = 30                // 30%
	Memory2ndThreshold uint64  = 700 * 1024 * 1024 // 700MB
)

package utils

const HB_CRON_JOB_TIME_SECOND = 60
const DETAIL_DATA_CRON_JOB_TIME_SECOND = 10
const HTTP_CLIENT_TIME_OUT = 30
const POST_HEART_BEAT_URI = "/agent-status"

const (
	AgentNameWin   = "agent.exe"
	AgentNameLinux = "agent"

	DaemonNameWin   = "telescope.exe"
	DaemonNameLinux = "telescope"
)

const GetOpenstackMetaDataUrl = "http://169.254.169.254/openstack/latest/meta_data.json"

const OpenStackURL4AKSK = "http://169.254.169.254/openstack/latest/securitykey"

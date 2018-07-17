package utils

const PER_FILE_EVENT_LOGS_MAX_TOTAL_SIZE = 5 * 1024 * 1024 //one Log Event can hold the max size of the logs content
const PER_FILE_EVENT_LOGS_MAX_NUMBER = 5000                //one log event can hold the max numbers of the logs
const CONTENT_LENGTH_LIMIT_PER_LOG_TEXT = 127 * 1024       //limitation of the content length per log text
const WINDOWS_OS_LOG_PER_COLLECT_MAX_NUMBER = 500          //windows os log 每次收集的最大日志条数

const LOG_File_VALID_DURATION = 7 * 24 * 3600 * 1000         //only extract the logs produced in recent 7 days
const LOG_FILE_SYSTEM_TIME_VALID_DURATION = 24 * 3600 * 1000 // 利用系统时间作为日志时间的日志文件有效期：当天

const RECORD_FILE_PATH = "/record.json" //Record file name
const WINDOWS_OS_LOG_RECORD_FILE_PATH = "/windows_os_log_record.json"

const LOG_EXTRACT_CRON_JOB_TIME_SECOND = 5 //the value represent how often the cron job run at a time,need change to be configured for user
const LOG_AGENT_HEART_BEAT_TIME_SECOND = 5 //the value represent how often heart beat send a package to server

const PUT_LOG_MAX_RETRY = 5
const PUT_LOG_RETRY_INTERVAL_SEC = 30
const PUT_LOG_RETRY_INTERVAL_MS = 200
const PUT_LOG_OVER_LIMIT_WAIT_MINITUES = 10

const SERVICE = "LTS"

const COLLECT_LOG_CRON_JOB_TIME_SECOND = 1

const LOGS_NEED_DROP_ERR_CODES = "LTS.0302,LTS.0303,LTS.0305"

const TIME_EXTRACT_MODE_SYSTEM = "SYSTEM_TIME"

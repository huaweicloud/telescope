package errs

type LtsError struct {
	Code    string
	Message string
}

var (
	NO_FILES_FOUNT                  = LtsError{Code: "LTS.AGENT.0001", Message: "There is no files under the extractor path."}
	GET_FILE_STAT_FAILED            = LtsError{Code: "LTS.AGENT.0002", Message: "Get file stat failed."}
	NO_ACCESS_TO_PUT_LOG            = LtsError{Code: "LTS.AGENT.0003", Message: "No access to put log."}
	SEND_LOG_DATA_FAILED            = LtsError{Code: "LTS.AGENT.0004", Message: "Send log data to server failed."}
	BAD_REQUEST_UNKONOW_ERR         = LtsError{Code: "LTS.AGENT.0005", Message: "Unknown reason to put log failed, and http code is 400"}
	AUTHORIZATION_FAILED_UNKOWN_ERR = LtsError{Code: "LTS.AGENT.0006", Message: "Unknown reason to authorize failed, and http code is 401"}
	ERR_WRITE_LTS_CONFIG_FILE       = LtsError{Code: "LTS.AGENT.0007", Message: "Failed to update local lts config file"}
)

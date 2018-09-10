package utils

import (
	"errors"
)

type logdErrors struct {
	NoConfigFileFound         error
	ConfigFileValidationError error
	NoMatchedFileFound        error
	NoTimeInTheLog            error
	NoTimeFormat              error
	AkskStrInvalid            error
}

// Errors ...
var Errors = logdErrors{
	NoConfigFileFound:         errors.New("no config file found, please check if missing"),
	ConfigFileValidationError: errors.New("config file validation failed, please re-check the content on it"),
	NoMatchedFileFound:        errors.New("no matched file found with a pattern"),
	NoTimeInTheLog:            errors.New("there is no time in the log text"),
	NoTimeFormat:              errors.New("no dateformat for the log agent in the configuration file"),
	AkskStrInvalid:            errors.New("aksk data is invalid"),
}

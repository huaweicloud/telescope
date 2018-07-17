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

var Errors = logdErrors{
	NoConfigFileFound:         errors.New("No config file found, please check if conf.json is missing."),
	ConfigFileValidationError: errors.New("Config file validation failed, please re-check the content on it."),
	NoMatchedFileFound:        errors.New("No matched file found with a pattern"),
	NoTimeInTheLog:            errors.New("There is no time in the log text"),
	NoTimeFormat:              errors.New("No dateformat for the log agent in the configuration file"),
	AkskStrInvalid:            errors.New("Aksk data is invalid"),
}

package utils

import (
	"bytes"
	"fmt"
	"net/url"
)

const (
	// FieldRequiredErrCode ...
	FieldRequiredErrCode  = "ParamRequiredError"     //the field value is required
	FieldMaxSizeErrCode   = "FieldMaxSizeError"      //the size of field error code
	FieldLenErrCode       = "FieldLengthError"       //the length of field error code
	FieldValueExistedCode = "FieldValueExistedError" //the field value existed code
)

// Validator ...
type Validator interface {
	Validate() error
	ParseJson(bytes []byte) error
}

// QueryStringParser ...
type QueryStringParser interface {
	Build(v url.Values)
	Parse() (interface{}, error)
}

// An ErrInvalidField represents an invalid field.
type ErrInvalidField interface {
	error

	//error code
	Code() string

	//error message
	Msg() string

	// Field name the error occurred on.
	Field() string

	// SetObject updates the object of the error,eg."User" is the object of the field "name".
	SetObject(string)
}

type errInvalidField struct {
	object string
	field  string
	code   string
	msg    string
}

// Code returns the error code for the type of invalid field.
func (e *errInvalidField) Code() string {
	return e.code
}

// Msg returns the reason the field was invalid, and its context.
func (e *errInvalidField) Msg() string {
	return fmt.Sprintf("%s: %s.", e.msg, e.field)
}

// Error returns the string version of the invalid parameter error.
func (e *errInvalidField) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.msg)
}

// Field Returns the field and object the error occurred.
func (e *errInvalidField) Field() string {
	field := e.object
	if len(field) > 0 {
		field += "."
	}
	field += e.field

	return field
}

// SetObject ...
func (e *errInvalidField) SetObject(obj string) {
	e.object = obj
}

// ErrInvalidFields ...
type ErrInvalidFields struct {
	Object string
	Errs   []ErrInvalidField
}

// Add a new field invalid error
func (fields *ErrInvalidFields) Add(field ErrInvalidField) {
	field.SetObject(fields.Object)
	fields.Errs = append(fields.Errs, field)

}

// Len get number of invalid field error
func (fields *ErrInvalidFields) Len() int {
	return len(fields.Errs)
}

// Error returns the string formatted form of the invalid parameters.
func (fieldsfields ErrInvalidFields) Error() string {
	w := &bytes.Buffer{}

	for _, err := range fieldsfields.Errs {
		fmt.Fprintf(w, "- %s\n", err.Msg())
	}

	return w.String()
}

// ErrFieldRequired An ErrFieldRequired means the field need value,not nil.
type ErrFieldRequired struct {
	errInvalidField
}

// NewErrFieldRequired creates a new required parameter error.
func NewErrFieldRequired(field string) *ErrFieldRequired {
	return &ErrFieldRequired{
		errInvalidField{
			code:  FieldRequiredErrCode,
			field: field,
			msg:   fmt.Sprintf("missing required field"),
		},
	}
}

// ErrFieldMaxSize field size over the max, eg.The max size of field "message" is 10KB
type ErrFieldMaxSize struct {
	errInvalidField
	max int
}

// NewErrFieldMaxSize ...
func NewErrFieldMaxSize(field string, size int) *ErrFieldMaxSize {
	return &ErrFieldMaxSize{
		errInvalidField: errInvalidField{
			code:  FieldMaxSizeErrCode,
			field: field,
			msg:   fmt.Sprintf("field size over the threshold %d", size),
		},
		max: size,
	}
}

// ErrFieldMaxLen field length over the threshold
type ErrFieldMaxLen struct {
	errInvalidField
	threshold int
}

// NewErrFieldMaxLen ...
func NewErrFieldMaxLen(field string, size int) *ErrFieldMaxLen {
	return &ErrFieldMaxLen{
		errInvalidField: errInvalidField{
			code:  FieldLenErrCode,
			field: field,
			msg:   fmt.Sprintf("field length over the threshold %d", size),
		},
		threshold: size,
	}
}

// ErrFieldValueExisted ...
type ErrFieldValueExisted struct {
	errInvalidField
}

// NewErrFieldValueExisted field value existed
func NewErrFieldValueExisted(field string) *ErrFieldValueExisted {
	return &ErrFieldValueExisted{
		errInvalidField: errInvalidField{
			code:  FieldValueExistedCode,
			field: field,
			msg:   "field value has been existed,please retype",
		},
	}
}

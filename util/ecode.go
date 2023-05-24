package util

import (
	"strconv"
)

type ErrorCode struct {
	code int
	msg  string
}

func (e ErrorCode) Error() string {
	return strconv.FormatInt(int64(e.Code()), 10)
}

// Code return error code
func (e ErrorCode) Code() int { return int(e.code) }

// Message return error message
func (e ErrorCode) Message() string {
	return e.msg
}

func NewErrorCode(code int, msg string) ErrorCode {
	return ErrorCode{code: code, msg: msg}
}

var (
	ParamsError           = NewErrorCode(10001, "Params Error")
	InvalidAccountAddress = NewErrorCode(10002, "Invalid Account Address")
	RecordNotFound        = NewErrorCode(10004, "Record Not Found")
)

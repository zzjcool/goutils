package zhttp

import (
	"net/http"

	"github.com/zzjcool/goutils/ferr"
)

/*
所有错误都应当设立对应的错误码，保证敏感信息不会暴露，未定义的错误类型统一返回ErrUnknown
五位错误码规则：
第一位错误的类型
第二、三位错误出现的模块
第四、五位为该模块内部的错误
*/
var (
	OK             = StatusItf(status{httpStatusCode: http.StatusOK, code: 0, message: "OK"})
	ErrUnknown     = StatusItf(status{httpStatusCode: http.StatusInternalServerError, code: 99999, message: "The error unknown"})
	ErrPermission  = StatusItf(status{httpStatusCode: http.StatusForbidden, code: 10001, message: "Need to check permissions"})
	ErrServerPanic = StatusItf(status{httpStatusCode: http.StatusInternalServerError, code: 10002, message: "Server error"})

	// system errors
	ErrAuthentication = StatusItf(status{httpStatusCode: http.StatusUnauthorized, code: 20001, message: "Authentication failed"})
	ErrCreate         = StatusItf(status{httpStatusCode: http.StatusBadRequest, code: 20002, message: "Create failed"})
	ErrParameterMatch = StatusItf(status{httpStatusCode: http.StatusBadRequest, code: 20003, message: "Parameter binding failed"})
	ErrDelete         = StatusItf(status{httpStatusCode: http.StatusBadRequest, code: 20004, message: "Delete failed"})
	ErrGet            = StatusItf(status{httpStatusCode: http.StatusBadRequest, code: 20005, message: "Get failed"})
	ErrUpdate         = StatusItf(status{httpStatusCode: http.StatusBadRequest, code: 20006, message: "Update failed"})
)

type StatusItf interface {
	error
	Message() string
	HttpStatusCode() int
	Code() int
}

type status struct {
	httpStatusCode int
	code           int
	message        string
}

func (s status) Error() string {
	return s.Message()
}

func (s status) Message() string {
	return s.message
}
func (s status) Code() int {
	return s.code
}
func (s status) HttpStatusCode() int {
	return s.httpStatusCode
}

type httpError struct {
	StatusItf
	ferr.Interface
}

func (h *httpError) Error() string {
	return h.Interface.Error()
}

func (h *httpError) Message() string {
	return h.Interface.Error()
}

func NewE(err error, sts StatusItf) *httpError {
	return &httpError{
		StatusItf: sts,
		Interface: ferr.Convert(err),
	}
}

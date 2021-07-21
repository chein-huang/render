package render

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	UnknownCode = "error.data.unknown-error"

	UnknownErrResponseMsg I18nID = "ResponseMsg.UnknownErr"
)

type BuError struct {
	Error_        error
	LogLevel_     Level
	ErrMsg_       interface{}
	ResponseMsg_  interface{}
	ResponseCode_ string
	NeedTrace_    bool
	HttpCode_     int
}

func (e *BuError) Error() string {
	if e.ErrMsg_ == nil {
		if e.Error_ == nil {
			return "nil"
		}
		return e.Error_.Error()
	}

	switch msg := e.ErrMsg_.(type) {
	case func() string:
		return msg()
	case func(err error) string:
		return msg(e.Error_)
	case string:
		return msg
	default:
		return fmt.Sprintf("%v", msg)
	}
}

func (e *BuError) HttpCode() int {
	return e.HttpCode_
}

func (e *BuError) ResponseCode() string {
	return e.ResponseCode_
}

func (e *BuError) ResponseMessage(languages []Language) string {
	switch msg := e.ResponseMsg_.(type) {
	case func() string:
		return msg()
	case func([]Language) string:
		return msg(languages)
	case I18nResource:
		return msg.String(languages)
	case *I18nResource:
		return msg.String(languages)
	case *I18nID:
		return msg.String(languages)
	case I18nID:
		return msg.String(languages)
	case string:
		return msg
	default:
		return fmt.Sprintf("%v", msg)
	}
}

func (e *BuError) LogLevel() Level {
	return e.LogLevel_
}

func (e *BuError) Is(err error) bool {
	unwrapErr := err
	for unwrapErr != nil {
		if pointerErr, ok := unwrapErr.(*BuError); ok && e.ResponseCode() == pointerErr.ResponseCode() {
			return true
		}
		unwrapErr = errors.Unwrap(unwrapErr)
	}
	return false
}

func (e *BuError) Unwrap() error {
	return e.Error_
}

func (e *BuError) LogTrace() bool {
	return e.NeedTrace_
}

func UnknownErr(err error) *BuError {
	return &BuError{
		Error_:        err,
		LogLevel_:     ErrorLevel,
		HttpCode_:     http.StatusInternalServerError,
		ResponseCode_: UnknownCode,
		ResponseMsg_:  UnknownErrResponseMsg.Resource(),
		NeedTrace_:    true,
	}
}

func GetBuError(err error) *BuError {
	unWrap := err
	for unWrap != nil {
		if ret, ok := unWrap.(*BuError); ok {
			return ret
		}
		unWrap = errors.Unwrap(unWrap)
	}
	return UnknownErr(err)
}

func SetError(srcErr, err error) error {
	GetBuError(srcErr).Error_ = err
	return srcErr
}

func SetLogLevel(err error, level Level) error {
	GetBuError(err).LogLevel_ = level
	return err
}

func SetErrMsg(err error, msg string) error {
	GetBuError(err).ErrMsg_ = msg
	return err
}

func SetResponseMsg(err error, msg string) error {
	GetBuError(err).ResponseMsg_ = msg
	return err
}

func SetResponseCode(err error, code string) error {
	GetBuError(err).ResponseCode_ = code
	return err
}

func SetNeedTrace(err error, trace bool) error {
	GetBuError(err).NeedTrace_ = trace
	return err
}

func SetHttpCode(err error, code int) error {
	GetBuError(err).HttpCode_ = code
	return err
}

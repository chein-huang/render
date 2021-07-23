package render_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"text/template"

	. "github.com/chein-huang/render"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	assert := require.New(t)
	var err error
	err = &BuError{
		Error_:  fmt.Errorf("msg by error"),
		ErrMsg_: nil,
	}
	assert.Error(err)
	assert.Equal("msg by error", err.Error())

	err = &BuError{
		ErrMsg_: func() string {
			return "msg by func() string"
		},
	}
	assert.Error(err)
	assert.Equal("msg by func() string", err.Error())

	err = &BuError{
		Error_: fmt.Errorf("an error"),
		ErrMsg_: func(err error) string {
			return "msg by func(err error) string. error: " + err.Error()
		},
	}
	assert.Error(err)
	assert.Equal("msg by func(err error) string. error: an error", err.Error())

	err = &BuError{
		ErrMsg_: "a string error",
	}
	assert.Error(err)
	assert.Equal("a string error", err.Error())

	Struct := struct {
		Msg string
	}{
		Msg: "struct",
	}
	err = &BuError{
		ErrMsg_: Struct,
	}
	assert.Error(err)
	assert.Equal(fmt.Sprintf("%v", Struct), err.Error())
}

func TestHttpCode(t *testing.T) {
	assert := require.New(t)
	var err interface{ HttpCode() int }
	err = &BuError{
		HttpCode_: http.StatusBadRequest,
	}
	assert.Equal(http.StatusBadRequest, err.HttpCode())
}

func TestResponseCode(t *testing.T) {
	assert := require.New(t)
	var err interface{ ResponseCode() string }
	err = &BuError{
		ResponseCode_: "error-code",
	}
	assert.Equal("error-code", err.ResponseCode())
}

func TestResponseMessage(t *testing.T) {
	assert := require.New(t)
	var err interface {
		ResponseMessage([]Language) string
	}
	err = &BuError{
		ResponseMsg_: func() string {
			return "func() string"
		},
	}
	assert.Equal("func() string", err.ResponseMessage(nil))

	err = &BuError{
		ResponseMsg_: func(language []Language) string {
			strs := []string{}
			for _, l := range language {
				strs = append(strs, l.Language)
			}
			return strings.Join(strs, ",")
		},
	}
	assert.Equal("1,2", err.ResponseMessage([]Language{{"1", 1}, {"2", 1}}))

	err = &BuError{
		ResponseMsg_: I18nResource{
			Default: GetTemplate("default", "default"),
			Map: map[string]*template.Template{
				"1": GetTemplate("2", "2"),
				"3": GetTemplate("4", "4"),
			},
		},
	}
	assert.Equal("2", err.ResponseMessage([]Language{{"1", 1}}))
	assert.Equal("4", err.ResponseMessage([]Language{{"3", 1}, {"5", 1}}))
	assert.Equal("default", err.ResponseMessage([]Language{{"5", 1}}))

	err = &BuError{
		ResponseMsg_: &I18nResource{
			Default: GetTemplate("default", "default"),
			Map: map[string]*template.Template{
				"1": GetTemplate("2", "2"),
				"3": GetTemplate("4", "4"),
			},
		},
	}
	assert.Equal("2", err.ResponseMessage([]Language{{"1", 1}}))
	assert.Equal("4", err.ResponseMessage([]Language{{"3", 1}, {"5", 1}}))
	assert.Equal("default", err.ResponseMessage([]Language{{"5", 1}}))

	err = &BuError{
		ResponseMsg_: "string error",
	}
	assert.Equal("string error", err.ResponseMessage(nil))

	Struct := struct {
		Msg string
	}{
		Msg: "struct",
	}
	err = &BuError{
		ResponseMsg_: Struct,
	}
	assert.Equal(fmt.Sprintf("%v", Struct), err.ResponseMessage(nil))
}

func TestLogLevel(t *testing.T) {
	assert := require.New(t)
	var err interface{ LogLevel() Level }
	err = &BuError{
		LogLevel_: ErrorLevel,
	}
	assert.Equal(ErrorLevel, err.LogLevel())
}

func TestIs(t *testing.T) {
	assert := require.New(t)
	var err error
	err = &BuError{
		ErrMsg_:       "an error",
		ResponseCode_: "error-code",
	}

	var err2 error
	err2 = &BuError{
		Error_: &BuError{
			ResponseCode_: "error-code",
		},
		ResponseCode_: "error-code2",
	}

	var err3 error
	err3 = &BuError{
		Error_: &BuError{
			ResponseCode_: "error-code2",
		},
		ResponseCode_: "error-code3",
	}
	assert.True(errors.Is(err, err2))
	assert.False(errors.Is(err, err3))
	assert.True(errors.Is(err3, err2))
	assert.True(errors.Is(err2, fmt.Errorf("wrap2-%w", fmt.Errorf("wrap-%w", err))))
	assert.False(errors.Is(err2, fmt.Errorf("%v", err)))
	assert.Equal(err, GetBuError(fmt.Errorf("wrap2-%w", fmt.Errorf("wrap-%w", err))))
	assert.Equal(UnknownErr(fmt.Errorf("wrap2-wrap")).Error(), GetBuError(fmt.Errorf("wrap2-%w", fmt.Errorf("wrap"))).Error())
}

func TestLogTrace(t *testing.T) {
	assert := require.New(t)
	var err interface{ LogTrace() bool }
	err = &BuError{
		NeedTrace_: true,
	}
	assert.Equal(true, err.LogTrace())
}

func TestUnknownErr(t *testing.T) {
	assert := require.New(t)
	err := UnknownErr(fmt.Errorf("unknown error"))
	assert.Equal("unknown error", err.Error())
	assert.Equal(http.StatusInternalServerError, err.HttpCode())
	assert.Equal(UnknownCode, err.ResponseCode())
	assert.Equal(TemplateString(UnknownErrResponseMsg.Resource().Default, nil), err.ResponseMessage(nil))
	assert.Equal(TemplateString(UnknownErrResponseMsg.Resource().Default, nil), err.ResponseMessage([]Language{{"unknown", 1}}))
	assert.Equal(TemplateString(UnknownErrResponseMsg.Resource().Map["zh_CN"], nil), err.ResponseMessage([]Language{{"zh_CN", 1}, {"unknown", 1}}))
	assert.Equal(TemplateString(UnknownErrResponseMsg.Resource().Map["en_US"], nil), err.ResponseMessage([]Language{{"en_US", 1}, {"unknown", 1}}))
	assert.Equal(ErrorLevel, err.LogLevel())
	assert.True(err.LogTrace())
}

func TestSetter(t *testing.T) {
	assert := require.New(t)

	unknownErr := UnknownErr(fmt.Errorf("unknown"))
	err := fmt.Errorf("wrap1-%w", fmt.Errorf("wrap2-%w", UnknownErr(fmt.Errorf("unknown"))))
	assert.NotEqual(unknownErr.Error(), GetBuError(SetError(err, fmt.Errorf("err2"))).Error())
	assert.NotEqual(unknownErr.LogLevel(), GetBuError(SetLogLevel(err, DebugLevel)))
	assert.NotEqual(unknownErr.Error(), GetBuError(SetErrMsg(err, "err msg")).Error())
	assert.NotEqual(unknownErr.ResponseMessage(nil), GetBuError(SetResponseMsg(err, "response msg")).ResponseMessage(nil))
	assert.NotEqual(unknownErr.ResponseCode(), GetBuError(SetResponseCode(err, "code")).ResponseCode())
	assert.NotEqual(unknownErr.LogTrace(), GetBuError(SetNeedTrace(err, false)).LogTrace())
	assert.NotEqual(unknownErr.HttpCode(), GetBuError(SetHttpCode(err, http.StatusBadRequest)).HttpCode())
}

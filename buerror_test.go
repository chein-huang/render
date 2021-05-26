package render_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/chein-huang/render"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	assert := require.New(t)
	var err error
	err = &render.BuError{
		Error_:  fmt.Errorf("msg by error"),
		ErrMsg_: nil,
	}
	assert.Error(err)
	assert.Equal("msg by error", err.Error())

	err = &render.BuError{
		ErrMsg_: func() string {
			return "msg by func() string"
		},
	}
	assert.Error(err)
	assert.Equal("msg by func() string", err.Error())

	err = &render.BuError{
		Error_: fmt.Errorf("an error"),
		ErrMsg_: func(err error) string {
			return "msg by func(err error) string. error: " + err.Error()
		},
	}
	assert.Error(err)
	assert.Equal("msg by func(err error) string. error: an error", err.Error())

	err = &render.BuError{
		ErrMsg_: "a string error",
	}
	assert.Error(err)
	assert.Equal("a string error", err.Error())

	Struct := struct {
		Msg string
	}{
		Msg: "struct",
	}
	err = &render.BuError{
		ErrMsg_: Struct,
	}
	assert.Error(err)
	assert.Equal(fmt.Sprintf("%v", Struct), err.Error())
}

func TestHttpCode(t *testing.T) {
	assert := require.New(t)
	var err interface{ HttpCode() int }
	err = &render.BuError{
		HttpCode_: http.StatusBadRequest,
	}
	assert.Equal(http.StatusBadRequest, err.HttpCode())
}

func TestResponseCode(t *testing.T) {
	assert := require.New(t)
	var err interface{ ResponseCode() string }
	err = &render.BuError{
		ResponseCode_: "error-code",
	}
	assert.Equal("error-code", err.ResponseCode())
}

func TestResponseMessage(t *testing.T) {
	assert := require.New(t)
	var err interface {
		ResponseMessage([]render.Language) string
	}
	err = &render.BuError{
		ResponseMsg_: func() string {
			return "func() string"
		},
	}
	assert.Equal("func() string", err.ResponseMessage(nil))

	err = &render.BuError{
		ResponseMsg_: func(language []render.Language) string {
			strs := []string{}
			for _, l := range language {
				strs = append(strs, l.Language)
			}
			return strings.Join(strs, ",")
		},
	}
	assert.Equal("1,2", err.ResponseMessage([]render.Language{{"1", 1}, {"2", 1}}))

	err = &render.BuError{
		ResponseMsg_: render.InternationalizationString{
			Default: "default",
			Map: map[string]string{
				"1": "2",
				"3": "4",
			},
		},
	}
	assert.Equal("2", err.ResponseMessage([]render.Language{{"1", 1}}))
	assert.Equal("4", err.ResponseMessage([]render.Language{{"3", 1}, {"5", 1}}))
	assert.Equal("default", err.ResponseMessage([]render.Language{{"5", 1}}))

	err = &render.BuError{
		ResponseMsg_: &render.InternationalizationString{
			Default: "default",
			Map: map[string]string{
				"1": "2",
				"3": "4",
			},
		},
	}
	assert.Equal("2", err.ResponseMessage([]render.Language{{"1", 1}}))
	assert.Equal("4", err.ResponseMessage([]render.Language{{"3", 1}, {"5", 1}}))
	assert.Equal("default", err.ResponseMessage([]render.Language{{"5", 1}}))

	err = &render.BuError{
		ResponseMsg_: "string error",
	}
	assert.Equal("string error", err.ResponseMessage(nil))

	Struct := struct {
		Msg string
	}{
		Msg: "struct",
	}
	err = &render.BuError{
		ResponseMsg_: Struct,
	}
	assert.Equal(fmt.Sprintf("%v", Struct), err.ResponseMessage(nil))
}

func TestLogLevel(t *testing.T) {
	assert := require.New(t)
	var err interface{ LogLevel() render.Level }
	err = &render.BuError{
		LogLevel_: render.ErrorLevel,
	}
	assert.Equal(render.ErrorLevel, err.LogLevel())
}

func TestIs(t *testing.T) {
	assert := require.New(t)
	var err error
	err = &render.BuError{
		ErrMsg_:       "an error",
		ResponseCode_: "error-code",
	}

	var err2 error
	err2 = &render.BuError{
		Error_: &render.BuError{
			ResponseCode_: "error-code",
		},
		ResponseCode_: "error-code2",
	}

	var err3 error
	err3 = &render.BuError{
		Error_: &render.BuError{
			ResponseCode_: "error-code2",
		},
		ResponseCode_: "error-code3",
	}
	assert.True(errors.Is(err, err2))
	assert.False(errors.Is(err, err3))
	assert.True(errors.Is(err2, fmt.Errorf("wrap2-%w", fmt.Errorf("wrap-%w", err))))
	assert.False(errors.Is(err2, fmt.Errorf("%v", err)))
	assert.Equal(err, render.GetBuError(fmt.Errorf("wrap2-%w", fmt.Errorf("wrap-%w", err))))
	assert.Equal(render.UnknownErr(fmt.Errorf("wrap2-wrap")).Error(), render.GetBuError(fmt.Errorf("wrap2-%w", fmt.Errorf("wrap"))).Error())
}

func TestLogTrace(t *testing.T) {
	assert := require.New(t)
	var err interface{ LogTrace() bool }
	err = &render.BuError{
		NeedTrace_: true,
	}
	assert.Equal(true, err.LogTrace())
}

func TestUnknownErr(t *testing.T) {
	assert := require.New(t)
	err := render.UnknownErr(fmt.Errorf("unknown error"))
	assert.Equal("unknown error", err.Error())
	assert.Equal(http.StatusInternalServerError, err.HttpCode())
	assert.Equal(render.UnknownCode, err.ResponseCode())
	assert.Equal(render.UnknownErrResponseMsg.Default, err.ResponseMessage(nil))
	assert.Equal(render.UnknownErrResponseMsg.Default, err.ResponseMessage([]render.Language{{"unknown", 1}}))
	assert.Equal(render.UnknownErrResponseMsg.Map["zh"], err.ResponseMessage([]render.Language{{"zh", 1}, {"unknown", 1}}))
	assert.Equal(render.UnknownErrResponseMsg.Map["en"], err.ResponseMessage([]render.Language{{"en", 1}, {"unknown", 1}}))
	assert.Equal(render.ErrorLevel, err.LogLevel())
	assert.True(err.LogTrace())
}

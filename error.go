package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/chein-huang/requestid"
)

var (
	OpenDetail = false
)

func RenderError(writer interface{}, err error) {
	buErr := GetBuError(err)
	httpCode := buErr.HttpCode()

	result := Result{
		Code: buErr.ResponseCode(),
		Msg:  buErr.ResponseMessage(GetLanguagesFrom(writer)),
	}

	if OpenDetail {
		result.Detail = fmt.Sprintf("%+v", err)
	}

	switch writer := writer.(type) {
	case *gin.Context:
		GinJSON(writer, httpCode, &result)
	case http.ResponseWriter:
		bytes, e := json.Marshal(result)
		if e != nil {
			panic(fmt.Sprintf("marshal failed. error: %v. value: %v", e, result))
		}

		_, e = writer.Write(bytes)
		if e != nil {
			panic(fmt.Sprintf("write failed. error: %v", e))
		}

		writer.WriteHeader(httpCode)
	default:
		panic(fmt.Sprintf("invalid writer: %T", err))
	}

	Log(writer, err, buErr.LogLevel(), buErr.LogTrace() || httpCode == http.StatusInternalServerError)
}

func Log(logger interface{}, err error, level Level, trace bool) {
	msg := "%v"

	if trace {
		msg = "%+v"
	}

	var logf func(level logrus.Level, format string, args ...interface{})
	switch logger := logger.(type) {
	case *logrus.Entry:
		logf = logger.Logf
	case *logrus.Logger:
		logf = logger.Logf
	case *gin.Context:
		logf = requestid.GetLogger(logger).Logf
	default:
		logf = logrus.StandardLogger().Logf
	}
	logf(logrus.Level(level), msg, err)
}

type Result struct {
	Code   string `json:"code"`
	Msg    string `json:"message"`
	Detail string `json:"detail,omitempty"`
}

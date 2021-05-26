package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	// . "gitlab.bj.sensetime.com/fdc/parrotsos-data-backend/libs/errorc"
	// . "gitlab.bj.sensetime.com/fdc/parrotsos-data-backend/libs/requestid"
)

func RenderError(writer interface{}, err error) {
	buErr := GetBuError(err)

	result, httpCode := Result{
		Code: buErr.ResponseCode(),
		Msg:  buErr.ResponseMessage(GetLanguagesFrom(writer)),
	}, buErr.HttpCode()

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

	Log(writer, err, httpCode == http.StatusInternalServerError)
}

func Log(logger interface{}, err error, trace bool) {
	msg := "%v"

	buErr := GetBuError(err)
	if buErr.LogTrace() || trace {
		msg = "%+v"
	}

	switch logger := logger.(type) {
	case *logrus.Entry:
		logger.Logf(logrus.Level(buErr.LogLevel()), msg, err)
	case *logrus.Logger:
		logger.Logf(logrus.Level(buErr.LogLevel()), msg, err)
	case *gin.Context:
		// GetLogger(logger).Logf(logrus.Level(level), msg, err)
	default:
		panic(fmt.Errorf("invalid logger type: %T", logger))
	}
}

type Result struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

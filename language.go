package render

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Language struct {
	Language string
	Q        float32
}

func GetLanguagesFrom(context interface{}) []Language {
	languageStr := ""
	switch c := context.(type) {
	case *gin.Context:
		languageStr = c.Request.Header.Get(http.CanonicalHeaderKey("Accept-Language"))
	case http.Request:
		languageStr = c.Header.Get(http.CanonicalHeaderKey("Accept-Language"))
	case http.Header:
		languageStr = c.Get(http.CanonicalHeaderKey("Accept-Language"))
	case string:
		languageStr = c
	default:
		panic(fmt.Sprintf("invalid context: %T", context))
	}
	return GetLanguages(languageStr)
}

func GetLanguages(header string) []Language {
	languageStrs := strings.Split(header, ",")
	languageList := []Language{}
	for _, languageStr := range languageStrs {
		items := strings.Split(strings.ReplaceAll(languageStr, " ", ""), ";")
		if items[0] == "" {
			continue
		}

		language := Language{
			Language: items[0],
			Q:        1,
		}

		if len(items) > 1 {
			val, err := strconv.ParseFloat(strings.TrimPrefix(items[1], "q="), 32)
			if err != nil {
				continue
			}
			language.Q = float32(val)
		}
		languageList = append(languageList, language)
	}
	sort.Slice(languageList, func(i, j int) bool {
		return languageList[i].Q > languageList[j].Q // 这里反过来降序
	})
	return languageList
}

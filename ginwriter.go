package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinPagy(c *gin.Context, page, rows, total uint, list interface{}) {
	var Pagy struct {
		Page  uint        `json:"page"`
		Rows  uint        `json:"rows"`
		Total uint        `json:"total"`
		List  interface{} `json:"list"`
	}
	Pagy.Page = page
	Pagy.Rows = rows
	Pagy.Total = total
	Pagy.List = list
	GinData(c, Pagy)
}

func GinData(c *gin.Context, data interface{}) {
	GinJSON(c, http.StatusOK, data)
}

func GinSuccess(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func GinJSON(c *gin.Context, httpCode int, data interface{}) {
	c.Set("JSONData", data)
	c.JSON(httpCode, data)
}

func GetJSONData(c *gin.Context) interface{} {
	data, _ := c.Get("JSONData")
	return data
}

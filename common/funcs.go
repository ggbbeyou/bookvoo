package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Ok     int         `json:"ok"`
	Reason string      `json:"reason,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func ResponseJson(c *gin.Context, ok int, reason string, data interface{}) {
	res := Response{
		Ok:     ok,
		Reason: reason,
		Data:   data,
	}
	c.JSON(http.StatusOK, res)
}

func Success(c *gin.Context, data interface{}) {
	ResponseJson(c, 1, "", data)
}

func Fail(c *gin.Context, reason string) {
	logrus.Debugf("[fail] %s, %s", c.Request.RequestURI, reason)
	ResponseJson(c, 0, reason, nil)
}

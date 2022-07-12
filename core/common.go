package core

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/core/base"
)

func D(s1 string) decimal.Decimal {
	ss1, _ := decimal.NewFromString(s1)
	return ss1
}

func Run(config string, router *gin.Engine) {
	base.RunMatching()
}

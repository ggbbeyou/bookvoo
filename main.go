package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/haoex/market"
	"github.com/yzimhao/haoex/tradecore"
	"github.com/yzimhao/haoex/views"
)

func main() {
	router := gin.Default()
	go tradecore.Run("./config.toml", router)
	go market.Run("./config.toml", router)

	//pages
	views.Run("./config.toml", router)
	time.Sleep(time.Second * time.Duration(3))
	router.Run(":8080")
}

package main

import (
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

	router.Run(":8080")
}

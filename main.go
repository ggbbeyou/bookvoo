package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/clearing"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/quotation"
	"github.com/yzimhao/bookvoo/user"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/bookvoo/user/orders"
	"github.com/yzimhao/bookvoo/views"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/yzimhao/utilgo"
	"github.com/yzimhao/utilgo/pack"
)

func main() {
	app := &cli.App{
		Name:  "bookVoo",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.toml", Usage: "config file"},
		},
		Action: func(c *cli.Context) error {
			appStart(c.String("config"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version",
				Action: func(ctx *cli.Context) error {
					pack.ShowVersion()
					return nil
				},
			},
			{
				Name:    "clean",
				Aliases: []string{"cl"},
				Usage:   "clean database",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func appStart(configPath string) {
	utilgo.ViperInit(configPath)

	//log level
	level, _ := logrus.ParseLevel(viper.GetString("main.log_level"))
	logrus.SetLevel(level)

	initModule()
	runModule()

}

//初始化各模块的数据库
func initModule() {

	//后面可以根据不同模块拆分到不同的数据库
	default_db := common.Default_db()
	default_rdc := common.Default_redis()

	base.Init(default_db, default_rdc)
	symbols.Init(default_db, default_rdc)
	//资产
	assets.Init(default_db, default_rdc)
	//订单
	orders.Init(default_db, default_rdc)
	//撮合
	match.Init(default_db, default_rdc)
	//结算
	clearing.Init(default_db, default_rdc)
	//k线行情系统
	quotation.Init(default_db, default_rdc)
}

func runModule() {
	//撮合服务
	match.Run()
	//结算服务
	clearing.Run()
	//用户中心
	user.Run()

	//http api相关接口服务
	router := gin.Default()
	quotation.Run(router)
	views.Run(router)

	viper.SetDefault("main.host", ":8080")
	router.Run(viper.GetString("main.host"))
}

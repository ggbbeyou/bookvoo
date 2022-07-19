package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/market"
	"github.com/yzimhao/bookvoo/user"
	"github.com/yzimhao/bookvoo/views"
	"xorm.io/xorm"

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
			start(c.String("config"))
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
					pack.ShowVersion()
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

func start(config string) {
	c := utilgo.ViperInit(config)
	router := gin.Default()

	default_db := func() *xorm.Engine {
		dsn := c.GetString("db.dsn")
		driver := c.GetString("db.driver")
		conn, err := xorm.NewEngine(driver, dsn)
		if err != nil {
			logrus.Panic(err)
		}
		return conn
	}()

	core.SetDbEngine(default_db)
	user.SetDbEngine(default_db)

	go core.Run(config, router)
	go user.Run(config, router)
	go market.RunWithGinRouter(config, router)
	//pages
	views.Run(config, router)

	c.SetDefault("main.host", ":8080")
	router.Run(c.GetString("main.host"))
}

func clean(config string) {
	//todo
}

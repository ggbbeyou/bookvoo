package base

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common/types"
	gowss "github.com/yzimhao/bookvoo/wss"
	"xorm.io/xorm"
)

var (
	Wss *gowss.Hub
	rdc *redis.Client
)

func Init(db *xorm.Engine, r *redis.Client) {
	symbols.Init(db, r)
	Wss = gowss.NewHub()
	rdc = r
	//go pushMsg()
}

func pushMsg() {
	for {
		ctx := context.Background()
		v := rdc.RPop(ctx, types.WsMessage.Format(nil))
		if v.Val() != "" {
			var vv gowss.MsgBody
			err := json.Unmarshal([]byte(v.Val()), &vv)
			if err != nil {
				logrus.Errorf("parse json error: %s", err)
				continue
			}

			Wss.Broadcast <- vv
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

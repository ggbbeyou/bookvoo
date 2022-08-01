package types

import (
	"fmt"
	"strings"
)

type key string

func format(s string, kvs map[string]string) string {
	ss := s
	for k, v := range kvs {
		ss = strings.Replace(ss, fmt.Sprintf("{%s}", k), v, -1)
	}

	if strings.Contains(ss, "{") && strings.Contains(ss, "}") {
		panic(fmt.Sprintf("please replace the param in %s", ss))
	}
	return ss
}

type RedisKey key

func (r RedisKey) Format(kvs map[string]string) string {
	return format(string(r), kvs)
}

func (r RedisKey) String() string {
	return string(r)
}

const (
	//消息推送队列
	WsMessage RedisKey = "message"
	//新订单通知到撮合系统的队列
	NewOrder RedisKey = "order.new.{symbol}"
	//结算成功推送
	TradeResult RedisKey = "trade.result.{symbol}"
	//推送到行情系统的队列
	MarketSubscribe RedisKey = "trade.result"
)

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

const (
	WsMessage   RedisKey = "message"
	NewOrder    RedisKey = "order.new.{symbol}"
	TradeResult RedisKey = "trade.result.{symbol}"
)

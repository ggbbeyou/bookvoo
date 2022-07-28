package types

import "strings"

type RedisKey string

const (
	Message     RedisKey = "message"
	NewOrder    RedisKey = "new_order_{symbol}"
	TradeResult RedisKey = "trade_result_{symbol}"
)

func (r RedisKey) Symbol(symbol string) string {
	return strings.Replace(string(r), "{symbol}", symbol, -1)
}

func (r RedisKey) String() string {
	if strings.Contains(string(r), "{symbol}") {
		panic("please replace the symbol")
	}
	return string(r)
}

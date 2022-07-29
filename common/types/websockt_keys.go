package types

type Subscribe string

func (r Subscribe) Format(kvs map[string]string) string {
	return format(string(r), kvs)
}

const (
	SubscribeUserId Subscribe = "user.{user_id}"

	SubscribeDepth       Subscribe = "depth.{symbol}"
	SubscribeTradeRecord Subscribe = "trade.record.{symbol}"

	SubscribeKline Subscribe = "kline.{period}.{symbol}"
)

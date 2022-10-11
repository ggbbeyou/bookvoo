package clearing

//结算订单的时候需要锁住订单，不能撤单
//待订单锁释放后，才能撤单

import (
	"context"
	"fmt"
	"time"
)

type ClearingLock struct {
	ctx       context.Context
	ask_id    string
	bid_id    string
	redis_key string
	ask_lock  string
	bid_lock  string
}

func NewClearingLock(ask, bid string) ClearingLock {
	return ClearingLock{
		ctx:       context.Background(),
		ask_id:    ask,
		bid_id:    bid,
		redis_key: fmt.Sprintf("clear_lock:%s.%s", ask, bid),
		ask_lock:  fmt.Sprintf("clear_lock:%s", ask),
		bid_lock:  fmt.Sprintf("clear_lock:%s", bid),
	}
}

func (l *ClearingLock) Lock() error {
	cmd := rdc.SetNX(l.ctx, l.redis_key, 1, time.Duration(0))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	cmd1 := rdc.Incr(l.ctx, l.ask_lock)
	if cmd1.Err() != nil {
		return cmd1.Err()
	}

	cmd2 := rdc.Incr(l.ctx, l.bid_lock)
	if cmd2.Err() != nil {
		return cmd2.Err()
	}
	return nil
}

func (l *ClearingLock) UnLock() error {
	cmd := rdc.Del(l.ctx, l.redis_key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	cmd1 := rdc.Decr(l.ctx, l.ask_lock)
	if cmd1.Err() != nil {
		return cmd1.Err()
	}

	if cmd1.Val() <= 0 {
		rdc.Del(l.ctx, l.ask_lock)
	}

	cmd2 := rdc.Decr(l.ctx, l.bid_lock)
	if cmd2.Err() != nil {
		return cmd2.Err()
	}

	if cmd2.Val() <= 0 {
		rdc.Del(l.ctx, l.bid_lock)
	}

	return nil
}

//判断一个订单是否还存在结算时候加的锁
func ClearingLockExist(order_id string) bool {
	key := fmt.Sprintf("clear_lock:%s", order_id)
	ctx := context.Background()
	cmd := rdc.Get(ctx, key)
	if val, _ := cmd.Int64(); val <= 0 {
		return false
	}
	return true
}

package models

import (
	"fmt"
	"strings"
	"time"
)

type TradeLog struct {
	Id       int64     `xorm:"autoincr pk"`
	Symbol   string    `xorm:"-"`
	At       time.Time `xorm:"notnull"`
	Price    string    `xorm:"decimal(40, 20) notnull"`
	Quantity string    `xorm:"decimal(40, 20) notnull"`
	Amount   string    `xorm:"decimal(40, 20) notnull"`

	AskId     string `xorm:"varchar(64) notnull unique(ask_bid)"`
	BidId     string `xorm:"varchar(64) notnull unique(ask_bid)"`
	tableName string `xorm:"-"`
}

func (t *TradeLog) SetTableName(symbol string) {
	t.tableName = fmt.Sprintf("kline_%s_trade_log", symbol)
}

func (t *TradeLog) TableName() string {
	return t.tableName
}

func (t *TradeLog) autoCreateTable() error {
	if t.Symbol == "" {
		return fmt.Errorf("symbol is empty")
	}
	t.SetTableName(t.Symbol)

	if v, ok := tradeLogTableMap.Load(t.tableName); ok && v != nil {
		return nil
	}

	exist, err := engine.IsTableExist(t.TableName())
	if err != nil {
		return err
	}

	if !exist {
		err := engine.CreateTables(t)
		if err != nil {
			return err
		}

		err = engine.CreateIndexes(t)
		if err != nil {
			return err
		}

		err = engine.CreateUniques(t)
		if err != nil {
			return err
		}
	}
	tradeLogTableMap.Store(t.tableName, true)
	return nil
}

func (t *TradeLog) Save() error {
	err := t.autoCreateTable()
	if err != nil {
		return err
	}
	_, err = engine.Table(t.tableName).Insert(&t)
	return err
}

func (t *TradeLog) Clean() {
	engine.Table(t.tableName).DropTable(t)
}

func (t *TradeLog) GetAt(period Period) (st, et time.Time) {
	at := t.At

	switch period {
	case PERIOD_M1:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute(), 0, 0, time.Local)
		et = st.Add(time.Duration(1) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M3:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%3, 0, 0, time.Local)
		et = st.Add(time.Duration(3) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M5:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%5, 0, 0, time.Local)
		et = st.Add(time.Duration(5) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M15:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%15, 0, 0, time.Local)
		et = st.Add(time.Duration(15) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M30:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%30, 0, 0, time.Local)
		et = st.Add(time.Duration(30) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_H1:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), 0, 0, 0, time.Local)
		et = st.Add(time.Duration(1) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H2:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%2, 0, 0, 0, time.Local)
		et = st.Add(time.Duration(2) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H4:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%4, 0, 0, 0, time.Local)
		et = st.Add(time.Duration(4) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H6:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%6, 0, 0, 0, time.Local)
		et = st.Add(time.Duration(6) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H8:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%8, 0, 0, 0, time.Local)
		et = st.Add(time.Duration(8) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H12:
		st = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%12, 0, 0, 0, time.Local)
		et = st.Add(time.Duration(12) * time.Hour).Add(time.Duration(-1) * time.Second)

	case PERIOD_D1:
		st = time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.Local)
		et = st.AddDate(0, 0, 1).Add(time.Duration(-1) * time.Second)
	case PERIOD_D3:
		st = time.Date(at.Year(), at.Month(), at.Day()-at.Day()%3, 0, 0, 0, 0, time.Local)
		et = st.AddDate(0, 0, 3).Add(time.Duration(-1) * time.Second)

	case PERIOD_W1:
		offset := int(time.Monday - at.Weekday())
		if offset > 0 {
			offset = -6
		}
		weekStart := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
		st = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.Local)
		et = st.AddDate(0, 0, 7).Add(time.Duration(-1) * time.Second)
	case PERIOD_MN:
		st = time.Date(at.Year(), at.Month(), 1, 0, 0, 0, 0, time.Local)
		et = st.AddDate(0, 1, 0).Add(time.Duration(-1) * time.Second)
	}
	return st, et
}

func PushTradeLog(symbol string, at time.Time, askId, bidId, price, qty, amount string) TradeLog {
	symbol = strings.ToLower(symbol)
	tl := TradeLog{
		Symbol:   symbol,
		At:       at,
		Price:    price,
		AskId:    askId,
		BidId:    bidId,
		Quantity: qty,
		Amount:   amount,
	}

	tl.Save()
	return tl
}

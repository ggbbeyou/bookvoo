package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Kline struct {
	Id         int64     `xorm:"autoincr pk" json:"-"`
	OpenAt     time.Time `xorm:"notnull timestamp unique(open_at)" json:"open_at"`            //开盘时间
	Open       string    `xorm:"decimal(40, 20) notnull" json:"open"`                         //开盘价
	High       string    `xorm:"decimal(40, 20) notnull" json:"high"`                         // 最高价
	Low        string    `xorm:"decimal(40, 20) notnull" json:"low"`                          //最低价
	Close      string    `xorm:"decimal(40, 20) notnull" json:"close"`                        //收盘价(当前K线未结束的即为最新价)
	Volume     string    `xorm:"decimal(40, 20) notnull" json:"volume"`                       //成交量
	CloseAt    time.Time `xorm:"timestamp notnull default CURRENT_TIMESTAMP" json:"close_at"` // 收盘时间
	Amount     string    `xorm:"decimal(40, 20) notnull" json:"amount"`                       //成交额
	TradeCnt   int64     `json:"trade_count"`                                                 //成交笔数
	CreateTime time.Time `xorm:"timestamp created" json:"-"`
	UpdateTime time.Time `xorm:"timestamp updated" json:"-"`

	symbol string `xorm:"-" json:"-"`
	period Period `xorm:"-" json:"-"`
}

func NewKline(symbol string, period Period) *Kline {

	kl := Kline{
		symbol: symbol,
		period: period,
	}
	return &kl
}

func (k *Kline) TableName() string {
	return fmt.Sprintf("kline_%s_%s", k.symbol, k.period)
}

func (k *Kline) ToJson() string {
	en, _ := json.Marshal(k)
	return string(en)
}

func (k *Kline) Save() error {
	err := k.autoCreateTable()
	if err != nil {
		return err
	}

	table := k.TableName()

	db := engine.NewSession()
	defer db.Close()

	// db.Begin()
	// defer db.Commit()

	old := Kline{}
	exist, _ := db.Table(table).Where("open_at=?", time2str(k.OpenAt)).Get(&old)
	if !exist {
		_, err := db.Table(table).Insert(&k)
		return err
	} else {
		_, err := db.Table(table).Where("open_at=?", time2str(k.OpenAt)).Update(k)
		return err
	}
}

func (k *Kline) autoCreateTable() error {
	if k.symbol == "" || k.period == "" {
		return fmt.Errorf("symbol or period is null")
	}

	if v, ok := klineTableMap.Load(k.TableName()); ok && v != nil {
		return nil
	}

	exist, err := engine.IsTableExist(k.TableName())
	if err != nil {
		return err
	}

	if !exist {
		err := engine.CreateTables(k)
		if err != nil {
			return err
		}
		err = engine.CreateIndexes(k)
		if err != nil {
			return err
		}
		err = engine.CreateUniques(k)
		if err != nil {
			return err
		}

	}
	klineTableMap.Store(k.TableName(), true)
	return nil
}

func time2str(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

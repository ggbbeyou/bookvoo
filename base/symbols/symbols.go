package symbols

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/common/types"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

type status int

const (
	StatusDisable status = 0
	StatusEnable  status = 1
)

type Symbol struct {
	Id           int       `xorm:"pk autoincr int" json:"id"`
	Symbol       string    `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name         string    `xorm:"varchar(250) notnull" json:"name"`
	ShowPrec     int       `xorm:"default(0)" json:"show_prec"`
	MinPrecision int       `xorm:"default(0)" json:"min_precision"`
	Standard     bool      `xorm:"default(0)" json:"standard"` //是否为本位币
	Status       status    `xorm:"default(0) notnull" json:"-"`
	CreateTime   time.Time `xorm:"timestamp created" json:"-"`
	UpdateTime   time.Time `xorm:"timestamp updated" json:"-"`
}

type Pairs struct {
	Id     int    `xorm:"pk autoincr int" json:"-"`
	Symbol string `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name   string `xorm:"varchar(250) notnull" json:"name"`

	TargetSymbolId   int `xorm:"default(0) unique(symbol_base)" json:"target_symbol_id"`   //交易物品
	StandardSymbolId int `xorm:"default(0) unique(symbol_base)" json:"standard_symbol_id"` //支付货币

	PricePrec      int             `xorm:"default(2)" json:"price_prec"`
	QtyPrec        int             `xorm:"default(0)" json:"qty_prec"`
	AllowMinQty    types.NumberStr `xorm:"decimal(40,20) notnull" json:"allow_min_qty"`
	AllowMaxQty    types.NumberStr `xorm:"decimal(40,20) notnull" json:"allow_max_qty"`
	AllowMinAmount types.NumberStr `xorm:"decimal(40,20) notnull" json:"allow_min_amount"`
	AllowMaxAmount types.NumberStr `xorm:"decimal(40,20) notnull" json:"allow_max_amount"`
	FeeRate        types.NumberStr `xorm:"decimal(40,20) notnull default(0)" json:"fee_rate"`

	Status     status    `xorm:"default(0) notnull" json:"-"`
	CreateTime time.Time `xorm:"timestamp created" json:"-"`
	UpdateTime time.Time `xorm:"timestamp updated" json:"-"`

	Target   Symbol `xorm:"-" json:"target"`
	Standard Symbol `xorm:"-" json:"standard"`
}

func (t *Pairs) FormatAmount(a string) string {
	q, _ := decimal.NewFromString(a)
	return q.StringFixedBank(int32(t.PricePrec))
}

func (t *Pairs) FormatQty(qty string) string {
	q, _ := decimal.NewFromString(qty)
	return q.StringFixedBank(int32(t.QtyPrec))
}

func GetPairBySymbol(symbol string) (*Pairs, error) {
	db := db_engine.NewSession()
	defer db.Close()

	item := Pairs{}
	has, err := db.Table(new(Pairs)).Where("symbol=?", symbol).Get(&item)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("not found trade pair %s", symbol)
	}

	item.Target = GetSymbolInfo(item.TargetSymbolId)
	item.Standard = GetSymbolInfo(item.StandardSymbolId)
	return &item, err
}

func GetSymbolInfo(id int) Symbol {
	var one Symbol
	db_engine.Table(new(Symbol)).Where("id=?", id).Get(&one)
	return one
}

func GetSymbolInfoBySymbol(symbol string) (*Symbol, error) {
	var one Symbol
	has, err := db_engine.Table(new(Symbol)).Where("symbol=?", symbol).Get(&one)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("not found symbol %s", symbol)
	}
	return &one, nil
}

func (s *Symbol) FormatNumber(n string) string {
	q, _ := decimal.NewFromString(n)
	return q.StringFixedBank(int32(s.ShowPrec))
}

func Init(db *xorm.Engine, rdc *redis.Client) {
	db_engine = db

	db_engine.Sync2(
		new(Symbol),
		new(Pairs),
	)
}

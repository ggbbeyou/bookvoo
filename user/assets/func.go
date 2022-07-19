package assets

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

const (
	ROOTUSERID = 0
)

func SetDbEngine(db *xorm.Engine) {
	db_engine = db

	//同步表结构
	err := db_engine.Sync2(
		new(Assets),
		new(assetsLog),
		new(assetFreezeRecord),
	)
	if err != nil {
		logrus.Errorf("sync2: %s", err)
	}
}

//解冻冻结余额
func unfreeze_balance(ses *xorm.Session, user_id int64, symbol_id int, business_id, unfreeze_amount, info string) (bool, error) {

	return true, nil
}

// //余额变动
// func balance_change(ses *xorm.Session, user_id int64, symbol_id int, amount string, info string) (bool, error) {
// 	item := Assets{UserId: user_id, SymbolId: symbol_id}
// 	has, err := ses.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).ForUpdate().Get(&item)

// 	if err != nil {
// 		return false, err
// 	}

// 	before := item.Total

// 	item.Total = number_add(item.Total, amount)
// 	item.Available = number_add(item.Available, amount)

// 	//检查余额是否为负数
// 	if check_amount_lt_zero(item.Total) {
// 		return false, fmt.Errorf("total balance less than zero")
// 	}
// 	if check_amount_lt_zero(item.Available) {
// 		return false, fmt.Errorf("available balance less than zero")
// 	}

// 	if !has {
// 		_, err = ses.Table(new(Assets)).Insert(&item)
// 	} else {
// 		_, err = ses.Table(new(assetsLog)).Where("user_id=? and symbol_id=?", user_id, symbol_id).Update(&item)
// 	}

// 	if err != nil {
// 		return false, err
// 	}

// 	//balance log
// 	lg := assetsLog{
// 		UserId:   user_id,
// 		SymbolId: symbol_id,
// 		Before:   before,
// 		Amount:   amount,
// 		After:    item.Total,
// 		Info:     info,
// 	}
// 	_, err = ses.Table(new(assetsLog)).Insert(&lg)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

func number_add(s1, s2 string) string {
	ss1, _ := decimal.NewFromString(s1)
	ss2, _ := decimal.NewFromString(s2)
	return ss1.Add(ss2).String()
}

func number_sub(s1, s2 string) string {
	ss1, _ := decimal.NewFromString(s1)
	ss2, _ := decimal.NewFromString(s2)
	return ss1.Sub(ss2).String()
}

func check_amount_lt_zero(s string) bool {
	ss, _ := decimal.NewFromString(s)
	if ss.Cmp(decimal.Zero) < 0 {
		return true
	}
	return false
}

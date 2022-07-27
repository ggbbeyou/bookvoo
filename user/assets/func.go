package assets

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
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

func UserAssets(user_id int64, symbol_id int) Assets {
	row := Assets{}
	db_engine.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).Get(&row)
	return row
}

func d(s string) decimal.Decimal {
	ss, _ := decimal.NewFromString(s)
	return ss
}

func number_add(s1, s2 string) string {
	return d(s1).Add(d(s2)).String()
}

func number_sub(s1, s2 string) string {
	return d(s1).Sub(d(s2)).String()
}

func check_number_lt_zero(s string) bool {
	if d(s).Cmp(decimal.Zero) < 0 {
		return true
	} else {
		return false
	}
}

func check_number_gt_zero(s string) bool {
	if d(s).Cmp(decimal.Zero) > 0 {
		return true
	} else {
		return false
	}
}

func check_number_eq_zero(s string) bool {
	if d(s).Cmp(decimal.Zero) == 0 {
		return true
	} else {
		return false
	}
}

func number(num string) string {
	return d(num).String()
}

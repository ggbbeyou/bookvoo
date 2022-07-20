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

func check_number_lt_zero(s string) bool {
	ss, _ := decimal.NewFromString(s)
	if ss.Cmp(decimal.Zero) < 0 {
		return true
	} else {
		return false
	}
}

func check_number_gt_zero(s string) bool {
	ss, _ := decimal.NewFromString(s)
	if ss.Cmp(decimal.Zero) > 0 {
		return true
	} else {
		return false
	}
}

func check_number_eq_zero(s string) bool {
	ss, _ := decimal.NewFromString(s)
	if ss.Cmp(decimal.Zero) == 0 {
		return true
	} else {
		return false
	}
}

func number(num string) string {
	ss, _ := decimal.NewFromString(num)
	return ss.String()
}

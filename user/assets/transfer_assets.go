package assets

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

func Transfer(db *xorm.Session, enable_transaction bool, from, to int64, symbol_id int, amount string, business_id string, info OpBehavior) (success bool, err error) {
	return transfer(db, enable_transaction, from, to, symbol_id, amount, business_id, info)
}

func transfer(db *xorm.Session, enable_transaction bool, from, to int64, symbol_id int, amount string, business_id string, info OpBehavior) (success bool, err error) {
	if enable_transaction {
		db.Begin()
		defer func() {
			if err != nil {
				logrus.Error(err)
				db.Rollback()
			} else {
				db.Commit()
			}
		}()
	}

	from_user := Assets{UserId: from, SymbolId: symbol_id}
	has_from, err := db.Table(new(Assets)).Where("user_id=? and symbol_id=?", from, symbol_id).ForUpdate().Get(&from_user)
	if err != nil {
		return false, err
	}
	//非根账户检查余额
	if from != ROOTUSERID {
		if check_number_lt_zero(from_user.Available) {
			return false, fmt.Errorf("available balance less than zero")
		}
	}

	to_user := Assets{UserId: to, SymbolId: symbol_id}
	has_to, err := db.Table(new(Assets)).Where("user_id=? and symbol_id=?", to, symbol_id).ForUpdate().Get(&to_user)
	if err != nil {
		return false, err
	}
	from_before := number(from_user.Total)
	from_user.Total = number_sub(from_user.Total, amount)
	from_user.Available = number_sub(from_user.Available, amount)
	if !has_from {
		from_user.Freeze = "0"
		_, err = db.Table(new(Assets)).Insert(&from_user)
	} else {
		_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", from, symbol_id).Update(&from_user)
	}
	if err != nil {
		return false, err
	}

	to_before := number(to_user.Total)
	to_user.Total = number_add(to_user.Total, amount)
	to_user.Available = number_add(to_user.Available, amount)
	if !has_to {
		to_user.Freeze = "0"
		_, err = db.Table(new(Assets)).Insert(&to_user)
	} else {
		_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", to, symbol_id).Update(&to_user)
	}
	if err != nil {
		return false, err
	}

	//双方日志
	from_log := assetsLog{
		UserId:   from,
		SymbolId: symbol_id,
		Before:   from_before,
		Amount:   "-" + amount,
		After:    from_user.Total,
		Info:     fmt.Sprintf("id: %s to: %d info: %s", business_id, to, info),
	}
	_, err = db.Table(new(assetsLog)).Insert(&from_log)
	if err != nil {
		return false, err
	}

	to_log := assetsLog{
		UserId:   to,
		SymbolId: symbol_id,
		Before:   to_before,
		Amount:   amount,
		After:    to_user.Total,
		Info:     fmt.Sprintf("id: %s from: %d info: %s", business_id, from, info),
	}
	_, err = db.Table(new(assetsLog)).Insert(&to_log)
	if err != nil {
		return false, err
	}
	return true, err
}

package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func QueryFreeze(db *xorm.Session, bid string) (*assetFreezeRecord, error) {
	row := assetFreezeRecord{}
	has, err := db.Table(new(assetFreezeRecord)).Where("business_id=?", bid).Get(&row)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, fmt.Errorf("failed to query frozen records")
	}

	return &row, nil
}

func FreezeAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, freeze_amount, business_id string, behavior OpBehavior) (success bool, err error) {
	return freezeAssets(db, enable_transaction, user_id, symbol_id, freeze_amount, business_id, behavior)
}

func FreezeTotalAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, business_id string, behavior OpBehavior) (success bool, err error) {
	return freezeAssets(db, enable_transaction, user_id, symbol_id, "0", business_id, behavior)
}

func freezeAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, freeze_amount, business_id string, behavior OpBehavior) (success bool, err error) {

	if check_number_lt_zero(freeze_amount) {
		return false, fmt.Errorf("freeze amount should be >= 0")
	}

	if enable_transaction {
		db.Begin()
		defer func() {
			if err != nil {
				db.Rollback()
			} else {
				db.Commit()
			}
		}()
	}

	item := Assets{UserId: user_id, SymbolId: symbol_id}
	_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).ForUpdate().Get(&item)
	if err != nil {
		return false, err
	}

	if d(freeze_amount).Equal(d("0")) {
		freeze_amount = d(item.Available).String()
	}

	item.Available = number_sub(item.Available, freeze_amount)
	item.Freeze = number_add(item.Freeze, freeze_amount)

	if check_number_lt_zero(item.Available) {
		return false, fmt.Errorf("available balance less than zero")
	}

	if check_number_lt_zero(item.Freeze) {
		return false, fmt.Errorf("freeze balance less than zero")
	}

	_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).AllCols().Update(&item)

	if err != nil {
		return false, err
	}

	//freeze log
	lg := assetFreezeRecord{
		UserId:       user_id,
		SymbolId:     symbol_id,
		Amount:       freeze_amount,
		FreezeAmount: freeze_amount,
		BusinessId:   business_id,
		Status:       FreezeStatusNew,
		Info:         string(behavior),
	}

	_, err = db.Table(new(assetFreezeRecord)).Insert(&lg)
	if err != nil {
		return false, err
	}

	return true, nil
}

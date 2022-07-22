package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func FreezeAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, freeze_amount, business_id string, info OpBehavior) (success bool, err error) {
	return freezeAssets(db, enable_transaction, user_id, symbol_id, freeze_amount, business_id, info)
}

func freezeAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, freeze_amount, business_id string, info OpBehavior) (success bool, err error) {

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
		freeze_amount = item.Available
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
		Info:         string(info),
	}
	_, err = db.Table(new(assetFreezeRecord)).Insert(&lg)
	if err != nil {
		return false, err
	}

	return true, nil
}

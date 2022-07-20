package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func freezeAssets(db *xorm.Session, user_id int64, symbol_id int, freeze_amount, business_id, info string) (success bool, err error) {

	if !check_number_gt_zero(freeze_amount) {
		return false, fmt.Errorf("freeze amount should be gt zero")
	}

	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	item := Assets{UserId: user_id, SymbolId: symbol_id}
	_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).ForUpdate().Get(&item)
	if err != nil {
		return false, err
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
		Info:         info,
	}
	_, err = db.Table(new(assetFreezeRecord)).Insert(&lg)
	if err != nil {
		return false, err
	}

	return true, nil
}

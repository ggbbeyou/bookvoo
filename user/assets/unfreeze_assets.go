package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func unfreezeAssets(db *xorm.Session, enable_transaction bool, user_id int64, symbol_id int, business_id, unfreeze_amount string) (success bool, err error) {
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

	if !check_number_gt_zero(unfreeze_amount) {
		return false, fmt.Errorf("unfreeze amount should be gt zero")
	}

	row := assetFreezeRecord{UserId: user_id, SymbolId: symbol_id, BusinessId: business_id}

	has, err := db.Table(new(assetFreezeRecord)).Where("user_id=? and symbol_id=? and business_id=?", user_id, symbol_id, business_id).Get(&row)
	if err != nil {
		return false, err
	}

	if !has {
		return false, fmt.Errorf("not found")
	}

	if row.Status == FreezeStatusDone {
		return false, fmt.Errorf("repeat unfreeze")
	}

	row.FreezeAmount = number_sub(row.FreezeAmount, unfreeze_amount)

	if check_number_lt_zero(row.FreezeAmount) {
		return false, fmt.Errorf("unfreeze amount must lt freeze amount")
	}

	if check_number_eq_zero(row.FreezeAmount) {
		row.Status = FreezeStatusDone
	}

	_, err = db.Table(new(assetFreezeRecord)).Where("user_id=? and symbol_id=? and business_id=?", user_id, symbol_id, business_id).AllCols().Update(&row)
	if err != nil {
		return false, err
	}

	//解冻资产为可用
	assets := Assets{UserId: user_id, SymbolId: symbol_id}
	_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).Get(&assets)
	if err != nil {
		return false, err
	}
	assets.Available = number_add(assets.Available, unfreeze_amount)
	assets.Freeze = number_sub(assets.Freeze, unfreeze_amount)

	if check_number_lt_zero(assets.Freeze) {
		return false, fmt.Errorf("freeze amount some wrong")
	}

	_, err = db.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).Update(&assets)
	if err != nil {
		return false, err
	}

	return true, nil
}

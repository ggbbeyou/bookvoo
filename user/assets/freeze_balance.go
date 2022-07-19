package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func freezeAssets(ses *xorm.Session, user_id int64, symbol_id int, freeze_amount, business_id, info string) (bool, error) {
	item := Assets{UserId: user_id, SymbolId: symbol_id}
	has, err := ses.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).ForUpdate().Get(&item)
	if err != nil {
		return false, err
	}

	item.Available = number_sub(item.Available, freeze_amount)
	item.Freezed = number_add(item.Freezed, freeze_amount)

	if check_amount_lt_zero(item.Available) {
		return false, fmt.Errorf("available balance less than zero")
	}

	if check_amount_lt_zero(item.Freezed) {
		return false, fmt.Errorf("freeze balance less than zero")
	}

	if !has {
		_, err = ses.Table(new(Assets)).Insert(&item)
	} else {
		_, err = ses.Table(new(Assets)).Where("user_id=? and symbol_id=?", user_id, symbol_id).Update(&item)
	}

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
	_, err = ses.Table(new(assetsLog)).Insert(&lg)
	if err != nil {
		return false, err
	}

	return true, nil
}

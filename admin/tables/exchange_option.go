package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetExchangeOptionTable(ctx *context.Context) table.Table {

	exchangeOption := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := exchangeOption.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Int).
		FieldFilterable()
	info.AddField("Symbol", "symbol", db.Varchar)
	info.AddField("Name", "name", db.Varchar)
	// info.AddField("Target_symbol_id", "target_symbol_id", db.Int)
	// info.AddField("Standard_symbol_id", "standard_symbol_id", db.Int)
	// info.AddField("Price_prec", "price_prec", db.Int)
	// info.AddField("Qty_prec", "qty_prec", db.Int)
	// info.AddField("Allow_min_qty", "allow_min_qty", db.Decimal)
	// info.AddField("Allow_max_qty", "allow_max_qty", db.Decimal)
	// info.AddField("Allow_min_amount", "allow_min_amount", db.Decimal)
	// info.AddField("Allow_max_amount", "allow_max_amount", db.Decimal)
	// info.AddField("Fee_rate", "fee_rate", db.Decimal)
	info.AddField("Status", "status", db.Int)
	info.AddField("Create_time", "create_time", db.Timestamp)
	info.AddField("Update_time", "update_time", db.Timestamp)

	info.SetTable("exchange_option").SetTitle("ExchangeOption").SetDescription("ExchangeOption")

	formList := exchangeOption.GetForm()
	formList.AddField("Id", "id", db.Int, form.Default)
	formList.AddField("Symbol", "symbol", db.Varchar, form.Text)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Target_symbol_id", "target_symbol_id", db.Int, form.Number)
	formList.AddField("Standard_symbol_id", "standard_symbol_id", db.Int, form.Number)
	formList.AddField("Price_prec", "price_prec", db.Int, form.Number)
	formList.AddField("Qty_prec", "qty_prec", db.Int, form.Number)
	formList.AddField("Allow_min_qty", "allow_min_qty", db.Decimal, form.Text)
	formList.AddField("Allow_max_qty", "allow_max_qty", db.Decimal, form.Text)
	formList.AddField("Allow_min_amount", "allow_min_amount", db.Decimal, form.Text)
	formList.AddField("Allow_max_amount", "allow_max_amount", db.Decimal, form.Text)
	formList.AddField("Fee_rate", "fee_rate", db.Decimal, form.Text)
	formList.AddField("Status", "status", db.Int, form.Number)
	formList.AddField("Create_time", "create_time", db.Timestamp, form.Datetime)
	formList.AddField("Update_time", "update_time", db.Timestamp, form.Datetime)

	formList.SetTable("exchange_option").SetTitle("ExchangeOption").SetDescription("ExchangeOption")

	return exchangeOption
}

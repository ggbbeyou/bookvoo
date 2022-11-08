package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAssetsTable(ctx *context.Context) table.Table {

	assets := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := assets.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("User_id", "user_id", db.Bigint)
	info.AddField("Symbol_id", "symbol_id", db.Int)
	info.AddField("Total", "total", db.Decimal)
	info.AddField("Freeze", "freeze", db.Decimal)
	info.AddField("Available", "available", db.Decimal)
	info.AddField("Create_time", "create_time", db.Timestamp)
	info.AddField("Update_time", "update_time", db.Timestamp)

	info.SetTable("assets").SetTitle("Assets").SetDescription("Assets")

	formList := assets.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("User_id", "user_id", db.Bigint, form.Number)
	formList.AddField("Symbol_id", "symbol_id", db.Int, form.Number)
	formList.AddField("Total", "total", db.Decimal, form.Text)
	formList.AddField("Freeze", "freeze", db.Decimal, form.Text)
	formList.AddField("Available", "available", db.Decimal, form.Text)
	formList.AddField("Create_time", "create_time", db.Timestamp, form.Datetime)
	formList.AddField("Update_time", "update_time", db.Timestamp, form.Datetime)

	formList.SetTable("assets").SetTitle("Assets").SetDescription("Assets")

	return assets
}

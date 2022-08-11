package assets

import (
	"time"
)

// 用户资产冻结记录
type FreezeStatus int
type OpBehavior string

const (
	FreezeStatusNew  FreezeStatus = 0
	FreezeStatusDone FreezeStatus = 1

	Behavior_Trade        OpBehavior = "trade"
	Behavior_Recharge     OpBehavior = "recharge"
	Behavior_Withdraw     OpBehavior = "withdraw"
	Behavior_Transfer     OpBehavior = "transfer"
	Behavior_DemoRecharge OpBehavior = "demo_recharge"
)

// 用户资产余额表
type Assets struct {
	Id         int64     `xorm:"pk autoincr bigint" json:"-"`
	UserId     int64     `xorm:"notnull unique(userid_symbol)" json:"-"`
	SymbolId   int       `xorm:"notnull unique(userid_symbol)" json:"-"`
	Total      string    `xorm:"decimal(40,20) default(0) notnull" json:"total"`
	Freeze     string    `xorm:"decimal(40,20) default(0) notnull" json:"freeze"`
	Available  string    `xorm:"decimal(40,20) default(0) notnull" json:"available"`
	CreateTime time.Time `xorm:"timestamp created" json:"-"`
	UpdateTime time.Time `xorm:"timestamp updated" json:"update_time"`
}

// 用户资产变动记录
type assetsLog struct {
	Id         int64      `xorm:"pk autoincr bigint"`
	UserId     int64      `xorm:"index notnull"`
	SymbolId   int        `xorm:"index notnull"`
	Before     string     `xorm:"decimal(40,20) default(0)"`               // 变动前
	Amount     string     `xorm:"decimal(40,20) default(0)"`               // 变动数
	After      string     `xorm:"decimal(40,20) default(0)"`               // 变动后
	BusinessId string     `xorm:"varchar(100) index(business_id) notnull"` //业务相关的id
	Behavior   OpBehavior `xorm:"varchar(15)"`
	Info       string     `xorm:"varchar(64)"`
	CreateTime time.Time  `xorm:"timestamp created"`
	UpdateTime time.Time  `xorm:"timestamp updated"`
}

type assetFreezeRecord struct {
	Id           int64        `xorm:"pk autoincr bigint"`
	UserId       int64        `xorm:"bigint index notnull"`
	SymbolId     int          `xorm:"index notnull"`
	Amount       string       `xorm:"decimal(40,20) default(0) notnull"`        // 冻结总量
	FreezeAmount string       `xorm:"decimal(40,20) default(0) notnull"`        // 冻结着的量
	Status       FreezeStatus `xorm:"tinyint(1)"`                               // 状态 冻结中, 已解冻
	BusinessId   string       `xorm:"varchar(100) unique(business_id) notnull"` //业务相关的id
	Info         string       `xorm:"varchar(64)"`
	CreateTime   time.Time    `xorm:"timestamp created"`
	UpdateTime   time.Time    `xorm:"timestamp updated"`
}

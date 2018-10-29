// echoRouteSaveflow
package main

import (
	"time"
)

func dbQueryShopById(id int64) (isExit bool, shop Shop) {

	querydatas := []Shop{}
	db.Table("t_shop").Select(" * ").Where(" shop_id = ? ", id).Find(&querydatas)
	if len(querydatas) > 0 {
		return true, querydatas[0]
	} else {
		return false, querydatas[0]
	}
}

type Shop struct {
	ShopId      int64     `gorm:"column:shop_id;primary_key" validate:"-"`
	ShopName    string    `gorm:"column:shop_name" validate:"required,max=64"`
	Address     string    `gorm:"column:address" validate:"required,max=64"`
	Operator_id int64     `gorm:"column:operator_id" validate:"-"`
	CreTime     time.Time `gorm:"column:cre_time" validate:"-"`
	UpdTime     time.Time `gorm:"column:upd_time" validate:"-"`
	Remark      string    `gorm:"column:remark" validate:"max=256"`
}
